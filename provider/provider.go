package provider

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kontena/terraform-provider-kontena/client"
)

func provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"access_token": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"kontena_grid": resourceKontenaGrid(),
			"kontena_node": resourceKontenaNode(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"kontena_token": resourceKontenaToken(),
		},
		ConfigureFunc: providerConfigure,
	}
}

type providerMeta struct {
	config        client.Config
	tokens        map[string]*client.Token
	defaultClient *client.Client
}

func (providerMeta *providerMeta) registerToken(id string, token *client.Token) {
	providerMeta.tokens[id] = token
}

func (providerMeta *providerMeta) connectClient(token *client.Token) (*client.Client, error) {
	var config = providerMeta.config
	config.Token = token

	log.Printf("[INFO] Kontena Resource: config %#v", config)

	return config.Connect()
}

func (providerMeta *providerMeta) resourceClient(rd *schema.ResourceData) (*client.Client, error) {
	if value, ok := rd.GetOk("kontena_token"); !ok {
		return providerMeta.defaultClient, nil
	} else if tokenID := value.(string); tokenID == "" {
		return nil, fmt.Errorf("Empty kontena_token given")
	} else if token, ok := providerMeta.tokens[tokenID]; !ok {
		return nil, fmt.Errorf("Missing kontena_token=%v", tokenID)
	} else if client, err := providerMeta.connectClient(token); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func resourceClient(rd *schema.ResourceData, meta interface{}) (*client.Client, error) {
	var providerMeta = meta.(*providerMeta)

	return providerMeta.resourceClient(rd)
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var meta = providerMeta{
		tokens: make(map[string]*client.Token),
		config: client.Config{
			URL: d.Get("url").(string),
		},
	}

	if value, ok := d.GetOk("access_token"); !ok {
		// no token configured
	} else if accessToken := value.(string); accessToken == "" {
		// empty token configured
	} else {
		meta.config.Token = client.MakeToken(value.(string))
	}

	log.Printf("[DEBUG] Kontena: config %#v", meta.config)

	if meta.config.Token == nil {
		// no token given, each resource must reference a kontena_token
	} else if client, err := meta.config.MakeClient(); err != nil {
		return nil, err
	} else if err := client.Ping(); err != nil {
		return nil, err
	} else {
		meta.defaultClient = client

		log.Printf("[INFO] Kontena: client %v", meta.defaultClient)
	}

	return &meta, nil
}

func Provider() terraform.ResourceProvider {
	return provider()
}
