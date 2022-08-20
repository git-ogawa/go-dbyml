package dbyml

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	buildkitInfo := NewBuildkitInfo()

	buildkitInfo.Output["type"] = "image"
	buildkitInfo.Output["name"] = "myregistry.com:5000/test:latest"
	buildkitInfo.Output["value"] = map[string]string{"type": "myregistry.com:5000/test:latest"}
	buildkitInfo.Cache["export"] = map[string]string{"type": "inline"}
	buildkitInfo.Cache["import"] = map[string]string{
		"type":  "registry",
		"value": "myregistry.com:5000/test:latest",
	}
	buildkitInfo.Platform = []string{"linux/amd64", "linux/arm64"}

	imageInfo := NewImageInfo()
	arg := "value1"
	imageInfo.BuildArgs = map[string]*string{"args1": &arg}
	imageInfo.Labels = map[string]string{"label1": "label_value"}

	expected := []string{
		"--output",
		"type=image,name=myregistry.com:5000/test:latest,push=true",
		"--export-cache",
		"type=inline",
		"--import-cache",
		"type=registry,ref=myregistry.com:5000/test:latest",
		"--opt",
		"platform=linux/amd64,linux/arm64",
		"--opt",
		"label:label1=label_value",
		"--opt",
		"build-arg:args1=value1",
	}
	cmd := buildkitInfo.ParseOptions(*imageInfo)
	assert.Equal(t, reflect.DeepEqual(cmd, expected), true)
}

func TestBuilderCreate(t *testing.T) {
	builder := NewBuilder()
	builder.Name = "gotest-builder"

	builder.Create()
	status, _ := builder.Inspect()
	assert.Equal(t, status.ContainerJSONBase.State.Status, "created")
	builder.Remove()
}

func TestBuilderUseExisting(t *testing.T) {
	builder := NewBuilder()
	builder.Name = "gotest-builder"

	builder.Create()
	assert.Equal(t, builder.Exists(), true)
	builder.SetContainerID()
	builder.Start()
	time.Sleep(time.Second * 2)
	status, _ := builder.Inspect()
	assert.Equal(t, status.ContainerJSONBase.State.Status, "running")
	builder.Remove()
}

func TestBuilderStop(t *testing.T) {
	builder := NewBuilder()
	builder.Name = "gotest-builder"

	builder.Create()
	builder.Start()
	builder.Stop()
	status, _ := builder.Inspect()
	assert.Equal(t, status.ContainerJSONBase.State.Status, "exited")
	builder.Remove()
}

// Build a image with buildkitd
func TestBuilderBuild(t *testing.T) {
	pwd, _ := os.Getwd()
	root, _ := filepath.Abs("../")
	os.Chdir(root)

	builder := NewBuilder()
	builder.Name = "gotest-builder"
	registry := NewRegistryInfo()

	builder.Setup(registry)
	builder.Start()
	time.Sleep(time.Second * 2)
	builder.CopyFiles("testdata/dockerfile_buildkit", "/tmp")
	builder.Build(true)
	builder.Remove()
	os.Chdir(pwd)
}

func TestImagePull(t *testing.T) {
	builder := NewBuilder()
	builder.Image.Exists()
	err := builder.Image.Pull()
	if err != nil {
		panic(err)
	}
}
