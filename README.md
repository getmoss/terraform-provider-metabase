# Terraform Metabase

Terraform provider for [Metabase](https://metabase.com/). It is built using the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) & [Metabase APIs](https://metabase.com/docs/latest/api-documentation.html#permissions).


This is the blog that I followed. It explains the concepts in detail:
https://learn.hashicorp.com/tutorials/terraform/provider-use?in=terraform/providers

### Developing a Provider
There are a bunch of structs in the SDK that you have to provide implementations for. `terraform` uses these when you plug-in your provider.
1. Create an instance of [Provider](https://pkg.go.dev/github.com/hashicorp/terraform/helper/schema#Provider). Under this, you have to override the following 4 members => 
    * `ConfigureContextFunc` to configure the provider.
    * `Schema` returns the supported fields.
    * `ResourcesMap` returns the supported resources.
    * `DataSourcesMap` returns the supported data sources.
2. To add a `resource`, create an instance of [Resource](https://pkg.go.dev/github.com/hashicorp/terraform/helper/schema#Resource). Override the necessary members =>
    * `Schema` returns the supported fields.
    * Provide mplementations of 4 `CRUD` methods + `ImportState` method as needed.
3. To add a `data source`, implement [Resource](https://pkg.go.dev/github.com/hashicorp/terraform/helper/schema#Resource), same as `resource` above. For a data type, 2 methods are required:
    * `Schema` returns the supported fields.
    * `ReadContext` that reads the resource from the API. This can delegate to the `ReadContext` of the `resource`.

#### To test
* `alias tf=terraform`
* Update `~/.terraformrc` with the following:
```tf
  dev_overrides {
      "github.com/getmoss/metabase" = "/Users/raj/go/bin"
  }
```
* Go to `examples` directory
* Create `local.tfvars` with your Metabase credentials
* Plan with `TF_LOG=debug tfp -var-file=local.tfvars`
* Apply with `TF_LOG=debug tf apply`


__________________
## Content from template =>
_This template repository is built on the [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk). The template repository built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) can be found at [terraform-provider-scaffolding-framework](https://github.com/hashicorp/terraform-provider-scaffolding-framework). See [Which SDK Should I Use?](https://www.terraform.io/docs/plugin/which-sdk.html) in the Terraform documentation for additional information._

This repository is a *template* for a [Terraform](https://www.terraform.io) provider. It is intended as a starting point for creating Terraform providers, containing:

 - A resource, and a data source (`internal/provider/`),
 - Examples (`examples/`) and generated documentation (`docs/`),
 - Miscellaneous meta files.
 
These files contain boilerplate code that you will need to edit to create your own Terraform provider. A full guide to creating Terraform providers can be found at [Writing Custom Providers](https://www.terraform.io/docs/extend/writing-custom-providers.html).

Please see the [GitHub template repository documentation](https://help.github.com/en/github/creating-cloning-and-archiving-repositories/creating-a-repository-from-a-template) for how to create a new repository from this template on GitHub.

Once you've written your provider, you'll want to [publish it on the Terraform Registry](https://www.terraform.io/docs/registry/providers/publishing.html) so that others can use it.


## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.15

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
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

```sh
$ make testacc
```
