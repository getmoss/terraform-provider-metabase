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
			"parent_id": {
				Description: "Parent collection id",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Description: "Collection name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"default_access": {
				Description: "Default access for all users",
				Type:        schema.TypeString,
				Default:     "none",
				Optional:    true,
			},
			"permissions": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"color": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "#31698A",
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
	parent_id := d.Get("parent_id").(int)
	name := d.Get("name").(string)
	default_access := d.Get("default_access").(string)
	permissions := d.Get("permissions").(map[string]interface{})
	color := d.Get("color").(string)
	archived := d.Get("archived").(bool)

	col := client.Collection{
		Id:       id,
		ParentId: parent_id,
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

	collectionGraph := client.CollectionGraph{}
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error reading the current collection '%s' permissions", name),
			Detail:   "Could not read the current collection permissions, unexpected error: " + err.Error(),
		})
		return diags
	}

	// Assign the permissions found above
	collectionGraph.Groups = createCollectionPermissions(permissions, fmt.Sprintf("%d", updated.Id), default_access)

	updated_cg, err := c.UpdateCollectionGraph(collectionGraph)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error updating collection '%s' permissions", name),
			Detail:   "Could not update the collection permissions, unexpected error: " + err.Error(),
		})
		return diags
	}

	if updated.ParentId != 0 {
		if err := d.Set("parent_id", updated.ParentId); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("name", updated.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("default_access", default_access); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("permissions", extractCollectionPermissions(updated_cg.Groups, d.Id())); err != nil {
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

	parent_id := d.Get("parent_id").(int)
	name := d.Get("name").(string)
	default_access := d.Get("default_access").(string)
	permissions := d.Get("permissions").(map[string]interface{})
	color := d.Get("color").(string)
	archived := d.Get("archived").(bool)
	col := client.Collection{
		ParentId: parent_id,
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

	// We need to fetch the current collection graph revision
	collectionGraph, err := c.GetCollectionGraph()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error reading the current collection '%s' permissions", name),
			Detail:   "Could not read the current collection permissions, unexpected error: " + err.Error(),
		})
		return diags
	}

	// Assign the permissions found above
	collectionGraph.Groups = createCollectionPermissions(permissions, fmt.Sprintf("%v", created.Id), default_access)

	updated, err := c.UpdateCollectionGraph(collectionGraph)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error updating collection '%s'", name),
			Detail:   "Could not update the collection permissions, unexpected error: " + err.Error(),
		})
		return diags
	}

	d.SetId(fmt.Sprint(created.Id))
	if col.ParentId != 0 {
		if err := d.Set("parent_id", created.ParentId); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("name", created.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("permissions", extractCollectionPermissions(updated.Groups, d.Id())); err != nil {
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

	cg, err := c.GetCollectionGraph()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error read collection '%s' permissions", col.Name),
			Detail:   "Could not update the collection permissions, unexpected error: " + err.Error(),
		})
		return diags
	}

	d.SetId(fmt.Sprintf("%v", col.Id))
	if col.ParentId != 0 {
		if err := d.Set("parent_id", col.ParentId); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("permissions", extractCollectionPermissions(cg.Groups, fmt.Sprintf("%v", col.Id))); err != nil {
		return diag.FromErr(err)
	}
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

func extractCollectionPermissions(cgGroups map[string]map[string]string, collectionId string) map[string]string {
	permissions := make(map[string]string)
	for groupId := range cgGroups {
		// Skip All Users group and & Admin group as settings permissions here is different
		if groupId == "1" || groupId == "2" {
			continue
		}
		if v, found := cgGroups[groupId][collectionId]; found {
			switch v {
			case "none":
				permissions[groupId] = "none"
			case "read":
				permissions[groupId] = "read"
			case "write":
				permissions[groupId] = "write"
			}
		}
	}
	return permissions
}

func createCollectionPermissions(p map[string]interface{}, collectionId string, defaultAccess string) map[string]map[string]string {
	permissions := make(map[string]map[string]string)
	for groupId, access := range p {
		// Skip All Users group and & Admin group as settings permissions here is different
		if groupId == "1" || groupId == "2" {
			continue
		}
		permissions[groupId] = map[string]string{}
		permissions[groupId][collectionId] = access.(string)
	}
	permissions["1"] = map[string]string{}
	permissions["1"][collectionId] = defaultAccess

	return permissions
}
