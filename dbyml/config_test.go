package dbyml

import (
	"encoding/base64"
	"encoding/json"

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
