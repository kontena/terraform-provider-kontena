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

func exchangeToken(rd *schema.ResourceData, providerMeta *providerMeta, code string) (*client.Token, error) {
	if token, err := providerMeta.config.ExchangeToken(code); err != nil {
		return nil, err
	} else {
		rd.Set("token", token.AccessToken)

		return token, nil
	}
}

func dataKontenaTokenExists(rd *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("[INFO] Kontena Token %v: Exists", rd.Id())

	return false, nil
}

func dataKontenaTokenRead(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(*providerMeta)
	var code = rd.Get("code").(string)
	var token = rd.Get("token").(string)

	log.Printf("[INFO] Kontena Token %v: Read code=%v token=%v", rd.Id(), code, token)

	var clientToken *client.Token

	if _, tokenOK := rd.GetOk("token"); !tokenOK || rd.HasChange("code") {
		log.Printf("[INFO] Kontena Token %v: Exchange code=%v", rd.Id(), code)

		if newToken, err := exchangeToken(rd, providerMeta, code); err != nil {
			return err
		} else {
			clientToken = newToken
		}
	} else {
		log.Printf("[INFO] Kontena Token %v: Use token=%v", rd.Id(), token)

		clientToken = client.MakeToken(token)
	}

	var id = tokenHash(clientToken)

	log.Printf("[INFO] Kontena Token %v: register: %#v", id, clientToken)

	rd.SetId(id)
	providerMeta.registerToken(id, clientToken)

	return nil
}
