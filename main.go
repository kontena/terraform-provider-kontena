package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/kontena/terraform-provider-kontena/kontena"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kontena.Provider,
	})
}
