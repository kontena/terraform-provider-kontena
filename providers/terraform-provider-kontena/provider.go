package main

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kontena/terraform-provider-kontena/client"
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
			"token": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"kontena_grid": resourceKontenaGrid(),
			"kontena_node": resourceKontenaNode(),
		},
		ConfigureFunc: providerConfigure,
	}
}

type providerMeta struct {
	config client.Config
	client *client.Client
}

func providerClient(meta interface{}) *client.Client {
	var providerMeta = meta.(*providerMeta)

	return providerMeta.client
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var meta = providerMeta{
		config: client.Config{
			URL:   d.Get("url").(string),
			Token: client.MakeToken(d.Get("token").(string)),
		},
	}

	log.Printf("[DEBUG] Kontena: config %#v", meta.config)

	if client, err := meta.config.Connect(); err != nil {
		return nil, err
	} else {
		meta.client = client

		log.Printf("[INFO] Kontena: client %v", meta.client)
	}

	return &meta, nil
}
