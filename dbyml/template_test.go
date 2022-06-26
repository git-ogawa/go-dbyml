package dbyml

import (
	"os"
	"testing"
)

func TestConfigTemplate(t *testing.T) {
	config := NewConfiguration()
	MakeTemplate(config)
	_, err := os.Stat("dbyml.yml")
	if err != nil {
		panic(err)
	}
	os.Remove("dbyml.yml")
}

func TestBuildkitToml(t *testing.T) {
	config := NewConfiguration()
	MakeBuildkitToml(&config.RegistryInfo)
	_, err := os.Stat("buildkitd.toml")
	if err != nil {
		panic(err)
	}
	os.Remove("buildkitd.toml")
}
