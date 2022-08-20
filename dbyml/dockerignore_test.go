package dbyml

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Read the content from .dockerignore
func TestReadFile(t *testing.T) {
	pwd, _ := os.Getwd()
	root, _ := filepath.Abs("../")
	os.Chdir(root)

	excludes, err := ReadDockerignore("testdata/dockerfile_ignore")
	if err != nil {
		panic(err)
	}

	expected := []string{"ignore.txt", "ignore_dir", "*/*tmp*"}
	assert.Equal(t, expected, excludes)
	os.Chdir(pwd)
}

// Read the content from .dockerignore but it dose not exist.
func TestReadFileEmpty(t *testing.T) {
	pwd, _ := os.Getwd()
	root, _ := filepath.Abs("../")
	os.Chdir(root)

	excludes, err := ReadDockerignore("testdata/dockerfile_standard")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, []string{}, excludes)
	os.Chdir(pwd)
}
