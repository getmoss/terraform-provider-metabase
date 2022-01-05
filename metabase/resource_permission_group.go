package metabase

import (
	"context"
	"errors"
	"strconv"
	"terraform-provider-metabase/client"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePermissionGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePermissionGroupCreate,
		ReadContext:   resourcePermissionGroupRead,
		DeleteContext: resourcePermissionGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourcePermissionGroupImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// "id": {
			// 	Type:     schema.TypeInt,
			// 	Computed: true,
			// },
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourcePermissionGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourcePermissionGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	name := d.Get("name").(string)
	id := d.Get("id").(int)

	var diags diag.Diagnostics

	var pg client.PermissionGroup
	// Get current values from API by Id if present
	if id != 0 {
		var err error
		pg, err = c.GetPermissionGroup(id)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading",
				Detail:   "Could not read permissionGroup=" + string(id) + ": " + err.Error(),
			})
			return diags
		}
	}

	if (pg == client.PermissionGroup{}) {
		// fetch all permission groups and find the one with the given name
		pgs, err := c.GetPermissionGroups()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading",
				Detail:   "Could not read permissionGroups: " + err.Error(),
			})
			return diags
		}
		for _, i := range pgs {
			if pg.Name == name {
				pg.Id = i.Id
				pg.Name = i.Name
				break
			}
		}
	}

	d.SetId(strconv.Itoa(pg.Id))
	if err := d.Set("id", pg.Id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", pg.Name); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourcePermissionGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourcePermissionGroupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if diag := resourcePermissionGroupRead(context.Background(), d, meta); diag.HasError() {
		return nil, errors.New("error resourcePermissionGroupImport")
	}
	return []*schema.ResourceData{d}, nil
}
