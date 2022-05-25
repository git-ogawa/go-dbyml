# go-dbyml

[![License](https://img.shields.io/github/license/git-ogawa/go-dbyml)](https://github.com/git-ogawa/go-dbyml/blob/main/LICENSE)
[![build](https://github.com/git-ogawa/go-dbyml/actions/workflows/build.yml/badge.svg?branch=develop)](https://github.com/git-ogawa/go-dbyml/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/git-ogawa/go-dbyml)](https://goreportcard.com/report/github.com/git-ogawa/go-dbyml)

Go-dbyml is a CLI tool to build a docker image with build options loaded from yaml. Instead of running the `docker build` with many options, write options in config file, build your docker image with them. It helps you to manage build process more readable and flexible.


Go-dbyml is a substitute of [git-ogawa/dbyml](https://github.com/git-ogawa/dbyml) with golang.


**This is a pre-alpha version now.**


## Install 
Install with `go install`.
```
go install github.com/git-ogawa/go-dbyml@latest
```

You can also manually download the binary from [release page](https://github.com/git-ogawa/go-dbyml/releases). Put it the directory in the $PATH.


## Usage
```
go-dbyml --init
```


```
go-dbyml
```




## Notes
Compared to [git-ogawa/dbyml](https://github.com/git-ogawa/dbyml), there are some differences regarding the fields in the configuration file and features.


### Build section
The following fields are currently enabled.

```yaml
build:
  target: ''
  no_cache: false
  verbose: true
```


### Registry section
The host and port are merge in host field, so set the format as "hostname:port" in the field.

```diff
registry:
    enabled: true
-    host: "myregistry" # Registry hostname or ip address 
-    port: "5000" # Registry port
+    host: "myregistry:5000"
```


### Buildx section
The multi-platform build using `buildx` does not support yet, so buildx section in config file does not work now. 
