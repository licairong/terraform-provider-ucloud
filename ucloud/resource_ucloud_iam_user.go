package ucloud

import (
	"fmt"
	//"strings"
	//"time"

	//"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	//"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

)

func resourceUCloudUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUCloudUserCreate,
		Read:   resourceUCloudUserRead,
		Update: resourceUCloudUserUpdate,
		Delete: resourceUCloudUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:         schema.TypeString,
				Required:     true,
			},
			"access_key_status": {
				Type: schema.TypeString,
				Optional: true,
				Default: "Active",
			},
			"login_profile_status": {
				Type:     schema.TypeString,
				Optional: true,
				Default: "Inactive",
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_key_secret": {
				Type:     schema.TypeString,
				Computed: true,
				Sensitive: true,
			},
		},
	}
}

func resourceUCloudUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.uaccountconn

	req := conn.NewGenericRequest()
	err := req.SetPayload(map[string]interface{}{
		"Action": "CreateUser",
		"UserName": d.Get("user_name").(string),
		"AccessKeyStatus": d.Get("access_key_status").(string),
		"LoginProfileStatus": d.Get("login_profile_status").(string),
		"DisplayName": d.Get("display_name").(string),
		"Email": d.Get("email").(string),
	})

	resp, err := conn.GenericInvoke(req)
	if err != nil {
		return fmt.Errorf("error on creating iam user, %s", err)
	}
	d.SetId(d.Get("user_name").(string))
	d.Set("access_key_id", resp.GetPayload()["AccessKeyID"])
	d.Set("access_key_secret", resp.GetPayload()["AccessKeySecret"])

	return resourceUCloudUserRead(d, meta)
}

func resourceUCloudUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.uaccountconn

	d.Partial(true)

	data := map[string]interface{}{
		"Action": "UpdateUser",
		"UserName": d.Id(),
	}

	isChanged := false

	if d.HasChange("display_name") {
		data["DisplayName"] = d.Get("display_name")
		isChanged = true
	}

	if d.HasChange("status") {
		data["Status"] = d.Get("status")
		isChanged = true
	}

	if d.HasChange("user_name") {
		data["NewUserName"] = d.Get("user_name")
		isChanged = true
	}

	req := conn.NewGenericRequest()
	_ = req.SetPayload(data)

	if isChanged {
		_, err := conn.GenericInvoke(req)
		if err != nil {
			return fmt.Errorf("error on update user %s, %s", d.Id(), err)
		}

		d.SetPartial("status")
		d.SetPartial("display_name")
		d.SetId(d.Get("user_name").(string))
	}

	d.Partial(false)

	return resourceUCloudUserRead(d, meta)
}

func resourceUCloudUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.uaccountconn

	req := conn.NewGenericRequest()
	err := req.SetPayload(map[string]interface{}{
		"Action": "ListUsers",
	})

	resp, err := conn.GenericInvoke(req)
	if err != nil {
		return fmt.Errorf("error on reading iam user %q, %s", d.Id(), err)
	}

	for _, u := range resp.GetPayload()["Users"].([]interface{}) {
		if d.Id() == u.(map[string]interface{})["UserName"].(string) {
			d.Set("email", u.(map[string]interface{})["Email"].(string))
			d.Set("display_name", u.(map[string]interface{})["DisplayName"])
			d.Set("status", u.(map[string]interface{})["Status"])
			d.Set("create_time", timestampToString(int(u.(map[string]interface{})["CreatedAt"].(float64))))
			return nil
		}
	}

	d.SetId("")

	return nil
}

func resourceUCloudUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.uaccountconn

	req := conn.NewGenericRequest()
	err := req.SetPayload(map[string]interface{}{
		"Action": "DeleteUser",
		"UserName": d.Id(),
	})

	_, err = conn.GenericInvoke(req)
	if err != nil {
		return err
	}
	return nil
}
