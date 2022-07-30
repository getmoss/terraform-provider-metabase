package metabase

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"terraform-provider-metabase/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCollection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCollectionCreate,
		ReadContext:   resourceCollectionRead,
		UpdateContext: resourceCollectionUpdate,
		DeleteContext: resourceCollectionArchive,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Collection name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"color": {
				Type:     schema.TypeString,
				Required: true,
			},
			"archived": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceCollectionUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id, _ := strconv.Atoi(d.Id())
	name := d.Get("name").(string)
	color := d.Get("color").(string)
	archived := d.Get("archived").(bool)

	col := client.Collection{
		Id:       id,
		Name:     name,
		Color:    color,
		Archived: archived,
	}

	// Update the collection
	updated, err := c.UpdateCollection(col)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error updating Collection '%v'", col),
			Detail:   "Could not update Collection, unexpected error: " + err.Error(),
		})
		return diags
	}

	if err := d.Set("name", updated.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("color", updated.Color); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("archived", updated.Archived); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCollectionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	color := d.Get("color").(string)
	archived := d.Get("archived").(bool)
	col := client.Collection{
		Name:     name,
		Color:    color,
		Archived: archived,
	}

	// Create the collection
	created, err := c.CreateCollection(col)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error creating collection '%s'", name),
			Detail:   "Could not create collection, unexpected error: " + err.Error(),
		})
		return diags
	}

	d.SetId(fmt.Sprint(created.Id))
	if err := d.Set("name", created.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("color", created.Color); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("archived", created.Archived); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCollectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	var diags diag.Diagnostics
	var id string
	if d.Id() == "root" {
		id = "root"
	} else {
		id = fmt.Sprint(d.Id())
	}

	log.Printf("[INFO] Finding collection by id '%s'", id)

	col, err := c.GetCollection(id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error reading collection with id '%s'", id),
			Detail:   "Could not read collection: " + err.Error(),
		})
		return diags
	}

	d.SetId(fmt.Sprintf("%v", col.Id))
	if err := d.Set("name", col.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("color", col.Color); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("archived", col.Archived); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// Metabase API does not implement a Delete for collection, the closer action it is to archive it.
func resourceCollectionArchive(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id, _ := strconv.Atoi(d.Id())

	archived := client.Collection{
		Id:       id,
		Name:     d.Get("name").(string),
		Color:    d.Get("color").(string),
		Archived: true,
	}

	a, err := c.UpdateCollection(archived)
	if err != nil {
		return diag.FromErr(err)
	}

	if !a.Archived {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error archived collection id '%s'", d.Id()),
			Detail:   "It was not possible to archive (delete) the collection: " + err.Error(),
		})
		return diags
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
