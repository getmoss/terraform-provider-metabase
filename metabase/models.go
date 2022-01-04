package metabase

type PermissionGroup struct {
	Id   int    `tfsdk:"id"`
	Name string `tfsdk:"name"`
	// Members []Member `tfsdk:"members"`
}

type Member struct{}
