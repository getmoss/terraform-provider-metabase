# Terraform Metabase

Terraform provider for [Metabase](https://metabase.com/). It is built using the [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) & [Metabase APIs](https://metabase.com/docs/latest/api-documentation.html#permissions).

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

**Ensure the host value doesn't contain a slash '/' at the end**

## Developing the Provider
Checkout [notes](./development.md)

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To compile for Silicon M1 run `GOOS=darwin GOARCH=arm64 go build -o terraform-provider-metabase .`

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
