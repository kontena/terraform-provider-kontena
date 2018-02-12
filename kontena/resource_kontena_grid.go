package kontena

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kontena/kontena-client-go/api"
	"github.com/kontena/kontena-client-go/client"
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
				Optional: true,
				Computed: true,
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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func makeKontenaGridCreateParams(rd *schema.ResourceData) api.GridPOST {
	var gridParams = api.GridPOST{
		Name:        rd.Get("name").(string),
		InitialSize: rd.Get("initial_size").(int),
	}

	if value, ok := rd.GetOk("token"); ok {
		var token = value.(string)

		gridParams.Token = token
	}

	if value, ok := rd.GetOk("subnet"); ok {
		var subnet = value.(string)

		gridParams.Subnet = subnet
	}

	if value, ok := rd.GetOk("supernet"); ok {
		var supernet = value.(string)

		gridParams.Supernet = supernet
	}

	if value, ok := rd.GetOk("default_affinity"); ok {
		var defaultAffinity = api.GridDefaultAffinity(value.([]string))

		gridParams.DefaultAffinity = &defaultAffinity
	}

	if value, ok := rd.GetOk("default_affinity"); ok {
		var defaultAffinity = api.GridDefaultAffinity(value.([]string))

		gridParams.DefaultAffinity = &defaultAffinity
	}

	return gridParams
}

func makeKontenaGridUpdateParams(rd *schema.ResourceData) api.GridPUT {
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

	return gridParams
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
	var apiClient = providerClient(meta)
	var gridParams = makeKontenaGridCreateParams(rd)

	log.Printf("[INFO] Kontena Grid: Create %#v", gridParams)

	if grid, err := apiClient.Grids.Create(gridParams); err != nil {
		return err
	} else {
		rd.SetId(grid.String())
		syncKontenaGrid(rd, grid)
	}

	return nil
}

func resourceKontenaGridRead(rd *schema.ResourceData, meta interface{}) error {
	var apiClient = providerClient(meta)

	if grid, err := apiClient.Grids.Get(rd.Id()); err == nil {
		log.Printf("[INFO] Kontena Grid %v: Read: %#v", rd.Id(), grid)

		syncKontenaGrid(rd, grid)
	} else if _, ok := err.(client.NotFoundError); ok {
		log.Printf("[INFO] Kontena Grid %v: Read gone", rd.Id())

		rd.SetId("")
	} else {
		log.Printf("[INFO] Kontena Grid %v: Read error: %v", rd.Id(), err)

		return err
	}

	return nil
}

func resourceKontenaGridUpdate(rd *schema.ResourceData, meta interface{}) error {
	var apiClient = providerClient(meta)
	var gridParams = makeKontenaGridUpdateParams(rd)

	log.Printf("[INFO] Kontena Grid %v: Update %#v", rd.Id(), gridParams)

	if grid, err := apiClient.Grids.Update(rd.Id(), gridParams); err != nil {
		return err
	} else {
		syncKontenaGrid(rd, grid)
	}

	return nil
}

func resourceKontenaGridDelete(rd *schema.ResourceData, meta interface{}) error {
	var apiClient = providerClient(meta)

	log.Printf("[INFO] Kontena Grid %v: Delete", rd.Id())

	if err := apiClient.Grids.Delete(rd.Id()); err != nil {
		return err
	}

	return nil
}
