package dbyml

import (
	"fmt"
)

// Module name and version
const (
	moduleName = "go-dbyml"
	version    = "v1.0.0"
)

// ShowVersion shows module version.
func ShowVersion() {
	fmt.Printf("%v %v\n", moduleName, version)
}
