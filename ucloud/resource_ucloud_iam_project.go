package ucloud

import (
	"fmt"
	"github.com/ucloud/ucloud-sdk-go/ucloud"

	//"strings"
	//"time"

	//"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	//"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

)

func resourceUCloudProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceUCloudProjectCreate,
		Read:   resourceUCloudProjectRead,
		Update: resourceUCloudProjectUpdate,
		Delete: resourceUCloudProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
			},
			"copy_image_id": {
				Type:         schema.TypeString,
				Optional:     true,
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
			},
			"default_project_id": {
				Type:         schema.TypeString,
				Optional:     true,
			},
			"target_image_id": {
				Type:         schema.TypeString,
				Computed:     true,
			},
			"is_default": {
				Type:         schema.TypeBool,
				Computed:     true,
			},
			"member_count": {
				Type:         schema.TypeInt,
				Computed:     true,
			},
			"resource_count": {
				Type:         schema.TypeInt,
				Computed:     true,
			},
			"parent_id": {
				Type:         schema.TypeString,
				Computed:     true,
			},
			"parent_name": {
				Type:         schema.TypeString,
				Computed:     true,
			},
			"create_time": {
				Type:         schema.TypeString,
				Computed:     true,
			},
		},
	}
}

func resourceUCloudProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.uaccountconn

	req := conn.NewGenericRequest()
	err := req.SetPayload(map[string]interface{}{
		"Action": "CreateProject",
		"ProjectName": d.Get("name").(string),
	})

	resp, err := conn.GenericInvoke(req)
	if err != nil {
		return fmt.Errorf("error on creating project, %s", err)
	}
	d.SetId(resp.GetPayload()["ProjectId"].(string))

	if d.Get("copy_image_id").(string) != "" {
		hostconn := client.uhostconn
		r := hostconn.NewCopyCustomImageRequest()
		r.Zone = ucloud.String(d.Get("zone").(string))
		r.ProjectId = ucloud.String(d.Get("default_project_id").(string))
		r.SourceImageId = ucloud.String(d.Get("copy_image_id").(string))
		r.TargetProjectId = ucloud.String(d.Id())

		rs, er := hostconn.CopyCustomImage(r)
		if er != nil {
			return er
		}
		d.Set("target_image_id", rs.TargetImageId)
	}

	return resourceUCloudProjectRead(d, meta)
}

func resourceUCloudProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceUCloudProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.uaccountconn

	req := conn.NewGetProjectListRequest()

	resp, err := conn.GetProjectList(req)
	if err != nil {
		return fmt.Errorf("error on reading project %q, %s", d.Id(), err)
	}

	for _, p := range resp.ProjectSet {
		if p.ProjectId == d.Id() {
			d.Set("is_default", p.IsDefault)
			d.Set("member_count", p.MemberCount)
			d.Set("resource_count", p.ResourceCount)
			d.Set("parent_id", p.ParentId)
			d.Set("parent_name", p.ParentName)
			d.Set("create_time", timestampToString(p.CreateTime))
			return nil
		}
	}

	d.SetId("")

	return nil
}

func resourceUCloudProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.uaccountconn

	req := conn.NewGenericRequest()
	err := req.SetPayload(map[string]interface{}{
		"Action": "DeleteProject",
		"ProjectID": d.Id(),
	})

	_, err = conn.GenericInvoke(req)
	if err != nil {
		return err
	}
	return nil
}
