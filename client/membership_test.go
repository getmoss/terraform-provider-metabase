package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupMembership(t *testing.T) {
	t.Run("Get Memberships", func(t *testing.T) {
		expected := Memberships{
			1: []Membership{
				{
					UserId:       1,
					GroupId:      2,
					MembershipId: 3,
				},
			},
		}
		url := "/api/permissions/membership"
		httpMethod := http.MethodGet
		svr := server(url, httpMethod, expected)
		defer svr.Close()

		c := Client{
			BaseURL:    svr.URL,
			HTTPClient: &http.Client{},
		}

		us, err := c.GetMemberships()

		assert.Nil(t, err)
		assert.Equal(t, expected, us)
	})

	t.Run("Create Membership", func(t *testing.T) {
		membershipToBeCreated := Membership{
			GroupId: 1,
			UserId:  2,
		}
		membershipId := 3
		expected := Membership{
			GroupId:      1,
			UserId:       2,
			MembershipId: membershipId,
		}
		groupMembership := []groupMembership{{
			UserId:       membershipToBeCreated.UserId,
			MembershipId: membershipId,
		}}
		url := "/api/permissions/membership"
		httpMethod := http.MethodPost
		svr := server(url, httpMethod, groupMembership)
		defer svr.Close()

		c := Client{
			BaseURL:    svr.URL,
			HTTPClient: &http.Client{},
		}

		us, err := c.CreateMembership(membershipToBeCreated)

		assert.Nil(t, err)
		assert.Equal(t, expected, us)
	})

	t.Run("Delete Membership", func(t *testing.T) {
		membershipId := 1

		mux := http.NewServeMux()
		mux.HandleFunc(fmt.Sprintf("/api/permissions/membership/%d", membershipId), func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodDelete:
				w.WriteHeader(http.StatusNoContent)
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
		})
		svr := httptest.NewServer(mux)
		defer svr.Close()

		c := Client{
			BaseURL:    svr.URL,
			HTTPClient: &http.Client{},
		}

		err := c.DeleteMembership(membershipId)

		assert.Nil(t, err)
	})
}
