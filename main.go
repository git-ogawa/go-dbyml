package main

import (
	"github.com/git-ogawa/go-dbyml/dbyml"
)

func main() {
	cli, exec := dbyml.GetArgs()
	if exec {
		cli.Parse()
	}
}
