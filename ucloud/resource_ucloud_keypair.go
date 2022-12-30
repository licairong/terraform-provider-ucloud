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

func resourceUCloudKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceUCloudKeyPairCreate,
		Read:   resourceUCloudKeyPairRead,
		Update: resourceUCloudKeyPairUpdate,
		Delete: resourceUCloudKeyPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key_pair_name": {
				Type:         schema.TypeString,
				Required:     true,
			},
			"private_key": {
				Type: schema.TypeString,
				Computed: true,
				Sensitive: true,
			},
			"key_pair_finger_print": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceUCloudKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.uhostconn

	req := conn.NewCreateUHostKeyPairRequest()
	req.KeyPairName = ucloud.String(d.Get("key_pair_name").(string))

	resp, err := conn.CreateUHostKeyPair(req)
	if err != nil {
		return fmt.Errorf("error on creating keypair, %s", err)
	}
	d.SetId(resp.KeyPair.KeyPairId)
	d.Set("private_key", resp.KeyPair.PrivateKeyBody)
	d.Set("key_pair_finger_print", resp.KeyPair.KeyPairFingerPrint)

	return resourceUCloudKeyPairRead(d, meta)
}

func resourceUCloudKeyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceUCloudKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.uhostconn

	req := conn.NewDescribeUHostKeyPairsRequest()
	req.KeyPairFingerPrint = ucloud.String(d.Get("key_pair_finger_print").(string))

	resp, err := conn.DescribeUHostKeyPairs(req)
	if err != nil {
		return fmt.Errorf("error on reading keypair %q, %s", d.Id(), err)
	}
	if len(resp.KeyPairs) < 1 {
		d.SetId("")
		return nil
	}

	d.Set("key_pair_name", resp.KeyPairs[0].KeyPairName)
	d.Set("key_pair_finger_print", resp.KeyPairs[0].KeyPairFingerPrint)
	d.Set("create_time", timestampToString(resp.KeyPairs[0].CreateTime))

	return nil
}

func resourceUCloudKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.uhostconn

	req := conn.NewDeleteUHostKeyPairsRequest()
	req.KeyPairIds = []string{
		d.Id(),
	}

	_, err := conn.DeleteUHostKeyPairs(req)
	if err != nil {
		return err
	}
	return nil
}
