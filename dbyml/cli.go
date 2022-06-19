// Dbyml is a CLI tool to build a docker image with the arguments loaded from configs in yaml.
//
// Usage
//
// The following command will generate a configuration file.
// 	dbyml --init
//
// The configuration file where the about image build is written.
//
// 	dbyml
//
// The options on image build can be written in the config.
// See https://github.com/git-ogawa/dbyml how to edit contents in config file.
package dbyml

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

// CLIoptions defines cli options.
type CLIoptions struct {
	// Path to config file.
	Config string

	// Whether to generate config file.
	Init bool
}

// GetArgs gets cli options from user inputs.
func GetArgs() (CLIoptions, bool) {
	desc := "Dbyml is a CLI tool to build a docker image with the arguments loaded from configs in yaml.\n\n"
	desc += "Passing the config file where the arguments are listed to build the image from your dockerfile,\n"
	desc += "push it to the docker registry.\n\n"
	desc += "To make sample config file, run the following command.\n"
	desc += "\n"
	desc += "$ dbyml --init\n"

	parser := argparse.NewParser("dbyml", desc)
	parser.HelpFunc = func(c *argparse.Command, msg interface{}) string {
		var help string
		help += fmt.Sprintln(c.GetDescription())
		help += "Optional arguments:\n"
		for _, arg := range c.GetArgs() {
			if arg.GetOpts() != nil {
				sopt := arg.GetSname()
				var prefix string
				var suffix string
				if sopt != "" {
					prefix = "-"
					suffix = ","
				} else {
					prefix = ""
					suffix = ""
				}
				sname := fmt.Sprintf("%v%s%v", prefix, sopt, suffix)
				lname := fmt.Sprintf("--%-15s", arg.GetLname())
				helpMsg := arg.GetOpts().Help
				help += fmt.Sprintf("  %3s %s %s\n", sname, lname, helpMsg)
			} else {
				help += fmt.Sprintf("Sname: %s, Lname: %s\n", arg.GetSname(), arg.GetLname())
			}
		}
		return help
	}

	Config := parser.String("c", "config", &argparse.Options{Help: "Path to config file."})
	Init := parser.Flag("", "init", &argparse.Options{Help: "Generate config."})
	Version := parser.Flag("v", "version", &argparse.Options{Help: "Show version."})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	if *Version {
		ShowVersion()
		return CLIoptions{}, false
	}

	return CLIoptions{*Config, *Init}, true
}

// Parse checks the input options, run actions according to the options.
func (options *CLIoptions) Parse() {
	if options.Init {
		config := NewConfiguration()
		MakeTemplate(config)
		return
	}
	if options.Config != "" {
		if exist := ConfigExists(options.Config); exist {
			ExecBuild(options.Config)
		} else {
			fmt.Printf("%v not found. Check the file exists.\n", options.Config)
		}
	} else {
		if exist := ConfigExists("dbyml.yml"); exist {
			ExecBuild("dbyml.yml")
		} else {
			msg := "Config file not found in the current directory.\nRun the following commands to generate config file."
			fmt.Println(msg)
			fmt.Println()
			fmt.Printf("%-20v Generate config file interactively or non-interactively.\n", "dbyml --init")
		}
	}
}

// ExecBuild run the build sequence.
func ExecBuild(path string) {
	config := LoadConfig(path)
	config.ImageInfo.Registry = config.RegistryInfo
	config.ImageInfo.BuildInfo = config.BuildInfo
	if config.BuildInfo.Verbose {
		config.ShowConfig()
	}

	fmt.Println()
	PrintCenter("Build start", 30, "-")
	fmt.Println()
	err := config.ImageInfo.Build()
	PrintCenter("Build finish", 30, "-")
	fmt.Println()
	if err != nil {
		fmt.Printf("Error has occurred during build: %v\n", err)
		fmt.Println("\x1b[31mBuild Failed\x1b[0m")
		os.Exit(1)
	}

	fmt.Printf("Image %v successfully built.\n", config.ImageInfo.ImageName)

	if config.RegistryInfo.Enabled {
		fmt.Println()
		PrintCenter("Push start", 30, "-")
		fmt.Println()
		err = config.ImageInfo.Push()
		fmt.Println()
		PrintCenter("Push finish", 30, "-")

		fmt.Println()
		if err != nil {
			fmt.Printf("Error has occurred during push: %v\n", err)
			fmt.Println("\x1b[31mPush Failed\x1b[0m")
			os.Exit(1)
		}

		fmt.Printf("Image %v successfully pushed.\n", config.ImageInfo.FullName)
	}
}
