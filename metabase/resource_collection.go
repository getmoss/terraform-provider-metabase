package metabase

import (
	"context"
	"fmt"
	"log"
	"sort"
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
			"read_access": {
				Description: "List of groups id with read access",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
				Computed: true,
			},
			"write_access": {
				Description: "List of groups id with write access",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
				Computed: true,
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
	read_access := d.Get("read_access").([]interface{})
	write_access := d.Get("write_access").([]interface{})
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

	// Update collection graph groups permissions
	permissions := map[string]map[string]string{}

	collectionGraph, err := c.GetCollectionGraph()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error reading the current collection '%s' permissions", name),
			Detail:   "Could not read the current collection permissions, unexpected error: " + err.Error(),
		})
		return diags
	}

	for i := 0; i < len(read_access); i++ {
		if len(permissions[fmt.Sprint(read_access[i])]) == 0 {
			permissions[fmt.Sprint(read_access[i])] = map[string]string{}
		}
		permissions[fmt.Sprint(read_access[i])][d.Id()] = "read"
	}

	for i := 0; i < len(write_access); i++ {
		if len(permissions[fmt.Sprint(write_access[i])]) == 0 {
			permissions[fmt.Sprint(write_access[i])] = map[string]string{}
		}
		permissions[fmt.Sprint(write_access[i])][d.Id()] = "write"
	}

	collectionGraph.Groups = permissions

	updated_cg, err := c.UpdateCollectionGraph(collectionGraph)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error updating collection '%s' permissions", name),
			Detail:   "Could not update the collection permissions, unexpected error: " + err.Error(),
		})
		return diags
	}

	var updated_read_access []int
	var updated_write_access []int

	for groupId, permission := range updated_cg.Groups {
		groupIdInt, _ := strconv.Atoi(groupId)

		if v, found := permission[d.Id()]; found {
			switch v {
			case "read":
				updated_read_access = append(updated_read_access, groupIdInt)
			case "write":
				updated_write_access = append(updated_write_access, groupIdInt)
			}
		}
	}

	sort.Ints(updated_read_access)
	sort.Ints(updated_write_access)

	if updated.ParentId != 0 {
		if err := d.Set("parent_id", updated.ParentId); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("name", updated.Name); err != nil {
		return diag.FromErr(err)
	}
	if len(updated_read_access) != 0 {
		if err := d.Set("read_access", updated_read_access); err != nil {
			return diag.FromErr(err)
		}
	}
	if len(updated_write_access) != 0 {
		if err := d.Set("write_access", updated_write_access); err != nil {
			return diag.FromErr(err)
		}
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
	read_access := d.Get("read_access").([]interface{})
	write_access := d.Get("write_access").([]interface{})
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

	// Assign collection graph groups permissions
	permissions := map[string]map[string]string{}

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

	for i := 0; i < len(read_access); i++ {
		if len(permissions[fmt.Sprint(read_access[i])]) == 0 {
			permissions[fmt.Sprint(read_access[i])] = map[string]string{}
		}
		permissions[fmt.Sprint(read_access[i])][d.Id()] = "read"
	}

	for i := 0; i < len(write_access); i++ {
		if len(permissions[fmt.Sprint(write_access[i])]) == 0 {
			permissions[fmt.Sprint(write_access[i])] = map[string]string{}
		}
		permissions[fmt.Sprint(write_access[i])][d.Id()] = "write"
	}

	collectionGraph.Groups = permissions

	updated, err := c.UpdateCollectionGraph(collectionGraph)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error updating collection '%s' permissions", name),
			Detail:   "Could not update the collection permissions, unexpected error: " + err.Error(),
		})
		return diags
	}

	var updated_read_access []int
	var updated_write_access []int

	for groupId, permission := range updated.Groups {
		groupIdInt, _ := strconv.Atoi(groupId)

		if v, found := permission[d.Id()]; found {
			switch v {
			case "read":
				updated_read_access = append(updated_read_access, groupIdInt)
			case "write":
				updated_write_access = append(updated_write_access, groupIdInt)
			}
		}
	}

	sort.Ints(updated_read_access)
	sort.Ints(updated_write_access)

	d.SetId(fmt.Sprint(created.Id))
	if col.ParentId != 0 {
		if err := d.Set("parent_id", created.ParentId); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("name", created.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("read_access", updated_read_access); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("write_access", updated_write_access); err != nil {
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

	var read_access []int
	var write_access []int

	for groupId := range cg.Groups {
		if v, found := cg.Groups[groupId][d.Id()]; found {
			groupIdInt, _ := strconv.Atoi(groupId)

			switch v {
			case "read":
				read_access = append(read_access, groupIdInt)
			case "write":
				write_access = append(write_access, groupIdInt)
			}
		}

	}

	sort.Ints(read_access)
	sort.Ints(write_access)

	d.SetId(fmt.Sprintf("%v", col.Id))
	if col.ParentId != 0 {
		if err := d.Set("parent_id", col.ParentId); err != nil {
			return diag.FromErr(err)
		}
	}
	if len(read_access) != 0 {
		if err := d.Set("read_access", read_access); err != nil {
			return diag.FromErr(err)
		}
	}
	if len(write_access) != 0 {
		if err := d.Set("write_access", write_access); err != nil {
			return diag.FromErr(err)
		}
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
