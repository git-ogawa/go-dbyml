package dbyml

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	dockerStrSlice "github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/moby/api/types/strslice"
	"github.com/moby/moby/pkg/stdcopy"
	"github.com/moby/term"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// The container image used in buildkitd container
const buildkitImageName = "moby/buildkit:v0.10.3"

// BuildkitInfo defines setting on build with buildkit.
type BuildkitInfo struct {
	Enabled  bool                   `yaml:"enabled"`
	Output   map[string]interface{} `yaml:"output"`
	Cache    map[string]interface{} `yaml:"cache"`
	Platform []string               `yaml:"platform"`
	Remove   bool                   `yaml:"remove"`
}

// NewBuildkitInfo makes BuildkitInfo object with default values.
func NewBuildkitInfo() *BuildkitInfo {
	build := new(BuildkitInfo)
	build.Enabled = false
	build.Output = map[string]interface{}{}
	build.Cache = map[string]interface{}{}
	build.Remove = true
	return build
}

// ParseOptions parses options related to buildkit and sets buildctl args.
func (buildkit *BuildkitInfo) ParseOptions(imageInfo ImageInfo) []string {
	var opts []string
	var cmd string

	// Output
	cmd = fmt.Sprintf("type=%s,name=%s,push=true", buildkit.Output["type"], buildkit.Output["name"])
	if buildkit.Output["insecure"] == true {
		cmd = fmt.Sprintf("%s,registry.insecure=true", cmd)
	}
	opts = append(opts, "--output", cmd)

	// Cache
	if len(buildkit.Cache) != 0 {
		exportCache := buildkit.Cache["export"].(map[string]string)["type"]
		if exportCache == "inline" {
			cmd = "type=inline"
		} else if exportCache == "registry" {
			cmd = fmt.Sprintf("type=registry,ref=%s", buildkit.Cache["export"].(map[string]string)["value"])
		}
		opts = append(opts, "--export-cache", cmd)

		importCache := buildkit.Cache["import"].(map[string]string)["type"]
		if importCache == "registry" {
			cmd = fmt.Sprintf(
				"type=registry,ref=%s",
				buildkit.Cache["import"].(map[string]string)["value"],
			)
		}
		opts = append(opts, "--import-cache", cmd)
	}

	// Platform
	if len(buildkit.Platform) != 0 {
		cmd = fmt.Sprintf("platform=%s", strings.Join(buildkit.Platform, ","))
		opts = append(opts, "--opt", cmd)
	}

	// Other imageInfo options
	if len(imageInfo.Labels) != 0 {
		for k, v := range imageInfo.Labels {
			cmd = fmt.Sprintf("label:%s=%s", k, v)
			opts = append(opts, "--opt", cmd)
		}
	}

	if len(imageInfo.BuildArgs) != 0 {
		for k, v := range imageInfo.BuildArgs {
			cmd = fmt.Sprintf("build-arg:%s=%s", k, *v)
			opts = append(opts, "--opt", cmd)
		}
	}

	return opts
}

// Builder describes a container information on buildkit
type Builder struct {
	Name           string                // The name of builder container
	Image          BuildkitImage         // The image of builder container
	ID             string                // The container ID
	Config         *container.Config     // Container config
	HostConfig     *container.HostConfig // Container host config
	Context        string                // The build context in builder
	DockerfilePath string                // The path to Dockerfile in builder
	Cmd            []string              // The command executed in the builder
	Client         *client.Client        // Docker client for connecting to builder
}

// NewBuilder creates a builder object with the default values.
func NewBuilder() (builder *Builder) {
	builder = new(Builder)
	builder.Name = "dbyml-buildkit-builder"
	builder.Image = BuildkitImage{buildkitImageName}
	builder.Context = "/tmp"
	builder.DockerfilePath = "/tmp"
	builder.Cmd = []string{
		"buildctl",
		"build",
		"--frontend",
		"dockerfile.v0",
		"--local",
		fmt.Sprintf("context=%s", builder.Context),
		"--local",
		fmt.Sprintf("dockerfile=%s", builder.DockerfilePath),
	}
	builder.Config = &container.Config{
		Image: builder.Image.Name,
		Entrypoint: dockerStrSlice.StrSlice(
			[]string{"buildkitd", "--config", "/etc/buildkitd.toml"},
		),
	}
	builder.HostConfig = &container.HostConfig{
		NetworkMode: "host",
		Privileged:  true,
	}
	builder.Client, _ = client.NewClientWithOpts(client.FromEnv)
	return builder
}

// AddCmd adds arguments passed to buildctl.
func (builder *Builder) AddCmd(cmd ...string) {
	builder.Cmd = append(builder.Cmd, cmd...)
}

