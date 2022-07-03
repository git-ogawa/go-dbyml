# go-dbyml

[![License](https://img.shields.io/github/license/git-ogawa/go-dbyml)](https://github.com/git-ogawa/go-dbyml/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/git-ogawa/go-dbyml)](https://github.com/git-ogawa/go-dbyml/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/git-ogawa/go-dbyml.svg)](https://pkg.go.dev/github.com/git-ogawa/go-dbyml)

[![build](https://github.com/git-ogawa/go-dbyml/actions/workflows/build.yml/badge.svg?branch=develop)](https://github.com/git-ogawa/go-dbyml/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/git-ogawa/go-dbyml)](https://goreportcard.com/report/github.com/git-ogawa/go-dbyml)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/86cff2a078d0455f945951c4474e9424)](https://www.codacy.com/gh/git-ogawa/go-dbyml/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=git-ogawa/go-dbyml&amp;utm_campaign=Badge_Grade)
[![git-ogawa](https://circleci.com/gh/git-ogawa/go-dbyml.svg?style=svg)](https://circleci.com/gh/git-ogawa/go-dbyml)

Go-dbyml is a CLI tool to build a docker image with build options loaded from yaml. Instead of running the `docker build` with many options, write options in config file, build your docker image with them. It helps you to manage build process more readable and flexible.

Go-dbyml is a substitute of [git-ogawa/dbyml](https://github.com/git-ogawa/dbyml) with golang.


# Table of Contents
- [Install](#install)
- [Preparation](#preparation)
- [Usage](#usage)
- [Build with buildkit](#build-with-buildkit)
- [Configuration](#configuration)
  - [Image](#image)
  - [Registry](#registry)
  - [Buildkit](#buildkit)
  - [Environment variables](#environment-variables)
  - [Examples](#examples)
- [Notes](#notes)


## Install
Install with `go install`.
```
go install github.com/git-ogawa/go-dbyml@latest
```

You can also manually download the binary from [release page](https://github.com/git-ogawa/go-dbyml/releases). Put it the directory in the $PATH.


## Preparation
To use go-dbyml, Docker Engine must be installed on host for build and run docker commands without root privileges (as non-root user) on client. Refer to [Manage Docker as a non-root user](https://docs.docker.com/engine/install/linux-postinstall/#manage-docker-as-a-non-root-user) or [Docker rootless mode](https://docs.docker.com/engine/security/rootless/) for non-root user setting.


## Usage
To build your image from Dockerfile, you must make `Dockerfile` and configuration file `dbyml.yml`. The file `dbyml.yml` is written in yaml syntax and contains the settings related to image build such as tags, build-args or label. To make the file, copy the contents from [examples/dbyml.yml](examples/dbyml.yml) or make by the following command.

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
REPOSITORY           TAG                 IMAGE ID       CREATED         SIZE
go-dbyml-sample      latest              cf55541823c7   5 hours ago     5.6MB
```


Go-dbyml has the following features for image build (these are the same as [git-ogawa/dbyml](https://github.com/git-ogawa/dbyml)).
- [Set build-args and labels in image](https://github.com/git-ogawa/dbyml#build-args-and-labels)
- [Push the image to a private registry](https://github.com/git-ogawa/dbyml#push-to-repository)

Most of features work by replacing `dbyml` with `go-dbyml`, but there are some differences (see [Notes](#notes)).


## Build with buildkit
Go-dbyml supports image build with [buildkit](https://github.com/moby/buildkit). The build with buildkit enables you to build multi-platform image and export and import build cache to external registry.


Additional settings are required in configuration file in order to enable buildkit. See [Buildkit](#buildkit).

# Configuration
The configuration about go-dbyml are managed by configuration file `dbyml.yml` written in yaml syntax. The Settings are automatically loaded from the file in current directory when run go-dbyml. Run `go-dbyml --init` to generate a config file from template.


The `dbyml.yml` consists of the following sections at the top level.

- Image
- Build
- Registry
- Buildkit

See the following description and [examples/dbyml.yml](examples/dbyml.yml) to know how to set these values.


## Image
The image section defines the basic settings about image to be built.

- `name`: The name of image.
- `tag`: The tag of image.
- `path`: Path to directory where Dockerfile exists.
- `dockerfile`: The filename of Dockerfile.
- `build_args`: The build-args used on build. These are passed as `docker build --build-arg [args]`.
- `label`: The labels used on build. These are passed as `docker build --label [labels]`.
- `docker_host`: URL to the Docker server.

## Registry
The registry section defines the registry information to which the built image is pushed.

- `enabled`: Set true to enable pushing image to a registry.
- `host`: Registry hostname including port.
- `insecure`: Set true to allow insecure server connections.
- `auth`: Credentials used when connect for auth-registry.

## Buildkit
The buildkit section defines the settings about buildkit. To build a image with buildkit, add the `buildkit` section in configuration file and set `enabled` to true.

### output
The output field sets output format of the image to be built. Only `Image` is supported now, which means the built image will be pushed the specified registry.

- `type`: Set `image`
- `name`: Set registry and image. e.g. `[registry]:[port]/[project]/[image]:[tag]`.
- `insecure`: Set true if push the image insecure registry such as insecure private registry. false otherwise.

### cache
The cache field sets import and export build cache. See [Cache](https://github.com/moby/buildkit#cache) for details.


#### export

- `type`: `inline` or `registry`.
- `value`: Set registry and image. e.g. `[registry]:[port]/[project]/[image]:[tag]` if type is registry.


#### import

- `type`: `registry`.
- `value`: Set registry and image. e.g. `[registry]:[port]/[project]/[image]:[tag]` if type is registry.


### platform
If you want to build a image supports multi-platform, Set the list of architectures to be supported in `platform` field.


## Environment variables
You can use environment variables in config file.

- To use the environment variable, set `${VARIABLE_NAME}`
- To use default value, set `${VARIABLE_NAME:-DEFAULT_VALUE}`.

```yaml
image:
  name: ${IMAGE_NAME}         # Set the value of ${IMAGE_NAME}.
  tag: ${TAG_NAME:-latest}    # Set latest if ${TAG_NAME} is undefined.
```

An error will be raised on build if the environment variable is undefined.

```yaml
image:
  name: ${UNDEFINED}        # >> ENV ${UNDEFINED} not defined.
```


## Examples
See [examples/dbyml.yml](examples/dbyml.yml) for an example of configuration.


## Notes
Compared to [git-ogawa/dbyml](https://github.com/git-ogawa/dbyml), there are some differences regarding the fields in the configuration file and features.


### ENV variables
The environment variable expression has been supported since v1.2.0.


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
Build with buildkit which works like buildx has been supported since v1.1.0.


### TLS section
The connection for docker host using TLS (HTTPS) dose not be supported yet.
