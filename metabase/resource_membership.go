package metabase

import (
	"context"
	"strconv"
	"terraform-provider-metabase/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceMembershipDelete(_ context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	membershipId, _ := strconv.Atoi(d.Id())

	if err := c.DeleteMembership(membershipId); err != nil {
		return diag.Errorf("error deleting membership: %s for membershipId=[%d]", err, membershipId)
	}
	return
}

func resourceMembershipCreate(_ context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	userId := d.Get("user_id").(int)
	groupId := d.Get("group_id").(int)
	m := client.Membership{
		UserId:  userId,
		GroupId: groupId,
	}

	created, err := c.CreateMembership(m)
	if err != nil {
		return diag.Errorf("error creating membership: %s for userId=[%d]", err, userId)
	}

	d.SetId(strconv.Itoa(created.MembershipId))
	if err := d.Set("user_id", userId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("membership_id", created.MembershipId); err != nil {
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

	memberships, err := c.GetMemberships()
	if err != nil {
		return diag.FromErr(err)
	}

	m = findMatchingMembership(memberships, membershipId)

	if m == (client.Membership{}) {
		return diag.Errorf("Could not find Membership by id [%d] in memberships[%+v]", membershipId, memberships)
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
			if membership.MembershipId == membershipId {
				return membership
			}
		}
	}
	return
}
