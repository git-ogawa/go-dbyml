package dbyml

import (
	"fmt"
	// "io/ioutil"
	"os"
	"text/template"
	// "gopkg.in/yaml.v2"
)

func MakeTemplate(config *Configuration) {
	templateFile := "templates/template.yml"

	tmpl := template.Must(template.ParseFiles(templateFile))

	file, _ := os.Create("dbyml.yml")
	err := tmpl.Execute(file, config)
	if err != nil {
		panic(err)
	}
	fmt.Println("Create dbyml.yml. Check the contents and edit it according to your docker image.")
}
