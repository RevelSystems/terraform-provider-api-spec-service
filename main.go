package main

import (
	"terraform-provider-api-spec-service/internal/provider"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	opts := &plugin.ServeOpts{ProviderFunc: provider.New()}
	plugin.Serve(opts)
}
