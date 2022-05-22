package dbyml

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNonInteractiveConfig(t *testing.T) {
	pwd, _ := os.Getwd()
	root, _ := filepath.Abs("../")
	os.Chdir(root)

	options := CLIoptions{"", true}
	// options.Parse()
	options.Parse()
	// stdout := extractStdout(t, options.Parse)
	_, err := os.Stat("dbyml.yml")
	if err == nil {
		os.Remove("dbyml.yml")
	}

	os.Chdir(pwd)
}
