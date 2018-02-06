package kontena

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
			"ssl_cert_pem": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl_cert_cn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"token": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"kontena_token":             resourceKontenaToken(),
			"kontena_grid":              resourceKontenaGrid(),
			"kontena_node":              resourceKontenaNode(),
			"kontena_external_registry": resourceKontenaExternalRegistry(),
		},
		ConfigureFunc: providerConfigure,
	}
}

type providerMeta struct {
	logger *Logger
	config client.Config
	client *client.Client
}

func providerClient(meta interface{}) *client.Client {
	var providerMeta = meta.(*providerMeta)

	return providerMeta.client
}

// Connect with token
func (providerMeta *providerMeta) connectClientWithToken(token *client.Token) (*client.Client, error) {
	var config = providerMeta.config
	config.Token = token

	providerMeta.logger.Debugf("connect %#v (token=%#v)", config, token)

	return config.Connect()
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var logger = Logger{}
	var meta = providerMeta{
		logger: &logger,
		config: client.Config{
			URL:           d.Get("url").(string),
			SSLCertPEM:    []byte(d.Get("ssl_cert_pem").(string)),
			SSLServerName: d.Get("ssl_cert_cn").(string),
			Logger:        &logger,
		},
	}

	logger.Debugf("config %#v", meta.config)

	if tokenValue, ok := d.GetOk("token"); !ok {
		log.Printf("[WARN] Missing token")
	} else if token, err := client.MakeToken(tokenValue.(string)); err != nil {
		return nil, fmt.Errorf("Invalid token: %v", err)
	} else {
		meta.config.Token = token
	}

	logger.Debugf("connect %v (token %v)", meta.config.URL, meta.config.Token)

	// do not test connection; provider can be configured without any url/token when planning
	if client, err := meta.config.MakeClient(); err != nil {
		return nil, err
	} else {
		meta.client = client

		logger.Infof("client %v", meta.client)
	}

	return &meta, nil
}
