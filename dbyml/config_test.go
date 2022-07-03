package dbyml

import (
	"encoding/base64"
	"encoding/json"
	"os"

	// "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {
	registry := NewRegistryInfo()
	expected := new(RegistryInfo)
	expected.Enabled = false

	assert.Equal(t, expected, registry)
}

func TestEncodeDecoder(t *testing.T) {
	registry := NewRegistryInfo()
	auth := map[string]string{"username": "docker", "password": "docker"}
	registry.Auth = auth
	encode := registry.BasicAuth()
	b, _ := base64.URLEncoding.DecodeString(encode)
	var decode map[string]string
	err := json.Unmarshal(b, &decode)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, auth, decode)
}

func TestParseEnv(t *testing.T) {
	os.Setenv("TEST", "en_US.UTF-8")
	os.Unsetenv("DUMMY")

	data := `
	test1: ${TEST}
	test2: ${DUMMY:-default}`
	res, err := parseEnv(data)
	if err != nil {
		panic(err)
	}

	expected := `
	test1: en_US.UTF-8
	test2: default`
	assert.Equal(t, res, expected)
	os.Unsetenv("TEST")

	_, err = parseEnv(data)
	if err != nil {
		assert.Equal(t, err.Error(), "ENV ${TEST} not defined.")
	}
}
