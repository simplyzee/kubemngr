# kubemngr

[![Releases](https://img.shields.io/github/release/zee-ahmed/kubemngr.svg?style=flat-square)](https://github.com/zee-ahmed/kubemngr/releases/latest) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](/LICENSE) <!-- [![SayThanks.io]()] --> [![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)[![Travis](https://img.shields.io/travis/zee-ahmed/kubemngr/master.svg?style=flat-square)](https://travis-ci.org/zee-ahmed/kubemngr)

## Introduction

kubemngr is a cli tool to manage kubectl binaries for developers who work with different versions of Kubernetes clusters within their environments. This way the developer can keep in sync with the kubernetes cluster version. This tool was written in Golang using [Cobra](https://github.com/spf13/cobra)

99% of the case - developers will not have a problem using `kubectl` across different versions of k8s clusters with api requests for getting pods, ingress and services etc. However, there is the small chance that some requests will return with:

```bash
> kubectl describe ing <name>
Error from server (NotFound): the server could not find the requested resource
```

This tool was written based on this experience and for learning opportunities/experiences.

## Install

via Go:
```
go get -u github.com/zee-ahmed/kubemngr
```
It can also be installed by downloading the binary from the Github release page [Github Releases](https://github.com/zee-ahmed/kubemngr/releases)

## Usage
```bash
> kubemngr --help
This tool is to help developers run different versions of kubectl within their workspace and to support working
with different versions of Kubernetes clusters.

Usage:
  kubemngr [command]

Available Commands:
  help        Help about any command
  install     A tool manage different kubectl versions inside a workspace.
  list        List installed kubectl binary versions. For available versions, see --remote
  remove      Remove a kubectl version from machine
  use         Use a specific version of one of the downloaded kubectl binaries
  version     Show the kubemngr client version

Flags:
  -h, --help     help for kubemngr
  -t, --toggle   Help message for toggle

Use "kubemngr [command] --help" for more information about a command.
```

## Contributing

Please raise an issue or pull request if you have any issues, questions or features.
