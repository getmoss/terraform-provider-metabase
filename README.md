# Terraform Metabase

Terraform provider for [Metabase](https://metabase.com/). It is built using the [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) & [Metabase APIs](https://metabase.com/docs/latest/api-documentation.html#permissions).


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
    * Provide implementations of 4 `CRUD` methods + `ImportState` method as needed.
3. To add a `data source`, implement [Resource](https://pkg.go.dev/github.com/hashicorp/terraform/helper/schema#Resource), same as `resource` above. For a data type, 2 methods are required:
    * `Schema` returns the supported fields.
    * `ReadContext` that reads the resource from the API. This can delegate to the `ReadContext` of the `resource`.
4. It's practical & a convention to create a client instead of directly handling HTTP in the provider. This helps with testing the client as well.

### Adding a new resource
1. Add User CRUD methods for the `client` with tests.
2. Add `resource` & `data` definitions to the provider using the client.
3. Manual testing with the `examples/resources.tf` file & code adjustment as necessary.

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
## Using the provider

Initialise the provider with `host`, `username` & `password`:

```tf
provider "metabase" {
  username = var.username
  password = var.password
  host     = var.host
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
