package dbyml

import (
	"fmt"
	"os"
	"text/template"
)

// ConfigurationTemplate is a template of dbyml settings.
const ConfigurationTemplate = `---
# file: dbyml.yml
# This is a configration file used by go-dbyml.
# The settings related to image build are written in yaml syntax.
# The values set in each field, so you can edit them according to your build settings.
# This is automatically generated.


# The image section manages docker image attributes.
image:
  # name: Image name. This field is required.
  name: {{ or .ImageInfo.Basename "go-dbyml-sample" }}

  # tag: Image tag.
  tag: {{ or .ImageInfo.Tag "latest" }}

  # path: Path to a directory containing Dockerfile.
  path: {{ or .ImageInfo.Path "." }}

  # dockerfile: Filename of Dockerfile.
  dockerfile: {{ or .ImageInfo.Dockerfile "Dockerfile" }}

  # build_args: Arguments corresponding to build-args of docker build.
  # Set list of key:value
  build_args:
    {{- range $k, $v := .ImageInfo.BuildArgs }}
    {{ $k }}: {{ $v }}
    {{- end }}

  # label: Arguments corresponding to label of docker build.
  # Use quotation when the key name includes "." such as com.test.name
  # Set list of key:value
  label:
    {{- range $k, $v := .ImageInfo.Labels }}
    {{ $k }}: {{ $v }}
    {{- end }}

  # docker_host: URL to the Docker server.
  # Set protocol:hostname:port for example unix:///var/run/docker.sock or tcp://127.0.0.1:1234.
  # Default to unix:///var/run/docker.sock
  docker_host: {{ or .ImageInfo.DockerHost "unix:/var/run/docker.sock" }}


# The build section manages some options on build such as using build-cache or showing build information.
build:
  # target: Name of the build-stage to build in a multi-stage Dockerfile.
  target: {{ or .BuildInfo.Target "''" }}

  # no_cache: Set true not to use build cache when build the image.
  # default: false
  no_cache: {{ or .BuildInfo.NoCache false }}

  # verbose: Set true to show build settings on build.
  # default: true
  verbose: {{ or .BuildInfo.Verbose true }}

# The registry section manages the information about registry to which the image push.
registry:
  # enabled: Enable push to a registry. Set false not to push the image
  # to the registry even if these fields are set.
  enabled: {{ or .RegistryInfo.Enabled false }}

  # host: Registry name or ip address and port.
  host: {{ or .RegistryInfo.Host "myregistry.com:5000" }}

  # project:
  # When set the project, the image to be pushed to the registry will be {host}:{port}/{project}/{name}:{tag},
  # otherwise {host}:{port}/{name}:{tag}
  project: {{ or .RegistryInfo.Project "" }}

  # auth: Credentials for a registry. This will be used when push a image to basic-auth registry.
  auth:
    {{- range $k, $v := .RegistryInfo.Auth }}
    {{ $k }}: {{ $v }}
    {{- end }}
`

// MakeTemplate makes a dbyml setting file from a template.
func MakeTemplate(config *Configuration) {
	tmpl := template.Must(template.New("ConfigurationTemplate").Parse(ConfigurationTemplate))

	file, _ := os.Create("dbyml.yml")
	err := tmpl.Execute(file, config)
	if err != nil {
		panic(err)
	}
	fmt.Println("Create dbyml.yml. Check the contents and edit it according to your docker image.")
}
