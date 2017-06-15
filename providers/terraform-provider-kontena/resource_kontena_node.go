package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kontena/terraform-provider-kontena/api"
	"github.com/kontena/terraform-provider-kontena/client"
)

func resourceKontenaNode() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// configured identifier
			"grid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// computed identifier
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			// updatable attributes
			"labels": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Computed: true,
			},

			// computed attributes
			"node_number": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"initial_node": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"public_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"overlay_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"agent_version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"docker_version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},

		Create: resourceKontenaNodeCreate,
		Read:   resourceKontenaNodeRead,
		Update: resourceKontenaNodeUpdate,
		Delete: resourceKontenaNodeDelete,
	}
}

func resourceKontenaNodeSync(rd *schema.ResourceData, node api.Node) {
	rd.Set("grid", node.Grid.Name)
	rd.Set("name", node.Name)
	rd.Set("id", node.ID)
	rd.Set("labels", node.Labels)

	rd.Set("node_number", node.NodeNumber)
	rd.Set("initial_node", node.InitialMember)
	rd.Set("public_ip", node.PublicIP)
	rd.Set("private_ip", node.PrivateIP)
	rd.Set("overlay_ip", node.OverlayIP)
	rd.Set("agent_version", node.AgentVersion)
	rd.Set("docker_version", node.DockerVersion)
}

// The Kontena nodes API does not support POST to create nodes
// Instead, nodes will be created once they connect for the first time
// The kontena_node create method just waits for a node with the given grid/name to show up
func resourceKontenaNodeCreate(rd *schema.ResourceData, meta interface{}) error {
	var apiClient = providerClient(meta)
	var nodeID = client.NodeID{
		Grid: rd.Get("grid").(string),
		Name: rd.Get("name").(string),
	}

	log.Printf("[INFO] Kontena Node: Create %v", nodeID)

	// Wait for node to show up
	for {
		if node, err := apiClient.Nodes.Get(nodeID); err == nil {
			rd.SetId(nodeID.String())
			resourceKontenaNodeSync(rd, node)

			break

		} else if _, ok := err.(client.NotFoundError); ok {
			continue

		} else {
			return err
		}
	}

	return nil
}

func resourceKontenaNodeRead(rd *schema.ResourceData, meta interface{}) error {
	var apiClient = providerClient(meta)

	nodeID, err := client.ParseNodeID(rd.Id())
	if err != nil {
		return fmt.Errorf("Invalid node ID %#v: %v", rd.Id(), err)
	}

	if node, err := apiClient.Nodes.Get(nodeID); err == nil {
		log.Printf("[INFO] Kontena Node %v: Read: %#v", rd.Id(), node)

		resourceKontenaNodeSync(rd, node)

	} else if _, ok := err.(client.NotFoundError); ok {
		log.Printf("[INFO] Kontena Grid %v: Read gone", rd.Id())

		rd.SetId("")

	} else {
		log.Printf("[INFO] Kontena Node %v: Read error: %v", rd.Id(), err)

		return err
	}

	return nil
}

func resourceKontenaNodeUpdate(rd *schema.ResourceData, meta interface{}) error {
	var apiClient = providerClient(meta)

	nodeID, err := client.ParseNodeID(rd.Id())
	if err != nil {
		return fmt.Errorf("Invalid node ID %#v: %v", rd.Id(), err)
	}

	var nodeParams = api.NodePUT{}

	if rd.HasChange("labels") {
		var labels = make(api.NodeLabels, 0)

		for _, value := range rd.Get("labels").([]interface{}) {
			labels = append(labels, value.(string))
		}

		nodeParams.Labels = &labels
	}

	log.Printf("[INFO] Kontena Grid %v: Update %#v", rd.Id(), nodeParams)

	if node, err := apiClient.Nodes.Update(nodeID, nodeParams); err != nil {
		return err
	} else {
		resourceKontenaNodeSync(rd, node)
	}

	return nil
}

func resourceKontenaNodeDelete(rd *schema.ResourceData, meta interface{}) error {
	var apiClient = providerClient(meta)

	nodeID, err := client.ParseNodeID(rd.Id())
	if err != nil {
		return fmt.Errorf("Invalid node ID %#v: %v", rd.Id(), err)
	}

	log.Printf("[INFO] Kontena Node %v: Delete", rd.Id())

	if err := apiClient.Nodes.Delete(nodeID); err != nil {
		return err
	}

	return nil
}