// Setup creates a builder container and copy setting toml into the builder.
func (builder *Builder) Setup(config *RegistryInfo) error {
	err := builder.Create()
	if err != nil {
		return err
	}

	path, err := MakeBuildkitToml(config)
	if err != nil {
		return err
	}

	if err = builder.CopyFiles(path, "/etc"); err != nil {
		return err
	}

	if err = os.Remove(path); err != nil {
		return err
	}

	return nil
}

// Exists checks if a builder container exists.
func (builder *Builder) Exists() bool {
	json, err := builder.Inspect()
	if err != nil {
		panic(err)
	}
	if json.ContainerJSONBase != nil {
		return true
	}
	return false
}

// SetContainerID sets container ID of a builder.
func (builder *Builder) SetContainerID() error {
	json, err := builder.Inspect()
	if err != nil {
		return err
	}
	builder.ID = json.ContainerJSONBase.ID
	return nil
}

// Inspect gets a builder information.
func (builder *Builder) Inspect() (types.ContainerJSON, error) {
	var json types.ContainerJSON
	ret, err := builder.Client.ContainerList(
		context.Background(),
		types.ContainerListOptions{All: true},
	)
	if err != nil {
		return json, err
	}

	target := "/" + builder.Name
	if len(ret) > 0 {
		for _, container := range ret {
			for _, name := range container.Names {
				if name == target {
					return builder.Client.ContainerInspect(context.Background(), container.ID)
				}
			}
		}
	}
	return json, nil
}

// Create creates a builder container.
func (builder *Builder) Create() error {
	body, err := builder.Client.ContainerCreate(
		context.Background(),
		builder.Config,
		builder.HostConfig,
		&network.NetworkingConfig{},
		&specs.Platform{},
		builder.Name,
	)
	if err != nil {
		return err
	}
	builder.ID = body.ID
	return nil
}

// Start starts a builder container.
func (builder *Builder) Start() error {
	return builder.Client.ContainerStart(
		context.Background(),
		builder.ID,
		types.ContainerStartOptions{},
	)
}

// Stop stops a builder container.
func (builder *Builder) Stop() error {
	timeout := time.Second * 60
	return builder.Client.ContainerStop(context.Background(), builder.ID, &timeout)
}

// Remove removes a builder container.
func (builder *Builder) Remove() error {
	return builder.Client.ContainerRemove(
		context.Background(),
		builder.ID,
		types.ContainerRemoveOptions{RemoveVolumes: true, RemoveLinks: false, Force: true},
	)
}

// CopyFiles copies directory in client to builder container.
// If the directory contains some other directories, copy them recursively.
func (builder *Builder) CopyFiles(path string, dst string) error {
	opts := types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
		CopyUIDGID:                false,
	}
	return builder.Client.CopyToContainer(
		context.Background(),
		builder.ID,
		dst,
		GetBuildkitContext(path),
		opts,
	)
}

// Build builds a image in a builder.
func (builder *Builder) Build(debug bool) error {
	if debug {
		fmt.Println("The following command will be run in buildkit container.")
		re := regexp.MustCompile(`\s{1}-{2}`)
		cmd := strings.Join(builder.Cmd, " ")
		cmd = re.ReplaceAllString(cmd, "\n\t--")
		fmt.Println(cmd)
	}
	return builder.Exec(builder.Cmd)
}

// Exec runs a command in buildkit container.
func (builder *Builder) Exec(cmd []string) error {
	execConfig := types.ExecConfig{
		Privileged:   true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Detach:       false,
		Cmd:          strslice.StrSlice(cmd),
	}

	// Create a new exec configuration to run an exec process.
	res, err := builder.Client.ContainerExecCreate(context.Background(), builder.Name, execConfig)
	if err != nil {
		return err
	}

	// Run the exec process and attach it.
	hijackRes, _ := builder.Client.ContainerExecAttach(
		context.Background(),
		res.ID,
		types.ExecStartCheck{},
	)
	defer func() error {
		return hijackRes.Conn.Close()
	}()
	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, hijackRes.Reader)
	if err != nil {
		return err
	}
	return nil
}

// BuildkitImage defines a docker image used in a builder container.
type BuildkitImage struct {
	Name string
}

// Exists checks if the image exists on host.
func (buildkit *BuildkitImage) Exists() bool {
	cli, _ := client.NewClientWithOpts(client.FromEnv)
	imgs, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	for _, img := range imgs {
		if contains(img.RepoTags, buildkit.Name) {
			return true
		}
	}
	return false
}

// Pull pulls a buildkit image from official dockerhub.
func (buildkit *BuildkitImage) Pull() error {
	cli, _ := client.NewClientWithOpts(client.FromEnv)
	ret, err := cli.ImagePull(context.Background(), buildkit.Name, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	termFd, isTerm := term.GetFdInfo(os.Stderr)
	return jsonmessage.DisplayJSONMessagesStream(ret, os.Stderr, termFd, isTerm, nil)
}

func contains(s []string, tag string) bool {
	for _, v := range s {
		if tag == v {
			return true
		}
	}
	return false
}
