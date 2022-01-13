package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Membership struct {
	UserId       int `json:"user_id"`
	GroupId      int `json:"group_id"`
	MembershipId int `json:"membership_id"`
}
type Memberships map[int][]Membership

type groupMembership struct {
	UserId       int `json:"user_id"`
	MembershipId int `json:"membership_id"`
}

func (c *Client) GetMemberships() (Memberships, error) {
	var memberships Memberships
	url := fmt.Sprintf("%s/api/permissions/membership", c.BaseURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return memberships, err
	}
	if err := c.sendRequest(req, &memberships); err != nil {
		return memberships, err
	}

	log.Printf("[DEBUG] Got memberships '%+v'", memberships)
	return memberships, nil
}

func (c *Client) CreateMembership(m Membership) (Membership, error) {
	var created Membership
	var gm []groupMembership
	url := fmt.Sprintf("%s/api/permissions/membership", c.BaseURL)
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(m)
	req, err := http.NewRequest(http.MethodPost, url, b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return created, err
	}
	if err := c.sendRequest(req, &gm); err != nil {
		return created, err
	}

	log.Printf("[INFO] Updated group membership '%+v' by creating '%+v'", gm, m)

	// Find the membership_id for the user
	for _, g := range gm {
		if g.UserId == m.UserId {
			created.MembershipId = g.MembershipId
			created.UserId = g.UserId
			created.GroupId = m.GroupId
			return created, nil
		}
	}

	return created, errors.New("something went wrong in membership creation")
}

func (c *Client) DeleteMembership(membershipId int) error {
	url := fmt.Sprintf("%s/api/permissions/membership/%d", c.BaseURL, membershipId)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return err
	}

	log.Printf("[INFO] Deleted membership by id='%d'", membershipId)
	return nil
}
