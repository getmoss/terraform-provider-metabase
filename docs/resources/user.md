---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "metabase_user Resource - terraform-provider-metabase"
subcategory: ""
description: |-
  
---

# metabase_user (Resource)



## Example Usage

```terraform
resource "metabase_user" "example" {
  email      = "john.doe@example.com"
  first_name = "John"
  last_name  = "Doe"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `email` (String)
- `first_name` (String)
- `last_name` (String)

### Read-Only

- `id` (String) The ID of this resource.
- `user_id` (Number)


