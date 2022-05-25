package dbyml

import (
	"fmt"
	"os"
	"text/template"
)

const ConfigurationTemplate = `---
# This is an auto-generated configuration file .
image:
  name: {{ or .ImageInfo.Basename "go-dbyml-sample" }}
  tag: {{ or .ImageInfo.Tag "latest" }}
  path: {{ or .ImageInfo.Path "." }}
  dockerfile: {{ or .ImageInfo.Dockerfile "Dockerfile" }}
  build_args:
    {{- range $k, $v := .ImageInfo.BuildArgs }}
    {{ $k }}: {{ $v }}
    {{- end }}
  label:
    {{- range $k, $v := .ImageInfo.Labels }}
    {{ $k }}: {{ $v }}
    {{- end }}
  docker_host: {{ or .ImageInfo.DockerHost "unix:/var/run/docker.sock" }}

build:
  target: {{ or .BuildInfo.Target "''" }}
  no_cache: {{ or .BuildInfo.NoCache false }}
  verbose: {{ or .BuildInfo.Verbose true }}

registry:
  enabled: {{ or .RegistryInfo.Enabled false }}
  host: {{ or .RegistryInfo.Host "myregistry.com:5000" }}
  project: {{ or .RegistryInfo.Project "" }}
  auth:
    {{- range $k, $v := .RegistryInfo.Auth }}
    {{ $k }}: {{ $v }}
    {{- end }}

buildx:
  enabled: false
  debug: false
  instance: multi-builder
  use_existing_instance: true
  platform:
    - linux/amd64
    - linux/arm64
  type: registry
  pull_output: true
  remove_instance: false
  driver_opt:
    network: host
  config:
    http: true
`

func MakeTemplate(config *Configuration) {
	tmpl := template.Must(template.New("ConfigurationTemplate").Parse(ConfigurationTemplate))

	file, _ := os.Create("dbyml.yml")
	err := tmpl.Execute(file, config)
	if err != nil {
		panic(err)
	}
	fmt.Println("Create dbyml.yml. Check the contents and edit it according to your docker image.")
}
