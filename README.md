# go-dbyml

[![License](https://img.shields.io/github/license/git-ogawa/go-dbyml)](https://github.com/git-ogawa/go-dbyml/blob/main/LICENSE)
[![build](https://github.com/git-ogawa/go-dbyml/actions/workflows/build.yml/badge.svg?branch=develop)](https://github.com/git-ogawa/go-dbyml/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/git-ogawa/go-dbyml)](https://goreportcard.com/report/github.com/git-ogawa/go-dbyml)
[![Release](https://img.shields.io/github/v/release/git-ogawa/go-dbyml)](https://github.com/git-ogawa/go-dbyml/releases)


Go-dbyml is a CLI tool to build a docker image with build options loaded from yaml. Instead of running the `docker build` with many options, write options in config file, build your docker image with them. It helps you to manage build process more readable and flexible.

Go-dbyml is a substitute of [git-ogawa/dbyml](https://github.com/git-ogawa/dbyml) with golang.


## Install 
Install with `go install`.
```
go install github.com/git-ogawa/go-dbyml@latest
```

You can also manually download the binary from [release page](https://github.com/git-ogawa/go-dbyml/releases). Put it the directory in the $PATH.


## Preparation
To use go-dbyml, Docker Engine must be installed on host for build and run docker commands without root privileges (as non-root user) on client. Refer to [Manage Docker as a non-root user](https://docs.docker.com/engine/install/linux-postinstall/#manage-docker-as-a-non-root-user) or [Docker rootless mode](https://docs.docker.com/engine/security/rootless/) for non-root user setting.


## Usage
To build your image from Dockerfile, you must make `Dockerfile` and configuration file `dbyml.yml`. The file `dbyml.yml` is written in yaml syntax and contains the settings related to image build such as tags, build-args or label. To make the file, copy the contents from [examples/dbyml.yml](https://github.com/git-ogawa/go-dbyml/blob/develop/examples/dbyml.yml) or make by the following command.

```
$ go-dbyml --init
```

You can edit the contents of the generated file in order to change settings about image build (name, tag, label or whether to push the image to a private registry). After that,run the `go-dbyml` in the same directory where the config file exists to build the image with the settings.
```
$ go-dbyml
```

After successfully building, the docker image will be created.
```
$ docker images
REPOSITORY                                                          TAG                 IMAGE ID       CREATED         SIZE
go-dbyml-sample                                                     latest              cf55541823c7   5 hours ago     5.6MB
```


Go-dbyml has the following features for image build (these are the same as [git-ogawa/dbyml](https://github.com/git-ogawa/dbyml)).
- [Set build-args and labels in image](https://github.com/git-ogawa/dbyml#build-args-and-labels)
- [Push the image to a private registry](https://github.com/git-ogawa/dbyml#push-to-repository)

Most of features work by replacing `dbyml` with `go-dbyml`, but there are some differences (see below).


## Notes
Compared to [git-ogawa/dbyml](https://github.com/git-ogawa/dbyml), there are some differences regarding the fields in the configuration file and features.


### ENV variables
The environment variable expression does not be supported yet.


### Build section
The following fields are currently enabled.

```yaml
build:
  target: ''
  no_cache: false
  verbose: true
```


### Registry section
The host and port has been merged in host field, so set the format as "hostname:port" in the field.

```diff
registry:
    enabled: true
-    host: "myregistry" # Registry hostname or ip address 
-    port: "5000" # Registry port
+    host: "myregistry:5000"
```


### Buildx section
The multi-platform build using `buildx` does not be supported yet, so buildx section in config file does not work now. 


### TLS section
The connection for docker host using TLS (HTTPS) dose not be supported yet.
