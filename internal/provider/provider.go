package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	allResources     = make(map[string]func() *schema.Resource)
	registerResource = makeRegisterResourceFunc(allResources, "resource")
)

func New() func() *schema.Provider {
	return func() *schema.Provider {
		provider := &schema.Provider{
			ResourcesMap: resourceFactoriesToMap(allResources),
			Schema: map[string]*schema.Schema{
				"m2m_token": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("M2M_TOKEN", nil),
					Description: `The M2M token that will be used to call API spec service. 
						Required only if ` + "`client_id`" + ` or ` + "`client_secret`" + ` is absent.
						This value can also be sourced from the ` + "`M2M_TOKEN`" + ` environment variable.
					`,
				},
				"client_id": {
					Type:         schema.TypeString,
					Optional:     true,
					RequiredWith: []string{"client_secret"},
					DefaultFunc:  schema.EnvDefaultFunc("CLIENT_ID", nil),
					Description: `The client id that will be used to issue M2M token. 
						Required only if ` + "`m2m_token`" + ` is absent. 
						This value can also be sourced from the ` + "`CLIENT_ID`" + ` environment variable.
					`,
				},
				"client_secret": {
					Type:         schema.TypeString,
					Optional:     true,
					Sensitive:    true,
					RequiredWith: []string{"client_id"},
					DefaultFunc:  schema.EnvDefaultFunc("CLIENT_SECRET", nil),
					Description: `The client secret that will be used to issue M2M token. 
						Required only if ` + "`m2m_token`" + ` is absent.
						This value can also be sourced from the ` + "`CLIENT_SECRET`" + ` environment variable.
					`,
				},
				"environment": {
					Type:         schema.TypeString,
					Optional:     true,
					Default:      "dev",
					Description:  `API spec service envrionemnt. Options: ` + "`dev (Default)`" + `,` + "` qa`" + `,` + "` prod`.",
					ValidateFunc: validation.StringInSlice([]string{"dev", "qa", "prod"}, false),
				},
			},
		}

		provider.ConfigureContextFunc = configure(provider)

		return provider
	}
}

func configure(p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		config := Config{
			M2MToken:     d.Get("m2m_token").(string),
			Environment:  d.Get("environment").(string),
			ClientID:     d.Get("client_id").(string),
			ClientSecret: d.Get("client_secret").(string),
		}

		client, err := config.Client()
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return client, nil
	}
}

func makeRegisterResourceFunc(factories map[string]func() *schema.Resource, resourceType string) func(name string, fn func() *schema.Resource) interface{} {
	return func(name string, fn func() *schema.Resource) interface{} {
		if strings.ToLower(name) != name {
			panic(fmt.Sprintf("cannot register %s %q: name must be lowercase", resourceType, name))
		}

		if _, exists := factories[name]; exists {
			panic(fmt.Sprintf("cannot register %s %q: a %s with the same name already exists", resourceType, name, resourceType))
		}

		factories[name] = fn

		return nil
	}
}

func resourceFactoriesToMap(factories map[string]func() *schema.Resource) map[string]*schema.Resource {
	resourcesMap := make(map[string]*schema.Resource)

	for name, fn := range factories {
		resourcesMap[name] = fn()
	}

	return resourcesMap
}
