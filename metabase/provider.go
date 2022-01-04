package metabase

import (
	"context"
	"terraform-provider-metabase/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			DataSourcesMap: map[string]*schema.Resource{
				"metabase_permission_group": dataSourcePermissionGroup(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"metabase_permission_group": resourcePermissionGroup(),
			},
			Schema: map[string]*schema.Schema{
				"host": {
					Type:        schema.TypeString,
					Description: "Hostname with protocol http/https",
					Required:    true,
				},
				"username": {
					Type:        schema.TypeString,
					Description: "User email of a super admin",
					Required:    true,
					Sensitive:   true,
				},
				"password": {
					Type:        schema.TypeString,
					Description: "User password",
					Required:    true,
					Sensitive:   true,
				},
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		// Setup a User-Agent for your API client (replace the provider name for yours):
		userAgent := p.UserAgent("terraform-provider-metabase", version)
		username := d.Get("username").(string)
		password := d.Get("password").(string)
		host := d.Get("host").(string)

		// Warning or errors can be collected in a slice type
		var diags diag.Diagnostics

		if username == "" || password == "" || host == "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create client",
				Detail:   "Please provide a valid username, password and host",
			})
			return nil, diags
		}

		c, err := client.NewClient(host, username, password, userAgent)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return c, diags
	}
}
