package dbyml

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Configuration defines the hierarchy of the settings in config file.
type Configuration struct {
	ImageInfo    ImageInfo    `yaml:"image"`
	BuildInfo    BuildInfo    `yaml:"build"`
	RegistryInfo RegistryInfo `yaml:"registry"`
	BuildkitInfo BuildkitInfo `yaml:"buildkit"`
}

// NewConfiguration makes Configuration struct with default values.
func NewConfiguration() *Configuration {
	config := new(Configuration)
	config.ImageInfo = *NewImageInfo()
	config.BuildInfo = *NewBuildInfo()
	config.RegistryInfo = *NewRegistryInfo()
	config.BuildkitInfo = *NewBuildkitInfo()
	return config
}

// ShowConfig shows the current Configuration to stdout.
func (config *Configuration) ShowConfig() {
	PrintCenter("Build info", 30, "-")
	config.ImageInfo.ShowProperties()
	fmt.Println()
	PrintCenter("Registry info", 30, "-")
	config.RegistryInfo.ShowProperties()
}

// BuildInfo defines some options related to setting or progress on image build.
type BuildInfo struct {
	Target  string `yaml:"target"`
	NoCache bool   `yaml:"no_cache"`
	Verbose bool   `yaml:"verbose"`
}

// NewBuildInfo makes Configuration struct with default values.
func NewBuildInfo() *BuildInfo {
	build := new(BuildInfo)
	build.Verbose = true
	build.NoCache = false
	return build
}

// LoadConfig loads the configuration from the path.
func LoadConfig(path string) (conf *Configuration) {
	conf = NewConfiguration()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		panic(err)
	}
	conf.ImageInfo.SetProperties()
	return conf
}

// ConfigExists checks if the input config exists.
func ConfigExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
