package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	allResources = make(map[string]func() *schema.Resource)
)

var registerResource = makeRegisterResourceFunc(allResources, "resource")

func New() func() *schema.Provider {
	return func() *schema.Provider {
		provider := &schema.Provider{
			ResourcesMap: resourceFactoriesToMap(allResources),
			Schema: map[string]*schema.Schema{
				"m2m_token": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("M2M_TOKEN", nil),
					Description: "The M2M token that will be used to call API spec service",
				},
				"environment": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "dev",
					Description: "API spec service envrionemnt [dev (Default), qa, prod]",
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
			M2MToken:    d.Get("m2m_token").(string),
			Environment: d.Get("environment").(string),
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
