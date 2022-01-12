package metabase

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"terraform-provider-metabase/client"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceUserUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	email := d.Get("email").(string)
	userId := d.Get("user_id").(int)
	u := client.User{
		Id:        userId,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	// Update the user
	updated, err := c.UpdateUser(u)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error updating User '%s'", email),
			Detail:   "Could not update User, unexpected error: " + err.Error(),
		})
		return diags
	}

	d.SetId(strconv.Itoa(updated.Id))
	if err := d.Set("first_name", updated.FirstName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_name", updated.LastName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email", updated.Email); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("user_id", updated.Id); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceUserCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	email := d.Get("email").(string)
	u := client.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
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

	id := getIdForRead(d)

	email := d.Get("email").(string)

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
		log.Printf("[INFO] Finding user by email '%s'", email)
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
			if i.Email == email {
				user = i
				break
			}
		}

		if (user == client.User{}) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading user",
				Detail:   fmt.Sprintf("Could not read user with email '%s'", email),
			})
			return diags
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
	if err := d.Set("email", user.Email); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func getIdForRead(d *schema.ResourceData) int {
	id, err := strconv.Atoi(d.Id()) // Used for import
	if err != nil || id == 0 {
		return d.Get("user_id").(int) // Used for data source
	}
	return id
}

func resourceUserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Get("user_id").(int)

	_, err := c.DeleteUser(id)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
