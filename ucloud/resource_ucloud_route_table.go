package ucloud

import (
	"fmt"
	//"strings"
	//"time"

	//"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	//"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	//"github.com/ucloud/ucloud-sdk-go/services/uhost"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

func resourceUCloudRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceUCloudRouteTableCreate,
		Read:   resourceUCloudRouteTableRead,
		Update: resourceUCloudRouteTableUpdate,
		Delete: resourceUCloudRouteTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
			},
			"vpc_id": {
				Type: schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"route_table_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceUCloudRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.vpcconn

	req := conn.NewCreateRouteTableRequest()
	req.ProjectId = ucloud.String(d.Get("project_id").(string))
	req.VPCId = ucloud.String(d.Get("vpc_id").(string))
	req.Name = ucloud.String(d.Get("name").(string))
	req.Tag = ucloud.String(d.Get("tag").(string))
	req.Remark = ucloud.String(d.Get("remark").(string))

	resp, err := conn.CreateRouteTable(req)
	if err != nil {
		return fmt.Errorf("error on creating route table, %s", err)
	}
	d.SetId(resp.RouteTableId)

	for _, subnet_id := range d.Get("subnet_ids").([]interface{}) {
		r := conn.NewAssociateRouteTableRequest()
		r.ProjectId = ucloud.String(d.Get("project_id").(string))
		r.SubnetId = ucloud.String(subnet_id.(string))
		r.RouteTableId = ucloud.String(d.Id())

		_, err := conn.AssociateRouteTable(r)
		if err != nil {
			return err
		}
	}

	itemsRaw := d.Get("rules").([]interface{})
	if len(itemsRaw) > 0 {
		r := conn.NewModifyRouteRuleRequest()
		r.ProjectId = ucloud.String(d.Get("project_id").(string))
		r.RouteTableId = ucloud.String(d.Id())
		items := make([]string, len(itemsRaw))
		for k, raw := range itemsRaw {
			items[k] = raw.(string)
		}
		r.RouteRule = items

		_, err := conn.ModifyRouteRule(r)
		if err != nil {
			return err
		}
	}

	return resourceUCloudRouteTableRead(d, meta)
}

func resourceUCloudRouteTableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.vpcconn

	d.Partial(true)

	req := conn.NewUpdateRouteTableAttributeRequest()
	req.ProjectId = ucloud.String(d.Get("project_id").(string))
	req.RouteTableId = ucloud.String(d.Id())

	isChanged := false

	if d.HasChange("name") {
		req.Name = ucloud.String(d.Get("name").(string))
		isChanged = true
	}

	if d.HasChange("remark") {
		req.Remark = ucloud.String(d.Get("remark").(string))
		isChanged = true
	}

	if d.HasChange("tag") {
		req.Tag = ucloud.String(d.Get("tag").(string))
		isChanged = true
	}

	if isChanged {
		_, err := conn.UpdateRouteTableAttribute(req)
		if err != nil {
			return fmt.Errorf("error on update user %s, %s", d.Id(), err)
		}

		d.SetPartial("name")
		d.SetPartial("remark")
		d.SetPartial("tag")
	}

	d.Partial(false)

	for _, subnet_id := range d.Get("subnet_ids").([]interface{}) {
		r := conn.NewAssociateRouteTableRequest()
		r.ProjectId = ucloud.String(d.Get("project_id").(string))
		r.SubnetId = ucloud.String(subnet_id.(string))
		r.RouteTableId = ucloud.String(d.Id())

		_, err := conn.AssociateRouteTable(r)
		if err != nil {
			return err
		}
	}

	itemsRaw := d.Get("rules").([]interface{})
	if len(itemsRaw) > 0 {
		r := conn.NewModifyRouteRuleRequest()
		r.ProjectId = ucloud.String(d.Get("project_id").(string))
		r.RouteTableId = ucloud.String(d.Id())
		items := make([]string, len(itemsRaw))
		for k, raw := range itemsRaw {
			items[k] = raw.(string)
		}
		r.RouteRule = items

		_, err := conn.ModifyRouteRule(r)
		if err != nil {
			return err
		}
	}

	return resourceUCloudRouteTableRead(d, meta)
}

func resourceUCloudRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.vpcconn

	req := conn.NewDescribeRouteTableRequest()
	req.ProjectId = ucloud.String(d.Get("project_id").(string))
	req.RouteTableId = ucloud.String(d.Id())

	resp, err := conn.DescribeRouteTable(req)
	if err != nil {
		return fmt.Errorf("error on reading keypair %q, %s", d.Id(), err)
	}
	if resp.TotalCount < 1 {
		d.SetId("")
		return nil
	}

	d.Set("tag", resp.RouteTables[0].Tag)
	d.Set("remark", resp.RouteTables[0].Remark)
	d.Set("vpc_id", resp.RouteTables[0].VPCId)
	d.Set("vpc_name", resp.RouteTables[0].VPCName)
	d.Set("subnet_count", resp.RouteTables[0].SubnetCount)
	d.Set("route_table_type", resp.RouteTables[0].RouteTableType)
	d.Set("create_time", timestampToString(resp.RouteTables[0].CreateTime))

	return nil
}

func resourceUCloudRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.vpcconn

	for _, subnet_id := range d.Get("subnet_ids").([]interface{}) {
		r := conn.NewDescribeRouteTableRequest()
		r.ProjectId = ucloud.String(d.Get("project_id").(string))
		r.VPCId = ucloud.String(d.Get("vpc_id").(string))

		rs, err := conn.DescribeRouteTable(r)
		if err != nil {
			return err
		}

		default_route_table_id := ""
		for _, route_table := range rs.RouteTables {
			if route_table.Tag == "Default" {
				default_route_table_id = route_table.RouteTableId
				break
			}
		}

		r2 := conn.NewAssociateRouteTableRequest()
		r2.ProjectId = ucloud.String(d.Get("project_id").(string))
		r2.SubnetId = ucloud.String(subnet_id.(string))
		r2.RouteTableId = ucloud.String(default_route_table_id)

		_, err = conn.AssociateRouteTable(r2)
		if err != nil {
			return err
		}
	}

	req := conn.NewDeleteRouteTableRequest()
	req.ProjectId = ucloud.String(d.Get("project_id").(string))
	req.RouteTableId = ucloud.String(d.Id())

	_, err := conn.DeleteRouteTable(req)
	if err != nil {
		return err
	}
	return nil
}
