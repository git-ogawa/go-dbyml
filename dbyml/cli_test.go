package dbyml

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigNotFound(t *testing.T) {
	pwd, _ := os.Getwd()
	root, _ := filepath.Abs("../")
	os.Chdir(root)

	options := CLIoptions{"", false}
	// options.Parse()
	stdout := extractStdout(t, options.Parse)
	expected := "Config file not found in the current directory.\n"
	expected += "Run the following commands to generate config file.\n\n"
	expected += "dbyml --init         Generate config file interactively or non-interactively."
	assert.Equal(t, expected, stdout)

	options = CLIoptions{"notexists.yml", false}
	options.Parse()
	stdout = extractStdout(t, options.Parse)
	expected = "notexists.yml not found. Check the file exists."
	assert.Equal(t, expected, stdout)

	os.Chdir(pwd)
}

// func TestCLIShowHelp(t *testing.T) {
// 	org := os.Args
// 	testArgs := []string{"dbyml", "-h"}
// 	os.Args = testArgs

// 	options := GetArgs()
// 	stdout := extractStdout(t, options.Parse)

// 	expected := "Dbyml is a CLI tool to build your docker image with the arguments loaded from configs in yaml.\n\n"
// 	expected += "Passing the config file where the arguments are listed to build the image from your dockerfile,\n"
// 	expected += "push it to the docker registry.\n\n"
// 	expected += "To make sample config file, run the following command.\n\n"
// 	expected += "$ dbyml --init\n\n"
// 	expected += "Optional arguments:\n"
// 	expected += "  -h, --help            Print help information\n"
// 	expected += "  -c, --config          Path to config file.\n"
// 	expected += "	  --init            Generate config.\n"
// 	expected += "  -v, --version         Show version.\n"

// 	assert.Equal(t, expected, stdout)

// 	os.Args = org
// }

func TestCLIBuild(t *testing.T) {
	pwd, _ := os.Getwd()
	root, _ := filepath.Abs("../")
	os.Chdir(root)

	// Set args
	org := os.Args
	path, _ := filepath.Abs("testdata/dbyml.yml")
	testArgs := []string{"dbyml", "-c", path}
	os.Args = testArgs

	options := GetArgs()
	options.Parse()

	os.Args = org
	os.Chdir(pwd)
}

func extractStdout(t *testing.T, fnc func()) string {
	t.Helper()

	orgStdout := os.Stdout
	defer func() {
		os.Stdout = orgStdout
	}()
	r, w, _ := os.Pipe()
	os.Stdout = w
	fnc()
	w.Close()
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to read buf: %v", err)
	}
	return strings.TrimRight(buf.String(), "\n")
}
