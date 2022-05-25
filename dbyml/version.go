package dbyml

import (
	"fmt"
)

// Module name and version
const (
	moduleName = "go-dbyml"
	version    = "v0.0.1"
)

// ShowVersion shows module version.
func ShowVersion() {
	fmt.Printf("%v %v\n", moduleName, version)
}
