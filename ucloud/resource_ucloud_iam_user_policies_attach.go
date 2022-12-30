package ucloud

import (
	"fmt"
	//"strings"
	//"time"

	//"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	//"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

)

func resourceUCloudUserPoliciesAttach() *schema.Resource {
	return &schema.Resource{
		Create: resourceUCloudUserPoliciesAttachCreate,
		Read:   resourceUCloudUserPoliciesAttachRead,
		Update: resourceUCloudUserPoliciesAttachUpdate,
		Delete: resourceUCloudUserPoliciesAttachDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:         schema.TypeString,
				Required:     true,
			},
			"policy_urns": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				ForceNew: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				Default: "Specified",
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceUCloudUserPoliciesAttachCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.uaccountconn

	req := conn.NewGenericRequest()
	err := req.SetPayload(map[string]interface{}{
		"Action": "AttachPoliciesToUser",
		"UserName": d.Get("user_name").(string),
		"Scope": d.Get("scope").(string),
		"ProjectID": d.Get("project_id").(string),
		"PolicyURNs": d.Get("policy_urns").(*schema.Set).List(),
	})

	_, err = conn.GenericInvoke(req)
	if err != nil {
		return fmt.Errorf("error on creating iam user project attach, %s", err)
	}
	d.SetId(d.Get("user_name").(string) + "+" + d.Get("project_id").(string))

	return resourceUCloudUserPoliciesAttachRead(d, meta)
}

func resourceUCloudUserPoliciesAttachUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceUCloudUserPoliciesAttachRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceUCloudUserPoliciesAttachDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
