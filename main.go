package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/kontena/terraform-provider-kontena/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
