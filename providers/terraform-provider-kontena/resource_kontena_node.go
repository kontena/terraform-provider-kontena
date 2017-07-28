package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kontena/kontena-client-go/api"
	"github.com/kontena/kontena-client-go/client"
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

			// updatable attributes
			"token": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"labels": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Computed: true,
			},

			// computed attributes
			"node_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
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

func setKontenaNode(rd *schema.ResourceData, node api.Node) {
	rd.Set("grid", node.Grid.Name)
	rd.Set("name", node.Name)
	rd.Set("id", node.ID)
	rd.Set("node_id", node.NodeID)
	rd.Set("labels", node.Labels)

	rd.Set("node_number", node.NodeNumber)
	rd.Set("initial_node", node.InitialMember)
	rd.Set("public_ip", node.PublicIP)
	rd.Set("private_ip", node.PrivateIP)
	rd.Set("overlay_ip", node.OverlayIP)
	rd.Set("agent_version", node.AgentVersion)
	rd.Set("docker_version", node.DockerVersion)
}

func getKontenaNodeLabels(rd *schema.ResourceData) *api.NodeLabels {
	var labels = make(api.NodeLabels, 0)

	for _, value := range rd.Get("labels").([]interface{}) {
		labels = append(labels, value.(string))
	}

	return &labels
}

// Get and sync node token from API
func readKontenaNodeToken(providerMeta *providerMeta, rd *schema.ResourceData) error {
	if nodeID, err := client.ParseNodeID(rd.Id()); err != nil {
		return fmt.Errorf("Invalid node ID %#v: %v", rd.Id(), err)

	} else if nodeToken, err := providerMeta.client.Nodes.GetToken(nodeID); err == nil {
		providerMeta.logger.Infof("Node %v: Read token: %#v", rd.Id(), nodeToken)

		rd.Set("token", nodeToken.Token) // XXX: current 1.4.0.pre5 returns null instead of HTTP 404

	} else if _, ok := err.(client.NotFoundError); ok {
		providerMeta.logger.Infof("Node %v: Read token gone", rd.Id())

		rd.Set("token", "")

	} else {
		providerMeta.logger.Warnf("Node %v: Read token error: %v", rd.Id(), err)

		return err
	}

	return nil
}

func resourceKontenaNodeCreate(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(*providerMeta)
	var gridName = rd.Get("grid").(string)
	var nodeName = rd.Get("name").(string)

	var nodeParams = api.NodePOST{
		Name:   nodeName,
		Token:  rd.Get("token").(string),
		Labels: getKontenaNodeLabels(rd),
	}

	providerMeta.logger.Infof("Node: Create %v/%v: %#v", gridName, nodeName, nodeParams)

	if node, err := providerMeta.client.Nodes.Create(gridName, nodeParams); err != nil {
		return fmt.Errorf("Node create: %v", err)

	} else if nodeID, err := client.ParseNodeID(node.ID); err != nil {
		return fmt.Errorf("Invalid nodeID %v: %v", node.ID, err)

	} else {
		rd.SetId(nodeID.String())
		setKontenaNode(rd, node)
	}

	if err := readKontenaNodeToken(providerMeta, rd); err != nil {
		return fmt.Errorf("Node token read: %v", err)
	}

	return nil
}

func resourceKontenaNodeRead(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(*providerMeta)

	if nodeID, err := client.ParseNodeID(rd.Id()); err != nil {
		return fmt.Errorf("Invalid node ID %#v: %v", rd.Id(), err)

	} else if node, err := providerMeta.client.Nodes.Get(nodeID); err == nil {
		providerMeta.logger.Infof("Node %v: Read: %#v", rd.Id(), node)

		setKontenaNode(rd, node)

	} else if _, ok := err.(client.NotFoundError); ok {
		providerMeta.logger.Infof("Node %v: Read gone", rd.Id())

		rd.SetId("")

	} else {
		providerMeta.logger.Warnf("Node %v: Read error: %v", rd.Id(), err)

		return err
	}

	if err := readKontenaNodeToken(providerMeta, rd); err != nil {
		return fmt.Errorf("Node token read: %v", err)
	}

	return nil
}

// TODO: token update
func resourceKontenaNodeUpdate(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(*providerMeta)
	var nodeParams = api.NodePUT{}

	if rd.HasChange("labels") {
		nodeParams.Labels = getKontenaNodeLabels(rd)
	}

	providerMeta.logger.Infof("Node %v: Update %#v", rd.Id(), nodeParams)

	if nodeID, err := client.ParseNodeID(rd.Id()); err != nil {
		return fmt.Errorf("Invalid node ID %#v: %v", rd.Id(), err)
	} else if node, err := providerMeta.client.Nodes.Update(nodeID, nodeParams); err != nil {
		return err
	} else {
		setKontenaNode(rd, node)
	}

	return nil
}

func resourceKontenaNodeDelete(rd *schema.ResourceData, meta interface{}) error {
	var providerMeta = meta.(*providerMeta)

	providerMeta.logger.Infof("Node %v: Delete", rd.Id())

	if nodeID, err := client.ParseNodeID(rd.Id()); err != nil {
		return fmt.Errorf("Invalid node ID %#v: %v", rd.Id(), err)
	} else if err := providerMeta.client.Nodes.Delete(nodeID); err != nil {
		return err
	}

	return nil
}
