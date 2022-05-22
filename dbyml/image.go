package dbyml

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// ImageInfo defines docker image information.
type ImageInfo struct {
	// Image name
	Basename string `yaml:"name"`

	// Image tag
	Tag string `yaml:"tag"`

	// Image name such as `go-dbyml:latest`
	ImageName string

	// Path to the directory where Dockerfile exists
	Path string `yaml:"path"`

	// Dockerfile filename
	Dockerfile string `yaml:"dockerfile"`

	// Build-args to be passed to image on build
	BuildArgs map[string]*string `yaml:"build_args"`

	// Labels to be passed to image on build
	Labels map[string]string `yaml:"label"`

	// Docker host such as "unix:/var/run/docker.sock"
	DockerHost string `yaml:"docker_host"`

	FilePath  string
	Registry  RegistryInfo
	BuildInfo BuildInfo
	FullName  string
}

func NewImageInfo() *ImageInfo {
	image := new(ImageInfo)
	image.Tag = "latest"
	image.Path = "."
	image.Dockerfile = "Dockerfile"
	image.DockerHost = ""
	image.FilePath = image.Path + "/" + image.Dockerfile
	image.Registry = *NewRegistryInfo()
	image.BuildInfo = *NewBuildInfo()
	return image
}

func (image *ImageInfo) SetProperties() {
	image.ImageName = image.Basename + ":" + image.Tag
	image.FilePath = image.Path + "/" + image.Dockerfile
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
			ShowMapElement(field.Name, value.MapRange())
		} else if kind == reflect.String && value.Interface() != "" {
			fmt.Printf("%-30v: %v\n", field.Name, value)
		}
	}
}

// Build runs image build.
func (image *ImageInfo) Build() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	buf := GetTarContext(image.FilePath)
	tar := bytes.NewReader(buf.Bytes())
	ctx := context.Background()

	options := types.ImageBuildOptions{
		Context:    tar,
		NoCache:    image.BuildInfo.NoCache,
		Dockerfile: image.FilePath,
		Remove:     true,
		BuildArgs:  image.BuildArgs,
		Labels:     image.Labels,
		Target:     image.BuildInfo.Target,
		Tags:       []string{image.ImageName},
	}

	res, err := cli.ImageBuild(ctx, tar, options)
	if err != nil {
		log.Fatal(err, " :unable to build docker image")
	}
	defer res.Body.Close()

	_, err = io.Copy(os.Stdout, res.Body)
	if err != nil {
		log.Fatal(err, " :unable to read image build response")
	}
}

func (image *ImageInfo) SetFullImageName() {
	if image.Registry.Project != "" {
		image.FullName = image.Registry.Host + "/" + image.Registry.Project + "/" + image.ImageName
	} else {
		image.FullName = image.Registry.Host + "/" + image.ImageName
	}
}

// AddTag add a tag containing the registry  image to a built image.
func (image *ImageInfo) AddTag() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	image.SetFullImageName()

	err = cli.ImageTag(ctx, image.ImageName, image.FullName)
	if err != nil {
		log.Fatal(err, " :unable to add tag")
	}
}

// Push runs image push to a registry.
func (image *ImageInfo) Push() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	image.AddTag()

	opts := types.ImagePushOptions{All: false, RegistryAuth: image.Registry.BasicAuth()}

	res, err := cli.ImagePush(ctx, image.FullName, opts)
	if err != nil {
		log.Fatal(err, " :unable to push docker image")
	}

	_, err = io.Copy(os.Stdout, res)
	if err != nil {
		log.Fatal(err, " :unable to read image build response")
	}
}
