package kontena

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kontena/kontena-client-go/api"
	"github.com/kontena/kontena-client-go/client"
)

func resourceKontenaExternalRegistry() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// configured identifier
			"grid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// computed attributes
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		Create: resourceKontenaExternalRegistryCreate,
		Read:   resourceKontenaExternalRegistryRead,
		Delete: resourceKontenaExternalRegistryDelete,
	}
}

func setKontenaExternalRegistry(rd *schema.ResourceData, externalRegistry api.ExternalRegistry) {
	rd.Set("name", externalRegistry.Name)
}

func resourceKontenaExternalRegistryCreate(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(*providerMeta)
	var gridName = rd.Get("grid").(string)
	var params = api.ExternalRegistryPOST{
		URL:      rd.Get("url").(string),
		Username: rd.Get("username").(string),
		Password: rd.Get("password").(string),
	}

	if email, ok := rd.GetOk("email"); ok {
		params.Email = email.(string)
	}

	providerMeta.logger.Infof("ExternalRegistry: Create %v: %#v", gridName, params)

	if externalRegistry, err := providerMeta.client.ExternalRegistries.Create(gridName, params); err != nil {
		return fmt.Errorf("ExternalRegistry create: %v", err)

	} else {
		rd.SetId(externalRegistry.ID)

		providerMeta.logger.Infof("ExternalRegistry %v: Create: %#v", rd.Id(), externalRegistry)

		setKontenaExternalRegistry(rd, externalRegistry)
	}

	return nil
}

func resourceKontenaExternalRegistryRead(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(*providerMeta)
	var id = rd.Id()

	if externalRegistry, err := providerMeta.client.ExternalRegistries.Get(id); err == nil {
		providerMeta.logger.Infof("ExternalRegistry %v: Read: %#v", rd.Id(), externalRegistry)

		setKontenaExternalRegistry(rd, externalRegistry)

	} else if _, ok := err.(client.NotFoundError); ok {
		providerMeta.logger.Infof("ExternalRegistry %v: Read gone", rd.Id())

		rd.SetId("")

	} else {
		providerMeta.logger.Warnf("ExternalRegistry %v: Read error: %v", rd.Id(), err)

		return err
	}

	return nil
}

func resourceKontenaExternalRegistryDelete(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(*providerMeta)
	var id = rd.Id()

	providerMeta.logger.Infof("ExternalRegistry %v: Delete", id)

	if err := providerMeta.client.ExternalRegistries.Delete(id); err != nil {
		return err
	}

	return nil
}
