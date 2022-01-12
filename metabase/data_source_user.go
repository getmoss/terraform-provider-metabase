package metabase

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Description: "User Id. At least one of `user_id` or `email` must be specified.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"email": {
				Description: "User email. At least one of `user_id` or `email` must be specified.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceUserRead(ctx, d, meta)
}
