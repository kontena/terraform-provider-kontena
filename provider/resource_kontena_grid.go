package provider

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kontena/terraform-provider-kontena/api"
)

func resourceKontenaGrid() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"initial_size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"token": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"default_affinity": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},

		Create: resourceKontenaGridCreate,
		Read:   resourceKontenaGridRead,
		Update: resourceKontenaGridUpdate,
		Delete: resourceKontenaGridDelete,
	}
}

func resourceKontenaGridCreate(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(providerMeta)
	var gridParams = api.GridPOST{
		Name:        rd.Get("name").(string),
		InitialSize: rd.Get("initial_size").(int),
	}

	if value, ok := rd.GetOk("token"); ok {
		var token = value.(string)

		gridParams.Token = &token
	}

	log.Printf("[INFO] Kontena Grid: Create %#v", gridParams)

	if grid, err := providerMeta.client.Grids.Create(gridParams); err != nil {
		return err
	} else {
		rd.SetId(grid.String())
	}

	return nil
}

func resourceKontenaGridRead(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(providerMeta)

	log.Printf("[INFO] Kontena Grid %v: Read", rd.Id())

	if grid, err := providerMeta.client.Grids.Get(rd.Id()); err != nil {
		return err
	} else {
		rd.Set("name", grid.Name)
		rd.Set("initial_size", grid.InitialSize)
		rd.Set("token", grid.Token)

		rd.Set("default_affinity", grid.DefaultAffinity)
	}

	return nil
}

func resourceKontenaGridUpdate(rd *schema.ResourceData, meta interface{}) error {

	var providerMeta = meta.(providerMeta)
	var gridParams = api.GridPUT{}

	if rd.HasChange("default_affinity") {
		var defaultAffinity = make(api.GridDefaultAffinity, 0)

		for _, value := range rd.Get("default_affinity").([]interface{}) {
			defaultAffinity = append(defaultAffinity, value.(string))
		}

		gridParams.DefaultAffinity = &defaultAffinity
	}

	log.Printf("[INFO] Kontena Grid %v: Update %#v (default_affinity=%#v)", rd.Id(), gridParams, gridParams.DefaultAffinity)

	if grid, err := providerMeta.client.Grids.Update(rd.Id(), gridParams); err != nil {
		return err
	} else {
		rd.Set("default_affinity", grid.DefaultAffinity)
	}

	return nil
}

func resourceKontenaGridDelete(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(providerMeta)

	log.Printf("[INFO] Kontena Grid %v: Delete", rd.Id())

	if err := providerMeta.client.Grids.Delete(rd.Id()); err != nil {
		return err
	}

	return nil
}
