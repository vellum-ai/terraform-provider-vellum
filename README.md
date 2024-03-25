# Vellum Terraform Provider

This repository is the [Terraform](https://www.terraform.io) provider. 

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## Code Generation
We use code generation where possible to auto-generate as much provider code as possible from an OpenAPI spec.
Where not possible, we generate the initial scaffolding and boilerplate and then hand-write the rest.

### Generating from OpenAPI Spec

Step 1) Install these two generator tools:

```shell
go install github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi@latest
go install github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework@latest
```

Step 2) If needed, update the OpenAPI spec found at `internal/provider/specs/openapi.yaml`.

Step 3) If defining a new resource, specify the mappings in `internal/provider/specs/generator_config.yml`.

Step 4) Run the following command to generate the provider code:

```shell
make generate-all
```

### Scaffolding a New Resource
Code generation is quite limited today (see [known limitations](https://github.com/hashicorp/terraform-plugin-codegen-openapi/blob/main/DESIGN.md#known-limitations).
If you need to define a resource and code generation doesn't work, you can at least scaffold its initial boilerplate
and then hand-write the rest using the following command:

```shell
make scaffold-resource name=my_resource
```