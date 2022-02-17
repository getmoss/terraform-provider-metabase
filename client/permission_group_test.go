package client

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermissionGroups(t *testing.T) {
	t.Run("Get PermissionGroups", func(t *testing.T) {
		expected := PermissionGroups{{
			Id:          1,
			Name:        "Test Group",
			MemberCount: 1,
		}}
		url := "/api/permissions/group"
		httpMethod := http.MethodGet
		svr := server(url, httpMethod, expected)
		defer svr.Close()
		c := Client{
			BaseURL:    svr.URL,
			HTTPClient: &http.Client{},
		}

		groups, err := c.GetPermissionGroups()

		assert.Nil(t, err)
		assert.Equal(t, expected, groups)
	})

	t.Run("Create PermissionGroup", func(t *testing.T) {
		expected := PermissionGroup{
			Id: 1,
		}
		url := "/api/permissions/group"
		httpMethod := http.MethodPost
		svr := server(url, httpMethod, expected)
		defer svr.Close()
		c := Client{
			BaseURL:    svr.URL,
			HTTPClient: &http.Client{},
		}

		group, err := c.CreatePermissionGroup("test-client")

		assert.Nil(t, err)
		assert.Equal(t, expected, group)
	})

	t.Run("Get PermissionGroup", func(t *testing.T) {
		groupId := 1
		expected := PermissionGroup{
			Id:   groupId,
			Name: "Test Group",
		}
		url := fmt.Sprintf("/api/permissions/group/%d", groupId)
		httpMethod := http.MethodGet
		svr := server(url, httpMethod, expected)
		defer svr.Close()
		c := Client{
			BaseURL:    svr.URL,
			HTTPClient: &http.Client{},
		}

		group, err := c.GetPermissionGroup(groupId)

		assert.Nil(t, err)
		assert.Equal(t, expected, group)
	})

	t.Run("Delete PermissionGroup", func(t *testing.T) {
		groupId := 1

		url := fmt.Sprintf("/api/permissions/group/%d", groupId)
		httpMethod := http.MethodDelete
		svr := server(url, httpMethod, nil)
		defer svr.Close()
		c := Client{
			BaseURL:    svr.URL,
			HTTPClient: &http.Client{},
		}

		err := c.DeletePermissionGroup(groupId)

		assert.Nil(t, err)
	})

	t.Run("Get groups & cache them so that subsequent calls are not making network call", func(t *testing.T) {
		name := "Test Group"
		expected := PermissionGroups{{
			Id:          1,
			Name:        name,
			MemberCount: 1,
		}}
		svr := server("/api/permissions/group", http.MethodGet, expected)
		c := Client{
			BaseURL:    svr.URL,
			HTTPClient: &http.Client{},
		}

		orig, errOrig := c.GetPermissionGroups()
		svr.Close() // close so the mock server is not running
		later, errLater := c.GetPermissionGroups()

		assert.Nil(t, errOrig)
		assert.Nil(t, errLater)
		assert.Equal(t, orig, later)
	})
}
