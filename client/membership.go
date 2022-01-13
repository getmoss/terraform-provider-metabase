package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type UserId int
type GroupId int
type MembershipId int

type Membership struct {
	UserId       UserId       `json:"user_id"`
	GroupId      GroupId      `json:"group_id"`
	MembershipId MembershipId `json:"membership_id"`
}
type Memberships map[UserId][]Membership

type groupMembership struct {
	UserId       UserId       `json:"user_id"`
	MembershipId MembershipId `json:"membership_id"`
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

func (c *Client) CreateMembership(m Membership) (MembershipId, error) {
	var gm []groupMembership
	url := fmt.Sprintf("%s/api/permissions/membership", c.BaseURL)
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(m)
	req, err := http.NewRequest(http.MethodPost, url, b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return 0, err
	}
	if err := c.sendRequest(req, &gm); err != nil {
		return 0, err
	}

	log.Printf("[INFO] Updated group membership '%+v'", gm)

	// Find the membership_id for the user
	for _, g := range gm {
		if g.UserId == m.UserId {
			return g.MembershipId, nil
		}
	}

	return 0, errors.New("something went wrong in membership creation")
}

func (c *Client) DeleteMembership(membershipId MembershipId) error {
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
