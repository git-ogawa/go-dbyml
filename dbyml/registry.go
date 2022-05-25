package dbyml

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/docker/docker/api/types"
)

// RegistryInfo defines a registry settings where a built image will be pushed.
type RegistryInfo struct {
	// Whether to enabled pushing to a registry
	Enabled bool `yaml:"enabled"`

	// Registry host such as `myregistry.com:5000`
	Host string `yaml:"host"`

	// Project name
	Project string `yaml:"project"`

	// credentials settings to a registry
	Auth map[string]string `yaml:"auth"`
}

// NewRegistryInfo creates a new RegistryInfo struct with default values.
func NewRegistryInfo() *RegistryInfo {
	registry := new(RegistryInfo)
	registry.Enabled = false
	return registry
}

// ShowProperties shows the current registry settings to stdout.
func (registry RegistryInfo) ShowProperties() {
	rv := reflect.ValueOf(registry)
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		value := rv.FieldByName(field.Name)
		if field.Name == "Auth" {
			fmt.Printf("%-30v: *********\n", field.Name)
		} else {
			fmt.Printf("%-30v: %v\n", field.Name, value)
		}
	}
}

// BasicAuth returns base64 the encoded credentials for the registry.
func (registry *RegistryInfo) BasicAuth() string {
	return GetAuthBase64(registry.Auth["username"], registry.Auth["password"])
}

// GetAuthBase64 encodes credentials for the registry with base64.
func GetAuthBase64(username string, password string) string {
	auth := types.AuthConfig{
		Username: username,
		Password: password,
	}
	authBytes, _ := json.Marshal(auth)
	return base64.URLEncoding.EncodeToString(authBytes)
}
