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

type User struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Users struct {
	Data []User `json:"data"`
}
