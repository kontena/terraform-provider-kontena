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

func resourceKontenaToken() *schema.Resource {
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

		Create: resourceKontenaTokenCreate,
		Read:   resourceKontenaTokenRead,
		Delete: resourceKontenaTokenDelete,
	}
}

func tokenID(token *client.Token) string {
	return fmt.Sprintf("%s", token.AccessToken)
}

func resourceKontenaTokenSync(rd *schema.ResourceData, providerMeta *providerMeta, token *client.Token) {
	rd.Set("token", token.AccessToken)
	rd.SetId(tokenID(token))
}

func resourceKontenaTokenCreate(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(*providerMeta)
	var code = rd.Get("code").(string)

	if clientToken, err := providerMeta.config.ExchangeToken(code); err != nil {
		return err

	} else {
		log.Printf("[DEBUG] Kontena Token: Create code=%v: %#v", code, clientToken)

		resourceKontenaTokenSync(rd, providerMeta, clientToken)
	}

	return nil
}

func resourceKontenaTokenRead(rd *schema.ResourceData, meta interface{}) error {
	var token = rd.Get("token").(string)

	var clientToken = client.MakeToken(token)

	log.Printf("[INFO] Kontena Token %v: Read: %#v", rd.Id(), clientToken)

	return nil
}

func resourceKontenaTokenDelete(rd *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Kontena Token %v: Delete", rd.Id())

	return nil
}
