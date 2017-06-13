package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
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
		},
	}
}
