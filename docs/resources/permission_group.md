---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "metabase_permission_group Resource - terraform-provider-metabase"
subcategory: ""
description: |-
  
---

# metabase_permission_group (Resource)



## Example Usage

```terraform
resource "metabase_permission_group" "example" {
  name = "my-department"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)

### Read-Only

- `group_id` (Number)
- `id` (String) The ID of this resource.


