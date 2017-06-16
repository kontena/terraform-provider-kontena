package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kontena/kontena-client-go/client"
)

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

func resourceKontenaTokenCreate(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(*providerMeta)
	var code = rd.Get("code").(string)

	if clientToken, err := providerMeta.config.ExchangeToken(code); err != nil {
		return err

	} else {
		log.Printf("[DEBUG] Kontena-OAuth2 Token: Create code=%v: %#v", code, clientToken)

		rd.Set("token", clientToken.AccessToken)
		rd.SetId(tokenID(clientToken))
	}

	return nil
}

func resourceKontenaTokenRead(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(*providerMeta)

	var code = rd.Get("code").(string)
	var token = rd.Get("token").(string)

	// check token still exists
	if clientToken, err := client.MakeToken(token); err != nil {
		// XXX: corrupt, force re-exchange?
		return err
	} else if apiClient, err := providerMeta.connectClient(clientToken); err != nil {
		return err
	} else if user, err := apiClient.Users.GetUser(); err == nil {
		log.Printf("[INFO] Kontena-OAuth2 Token %v: Read code=%v token=%v ok: %#v", rd.Id(), code, token, user)
	} else if forbiddenError, ok := err.(client.ForbiddenError); ok {
		log.Printf("[INFO] Kontena-OAuth2 Token %v: Read code=%v token=%v gone: %v", rd.Id(), code, token, forbiddenError)

		rd.SetId("")

	} else {
		log.Printf("[INFO] Kontena-OAuth2 Token %v: Read code=%v token=%v err: %v", rd.Id(), code, token, err)

		return err
	}

	return nil
}

func resourceKontenaTokenDelete(rd *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Kontena-OAuth2 Token %v: Delete", rd.Id())

	// TODO: revoke

	return nil
}
