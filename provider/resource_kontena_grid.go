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
				ForceNew: true,
			},
			"initial_size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"token": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"subnet": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"supernet": &schema.Schema{
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
			"trusted_subnets": &schema.Schema{
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

func syncKontenaGrid(rd *schema.ResourceData, grid api.Grid) {
	rd.Set("name", grid.Name)
	rd.Set("initial_size", grid.InitialSize)
	rd.Set("token", grid.Token)
	rd.Set("subnet", grid.Subnet)
	rd.Set("supernet", grid.Supernet)
	rd.Set("default_affinity", grid.DefaultAffinity)
	rd.Set("trusted_subnets", grid.TrustedSubnets)
}

func resourceKontenaGridCreate(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(providerMeta)
	var gridParams = api.GridPOST{
		Name:        rd.Get("name").(string),
		InitialSize: rd.Get("initial_size").(int),
	}

	// XXX: no POST support for default_affinity or trusted_subnets
	//     can only update it on a later apply
	if value, ok := rd.GetOk("token"); ok {
		var token = value.(string)

		gridParams.Token = &token
	}
	if value, ok := rd.GetOk("subnet"); ok {
		var subnet = value.(string)

		gridParams.Subnet = &subnet
	}
	if value, ok := rd.GetOk("supernet"); ok {
		var supernet = value.(string)

		gridParams.Supernet = &supernet
	}

	log.Printf("[INFO] Kontena Grid: Create %#v", gridParams)

	if grid, err := providerMeta.client.Grids.Create(gridParams); err != nil {
		return err
	} else {
		rd.SetId(grid.String())
		syncKontenaGrid(rd, grid)
	}

	return nil
}

func resourceKontenaGridRead(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(providerMeta)

	log.Printf("[INFO] Kontena Grid %v: Read", rd.Id())

	if grid, err := providerMeta.client.Grids.Get(rd.Id()); err != nil {
		return err
	} else {
		syncKontenaGrid(rd, grid)
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
	if rd.HasChange("trusted_subnets") {
		var trustedSubnets = make(api.GridTrustedSubnets, 0)

		for _, value := range rd.Get("trusted_subnets").([]interface{}) {
			trustedSubnets = append(trustedSubnets, value.(string))
		}

		gridParams.TrustedSubnets = &trustedSubnets
	}

	log.Printf("[INFO] Kontena Grid %v: Update %#v", rd.Id(), gridParams)

	if grid, err := providerMeta.client.Grids.Update(rd.Id(), gridParams); err != nil {
		return err
	} else {
		syncKontenaGrid(rd, grid)
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
