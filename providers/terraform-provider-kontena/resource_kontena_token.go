package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kontena/kontena-client-go/api"
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
			"user": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"roles": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func readKontenaToken(providerMeta *providerMeta, rd *schema.ResourceData) error {
	var code = rd.Get("code").(string)
	var token = rd.Get("token").(string)
	var user api.User

	// check token still exists
	if clientToken, err := client.MakeToken(token); err != nil {
		// XXX: corrupt, force re-exchange?
		return err
	} else if apiClient, err := providerMeta.connectClientWithToken(clientToken); err != nil {
		return err
	} else if getUser, err := apiClient.Users.GetUser(); err == nil {
		providerMeta.logger.Infof("Token %v: Read code=%v token=%v ok: %#v", rd.Id(), code, token, getUser)

		user = getUser

	} else if forbiddenError, ok := err.(client.ForbiddenError); ok {
		providerMeta.logger.Infof("Token %v: Read code=%v token=%v gone: %v", rd.Id(), code, token, forbiddenError)

		rd.SetId("")

		return nil

	} else {
		providerMeta.logger.Infof("Token %v: Read code=%v token=%v err: %v", rd.Id(), code, token, err)

		return err
	}

	rd.Set("user", user.Name)
	rd.Set("email", user.Email)

	var roles []string
	for _, apiRole := range user.Roles {
		roles = append(roles, apiRole.Name)
	}
	rd.Set("roles", roles)

	return nil
}

func resourceKontenaTokenCreate(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(*providerMeta)
	var code = rd.Get("code").(string)

	if clientToken, err := providerMeta.config.ExchangeToken(code); err != nil {
		return err

	} else {
		providerMeta.logger.Infof("Token: Create code=%v: %#v", code, clientToken)

		rd.Set("token", clientToken.AccessToken)
		rd.SetId(tokenID(clientToken))
	}

	// sync
	return readKontenaToken(providerMeta, rd)
}

func resourceKontenaTokenRead(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(*providerMeta)

	return readKontenaToken(providerMeta, rd)
}

func resourceKontenaTokenDelete(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(*providerMeta)

	providerMeta.logger.Infof("Token %v: Delete", rd.Id())

	// TODO: revoke

	return nil
}
