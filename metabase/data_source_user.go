package metabase

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Description: "Group Id. At least one of `group_id` or `name` must be specified.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"name": {
				Description: "Group Name. At least one of `group_id` or `name` must be specified.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourcePermissionGroupRead(ctx, d, meta)
}
