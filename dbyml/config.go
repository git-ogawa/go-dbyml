package dbyml

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

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
	rep, err := parseEnv(string(data))
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal([]byte(rep), &conf); err != nil {
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

func parseEnv(data string) (string, error) {
	matches := regexp.MustCompile(`\${.*}`).FindAllString(data, -1)
	res, err := getEnvs(matches)
	if err != nil {
		return "", err
	}
	for k, v := range res {
		data = regexp.MustCompile(fmt.Sprintf(`\%v`, k)).ReplaceAllString(data, v)
	}
	return data, nil
}

// getEnvs replaces the specified variables in environment variables.
// If the environment variable is not defined and the default value is not set, returns error.
func getEnvs(envs []string) (map[string]string, error) {
	var target, def string
	res := map[string]string{}
	for _, env := range envs {
		e := regexp.MustCompile("[${}]").ReplaceAllString(env, "")

		// Split env and its default value.
		arr := regexp.MustCompile("(:-)").Split(e, -1)
		if len(arr) == 1 {
			target = arr[0]
			def = ""
		} else if len(arr) == 2 {
			target = arr[0]
			def = regexp.MustCompile("[${}]").ReplaceAllString(arr[1], "")
		}

		// Get the env value.
		rep := os.Getenv(target)
		if rep == "" {
			if def == "" {
				return res, fmt.Errorf(fmt.Sprintf("ENV %v not defined.", env))
			}
			rep = def
		}
		res[env] = rep
	}
	return res, nil
}
