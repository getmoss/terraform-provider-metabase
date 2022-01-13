package metabase

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"terraform-provider-metabase/client"
)

func resourceMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMembershipCreate,
		ReadContext:   resourceMembershipRead,
		DeleteContext: resourceMembershipDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"membership_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceMembershipDelete(_ context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	membershipId, _ := strconv.Atoi(d.Id())

	if err := c.DeleteMembership(client.MembershipId(membershipId)); err != nil {
		return diag.Errorf("error deleting membership: %s for membershipId=[%d]", err, membershipId)
	}
	return
}

func resourceMembershipCreate(_ context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	userId := d.Get("user_id").(int)
	groupId := d.Get("group_id").(int)
	m := client.Membership{
		UserId:  client.UserId(userId),
		GroupId: client.GroupId(groupId),
	}

	created, err := c.CreateMembership(m)
	if err != nil {
		return diag.Errorf("error creating membership: %s for userId=[%d]", err, userId)
	}

	d.SetId(strconv.Itoa(int(created)))
	if err := d.Set("user_id", userId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("membership_id", created); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("group_id", groupId); err != nil {
		return diag.FromErr(err)
	}
	return
}

func resourceMembershipRead(_ context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	membershipId, _ := strconv.Atoi(d.Id())

	c := meta.(*client.Client)
	var m client.Membership
	if memberships, err := c.GetMemberships(); err != nil {
		m = findMatchingMembership(memberships, membershipId)

		if m == (client.Membership{}) {
			return diag.Errorf("Could not find Membership by id [%d] in memberships[%+v]", membershipId, memberships)
		}
	}

	if err := d.Set("user_id", m.UserId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("group_id", m.GroupId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("membership_id", m.MembershipId); err != nil {
		return diag.FromErr(err)
	}

	return
}

func findMatchingMembership(memberships client.Memberships, membershipId int) (m client.Membership) {
	for _, userMemberships := range memberships {
		for _, membership := range userMemberships {
			if membership.MembershipId == client.MembershipId(membershipId) {
				return membership
			}
		}
	}
	return
}
