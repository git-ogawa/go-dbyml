package dbyml

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
)

// ImageInfo defines docker image information.
type ImageInfo struct {
	Basename   string             `yaml:"name"`        // Image name
	Tag        string             `yaml:"tag"`         // Image tag
	ImageName  string             `yaml:"image_name"`  // Image name such as `go-dbyml:latest`
	Context    string             `yaml:"path"`        // Path to the directory where Dockerfile exists, is equivalent to build context
	Dockerfile string             `yaml:"dockerfile"`  // Dockerfile filename
	BuildArgs  map[string]*string `yaml:"build_args"`  // Build-args to be passed to image on build
	Labels     map[string]string  `yaml:"label"`       // Labels to be passed to image on build
	DockerHost string             `yaml:"docker_host"` // Docker host such as "unix:///var/run/docker.sock"

	DockerfilePath string
	Registry       RegistryInfo
	BuildInfo      BuildInfo
	FullName       string
	DockerClient   *client.Client
}

// NewImageInfo creates a new ImageInfo struct with default values.
func NewImageInfo() *ImageInfo {
	image := new(ImageInfo)
	image.Tag = "latest"
	image.Context = "."
	image.Dockerfile = "Dockerfile"
	image.DockerHost = "unix:///var/run/docker.sock"
	image.DockerfilePath = image.Context + "/" + image.Dockerfile
	image.Registry = *NewRegistryInfo()
	image.BuildInfo = *NewBuildInfo()
	return image
}

// SetProperties sets some properties when build an image.
func (image *ImageInfo) SetProperties() {
	image.ImageName = image.Basename + ":" + image.Tag
	image.DockerfilePath = image.Context + "/" + image.Dockerfile
	image.SetDockerClient()
}

// SetDockerClient initializes docker api client for the specified host.
func (image *ImageInfo) SetDockerClient() {
	image.DockerClient, _ = client.NewClientWithOpts(
		client.WithHost(image.DockerHost),
		client.WithAPIVersionNegotiation(),
	)
}

// ShowProperties shows the current settings related to image build.
func (image ImageInfo) ShowProperties() {
	rv := reflect.ValueOf(image)
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		kind := field.Type.Kind()
		value := rv.FieldByName(field.Name)
		if kind == reflect.Map || kind == reflect.Slice {
			showMapElement(field.Name, value.MapRange())
		} else if kind == reflect.String && value.Interface() != "" {
			fmt.Printf("%-30v: %v\n", field.Name, value)
		}
	}
}

// Build runs image build.
func (image *ImageInfo) Build() error {
	buf := GetBuildContext(image.Context)
	tar := bytes.NewReader(buf.Bytes())
	ctx := context.Background()

	options := types.ImageBuildOptions{
		NoCache:    image.BuildInfo.NoCache,
		Dockerfile: image.DockerfilePath,
		Remove:     true,
		BuildArgs:  image.BuildArgs,
		Labels:     image.Labels,
		Target:     image.BuildInfo.Target,
		Tags:       []string{image.ImageName},
	}

	res, err := image.DockerClient.ImageBuild(ctx, tar, options)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	termFd, isTerm := term.GetFdInfo(os.Stderr)
	err = jsonmessage.DisplayJSONMessagesStream(res.Body, os.Stderr, termFd, isTerm, nil)
	return err
}

// SetFullImageName sets image name for pushing to a registry.
func (image *ImageInfo) SetFullImageName() {
	if image.Registry.Project != "" {
		image.FullName = image.Registry.Host + "/" + image.Registry.Project + "/" + image.ImageName
	} else {
		image.FullName = image.Registry.Host + "/" + image.ImageName
	}
}

// Push runs image push to a registry.
func (image *ImageInfo) Push() error {
	ctx := context.Background()
	image.AddTag()

	opts := types.ImagePushOptions{All: false, RegistryAuth: image.Registry.BasicAuth()}

	res, err := image.DockerClient.ImagePush(ctx, image.FullName, opts)
	if err != nil {
		return err
	}
	defer res.Close()

	termFd, isTerm := term.GetFdInfo(os.Stderr)
	err = jsonmessage.DisplayJSONMessagesStream(res, os.Stderr, termFd, isTerm, nil)
	return err
}

// AddTag adds a tag containing the registry name to a built image.
func (image *ImageInfo) AddTag() {
	ctx := context.Background()
	image.SetFullImageName()

	err := image.DockerClient.ImageTag(ctx, image.ImageName, image.FullName)
	if err != nil {
		panic(err)
	}
}
