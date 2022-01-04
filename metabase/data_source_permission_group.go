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
			"id": {
				Description: "Group Id. At least one of `id` or `name` must be specified.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"name": {
				Description: "Group Name. At least one of `id` or `name` must be specified.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func dataSourcePermissionGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourcePermissionGroupRead(ctx, d, meta)
}
