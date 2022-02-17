package metabase

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePermissionGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePermissionGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Group Name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"group_id": {
				Description: "Group Id",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func dataSourcePermissionGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Get("name").(string)
	d.SetId(id)
	return resourcePermissionGroupRead(ctx, d, meta)
}
