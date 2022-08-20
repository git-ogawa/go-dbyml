package dbyml

import (
	"fmt"
	"os"
	"text/template"
)

// ConfigurationTemplate is a template of dbyml settings.
const ConfigurationTemplate = `---
# file: dbyml.yml
# This is a configuration file used by go-dbyml.
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
  path: {{ or .ImageInfo.Context "." }}

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

# The buildkit section manages the settings when build with buildkit.
buildkit:
  # enabled: Set true to enable build with buildkit.
  enabled: {{ or .BuildkitInfo.Enabled false }}
  # output: The output field sets output format of the image to be built
  output:
    # type: Type of output image.
    type: {{ or .BuildkitInfo.Output.Type "image" }}
    # name: Set registry and image. e.g. [registry]:[port]/[project]/[image]:[tag]
    name: {{ or .BuildkitInfo.Output.Name "myregistry.com/go-dbyml-sample:latest" }}
    # insecure: Set true if push the image insecure registry such as insecure private registry
    insecure: {{ or .BuildkitInfo.Output.Insecure false }}
  # The cache field sets import and export build cache.
  cache:
    # export: Export settings of build cache.
    # Type must be either inline or registry. Set inline to export the cache embed with the image and pushing them to registry together.
    # Set registry to export build cache to the specified registry. In this case, set the registry in value field. e.g. [registry]:[port]/[project]/[image]:[tag]
    export:
      type: {{ or .BuildkitInfo.Cache.Export.Type "inline" }}
      value: {{ or .BuildkitInfo.Cache.Export.Value "''" }}
    # import: Import settings of build cache.
    # Type must be registry. Set registry to import build cache from the specified registry. In this case, set the registry in value field. e.g. [registry]:[port]/[project]/[image]:[tag].
    import:
      type: {{ or .BuildkitInfo.Cache.Import.Type "registry" }}
      value: {{ or .BuildkitInfo.Cache.Import.Value "myregistry.com/go-dbyml-sample:latest" }}
  # platform: Set the list of architectures if want to build  a image that support multi-platform.
  platform:
  # remove: Set true to remove a builder container after build is successfully completed.
  remove: {{ or .BuildkitInfo.Remove true }}
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

// BuildkitdTomlTemplate is a template of buildkit settings.
const BuildkitdTomlTemplate = `[registry."{{ .Host }}"]
  {{- if .Insecure }}
  insecure = true
  {{- end }}
  {{- if .Auth.ca_cert }}
  ca_cert = ["{{ .Auth.ca_cert }}"]
  {{- end }}
`

// MakeBuildkitToml makes buildkitd.toml from a template.
func MakeBuildkitToml(config *RegistryInfo) (string, error) {
	tmpl := template.Must(template.New("BuildkitdTomlTemplate").Parse(BuildkitdTomlTemplate))

	file, _ := os.Create("buildkitd.toml")
	err := tmpl.Execute(file, config)
	if err != nil {
		return "", err
	}
	return "buildkitd.toml", nil
}
