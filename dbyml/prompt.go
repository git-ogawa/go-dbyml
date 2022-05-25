package dbyml

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
)

func CommonValidate(in string) error {
	if len(in) > 255 {
		return errors.New("word length must be less than 255")
	} else {
		return nil
	}
}

func NameValidate(in string) error {
	if len(in) == 0 || strings.Contains(in, ":") {
		return errors.New("image name cannot be empty or contain `:`")
	} else {
		return nil
	}
}

func TagValidate(in string) error {
	if strings.Contains(in, ":") {
		return errors.New("tag cannot contain `:`")
	} else {
		return nil
	}
}

func AnswerValidate(in string) error {
	r := regexp.MustCompile(in)
	if r.MatchString("y|yes|n|no") && in != "" {
		return nil
	} else {
		return errors.New("the answer must be yes or no (y/n)")
	}
}

func MapValidate(in string) error {
	arr := strings.Split(in, ":")
	if len(arr) == 2 && arr[0] != "" && arr[1] != "" {
		return nil
	} else {
		return errors.New("the answer must be key:value")
	}
}

// SelectPromptType determines which method is used to generate config file.
func SelectPromptType() (ret string, err error) {
	label := "Choose generation type."
	items := []string{"Interactive", "Non-interactive"}

	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, ret, err = prompt.Run()
	return ret, err
}

func GetString(validate promptui.ValidateFunc, label string, def string) (result string) {
	var msg string
	if def != "" {
		msg = fmt.Sprintf("%v Default: %v", label, def)
	} else {
		msg = fmt.Sprintf("%v", label)
	}

	prompt := promptui.Prompt{
		Label:    msg,
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	return result
}

func GetMap(validate promptui.ValidateFunc, label string, morePrompt string) map[string]string {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}
	result := map[string]string{}
	for {
		ret, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
		}
		k, v := splitString(ret)
		result[k] = v
		if !IsContinue(AnswerValidate, morePrompt) {
			break
		}
	}
	return result
}

func GetMapPointer(validate promptui.ValidateFunc, label string, morePrompt string) map[string]*string {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}
	result := map[string]*string{}
	for {
		ret, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
		}
		k, v := splitStringPointer(ret)
		result[k] = v
		if !IsContinue(AnswerValidate, morePrompt) {
			break
		}
	}
	return result
}

func IsContinue(validate promptui.ValidateFunc, label string) bool {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("is continue %v\n", err)
		// return false
	}
	if result == "y" || result == "yes" {
		return true
	} else {
		return false
	}
}

// GetUserInput gets user input interactively, generates config file from the template.
func GetUserInput() *Configuration {
	config := NewConfiguration()
	ret, _ := SelectPromptType()
	if ret == "Non-interactive" {
		return config
	}

	fmt.Println("Type values you want to set then press Enter. To use default value, press Enter with no input.")

	config.ImageInfo.Basename = GetString(NameValidate, "Input an image basename", "")
	config.ImageInfo.Tag = GetString(TagValidate, "Input a image tag.", "latest")
	config.ImageInfo.Path = GetString(CommonValidate, "Input path to the directory where Dockerfile exists.", ".")
	if IsContinue(AnswerValidate, "Do you set build args?") {
		config.ImageInfo.BuildArgs = GetMapPointer(MapValidate, "Input build args as key:value", "Do you set more build args?")
	}
	if IsContinue(AnswerValidate, "Do you set label?") {
		config.ImageInfo.Labels = GetMap(MapValidate, "Input labels as key:value.", "Do you set more labels?")
	}
	if IsContinue(AnswerValidate, "Do you set registry information?") {
		config.RegistryInfo.Enabled = true
		config.RegistryInfo.Host = GetString(CommonValidate, "Input registry host:port.", "myregistry:5000")
		config.RegistryInfo.Project = GetString(CommonValidate, "Input project name.", "")
		username := GetString(CommonValidate, "Input username to login registry.", "''")
		password := GetString(CommonValidate, "Input password to login registry.", "''")
		config.RegistryInfo.Auth = map[string]string{"username": username, "password": password}
	}
	return config
}
