package main

import (
	"fmt"
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
			URL: d.Get("url").(string),
		},
	}

	log.Printf("[DEBUG] Kontena: config %#v", meta.config)

	if tokenValue, ok := d.GetOk("token"); !ok {
		log.Printf("[WARN] Missing token")
	} else if token, err := client.MakeToken(tokenValue.(string)); err != nil {
		return nil, fmt.Errorf("Invalid token: %v", err)
	} else {
		meta.config.Token = token
	}

	log.Printf("[DEBUG] Kontena: connect %v (token %v)", meta.config.URL, meta.config.Token)

	// do not test connection; provider can be configured without any url/token when planning
	if client, err := meta.config.MakeClient(); err != nil {
		return nil, err
	} else {
		meta.client = client

		log.Printf("[INFO] Kontena: client %v", meta.client)
	}

	return &meta, nil
}
