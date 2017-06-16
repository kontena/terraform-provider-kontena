package main

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kontena/kontena-client-go/client"
)

func Provider() terraform.ResourceProvider {
	return provider()
}

func provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"kontena-oauth2_token": resourceKontenaToken(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var meta = providerMeta{
		config: client.Config{
			URL: d.Get("url").(string),
		},
	}

	log.Printf("[DEBUG] Kontena-OAuth2: config %#v", meta.config)

	return &meta, nil
}

type providerMeta struct {
	config client.Config
}

func (providerMeta *providerMeta) connectClient(token *client.Token) (*client.Client, error) {
	var config = providerMeta.config
	config.Token = token

	log.Printf("[DEBUG] Kontena-OAuth2: connect %#v (token=%#v)", config, token)

	return config.Connect()
}
