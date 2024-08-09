# Tubi fork - Terraform Runscope Provider

The Runscope provider is used to interact with the resources
supported by [Runscope](https://runscope.com/).

This is a Tubi fork of the [sport24ru/runscope](https://github.com/sport24ru/terraform-provider-runscope) provider to ease `darwin_arm64` builds.

## Build Instructions
* goreleaser is a dependency, you can install it with `$ brew install goreleaser`
* `$ goreleaser release --snapshot` will build a snapshot locally
* `$ goreleaser release` will build and version against the commits tag


## Usage

Read the [documentation on Terraform Registry site](https://registry.terraform.io/providers/sport24ru/runscope/latest/docs).
