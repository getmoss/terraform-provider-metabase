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
      "github.com/getmoss/metabase" = "<GOPATH>/bin"
  }
```
* Go to `examples` directory
* Create `local.tfvars` with your Metabase credentials
* Plan with `TF_LOG=debug tfp -var-file=local.tfvars`
* Apply with `TF_LOG=debug tf apply`

### References
https://learn.hashicorp.com/tutorials/terraform/provider-use?in=terraform/providers