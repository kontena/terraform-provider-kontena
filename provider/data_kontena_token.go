package provider

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kontena/terraform-provider-kontena/client"
)

var kontenaTokenSchema = &schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
}

func dataKontenaToken() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"code": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"token": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		Exists: dataKontenaTokenExists,
		Read:   dataKontenaTokenRead,
	}
}

func tokenHash(token *client.Token) string {
	return fmt.Sprintf("%s", token.AccessToken)
}

func dataKontenaTokenExists(rd *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("[INFO] Kontena Token %v: Exists?", rd.Id())

	return false, nil
}

func dataKontenaTokenRead(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(*providerMeta)
	var code = rd.Get("code").(string)
	var token = rd.Get("token").(string)

	log.Printf("[DEBUG] Kontena Token %v: Read code=%v token=%v", rd.Id(), code, token)

	var clientToken *client.Token

	if rd.Id() != "" {
		log.Printf("[INFO] Kontena Token %v: Use token=%v", rd.Id(), token)

		clientToken = client.MakeToken(token)

	} else if exchangeToken, err := providerMeta.config.ExchangeToken(code); err != nil {
		return err

	} else {
		log.Printf("[INFO] Kontena Token %v: Exchange code=%v", rd.Id(), code)

		clientToken = exchangeToken

		rd.Set("token", exchangeToken.AccessToken)
	}

	var id = tokenHash(clientToken)

	log.Printf("[INFO] Kontena Token %v: register: %#v", id, clientToken)

	rd.SetId(id)
	providerMeta.registerToken(id, clientToken)

	return nil
}
