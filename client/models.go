package client

type Member struct{}
type PermissionGroup struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	// Members []Member `json:"members"`
}

type PermissionGroups []struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	MemberCount int    `json:"member_count"`
}
