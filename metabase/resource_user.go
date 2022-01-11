package metabase

import (
	"context"
	"fmt"
	"strconv"
	"terraform-provider-metabase/client"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		DeleteContext: resourceUserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceUserCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	email := d.Get("email").(string)
	userId := d.Get("user_id").(int)
	u := client.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	if userId != 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "user_id is not empty",
			Detail:   "user_id cannot be specified",
		})
		return diags
	}

	// Create the user
	created, err := c.CreateUser(u)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error creating User '%s'", email),
			Detail:   "Could not create User, unexpected error: " + err.Error(),
		})
		return diags
	}

	d.SetId(strconv.Itoa(created.Id))
	if err := d.Set("first_name", created.FirstName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_name", created.LastName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email", created.Email); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("user_id", created.Id); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	email := d.Get("email").(string)
	id := d.Get("user_id").(int)

	var diags diag.Diagnostics

	var user client.User
	// Get current values from API by Id if present
	if id != 0 {
		var err error
		user, err = c.GetUser(id)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading",
				Detail:   "Could not read user=" + strconv.Itoa(id) + ": " + err.Error(),
			})
			return diags
		}
	}

	if (user == client.User{}) {
		// fetch all users and find the one with the given email
		users, err := c.GetUsers()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading",
				Detail:   "Could not read users: " + err.Error(),
			})
			return diags
		}
		for _, i := range users.Data {
			if user.Email == email {
				user = i
				break
			}
		}
	}

	d.SetId(strconv.Itoa(user.Id))
	if err := d.Set("user_id", user.Id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("first_name", user.FirstName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_name", user.LastName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email", user.LastName); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceUserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Get("group_id").(int)

	_, err := c.DeleteUser(id)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
