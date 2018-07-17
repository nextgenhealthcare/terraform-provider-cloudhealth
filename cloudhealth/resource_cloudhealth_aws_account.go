package cloudhealth

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nextgenhealthcare/cloudhealth-sdk-go"
	"strconv"
)

func resourceCloudHealthAwsAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudHealthAwsAccountCreate,
		Read:   resourceCloudHealthAwsAccountRead,
		Update: resourceCloudHealthAwsAccountUpdate,
		Delete: resourceCloudHealthAwsAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"authentication": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"access_key": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"authentication.assume_role_arn", "authentication.assume_role_external_id"},
						},
						"secret_key": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"authentication.assume_role_arn", "authentication.assume_role_external_id"},
						},
						"assume_role_arn": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"authentication.access_key", "authentication.secret_key"},
						},
						"assume_role_external_id": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"authentication.access_key", "authentication.secret_key"},
						},
					},
				},
			},
		},
	}
}

func resourceCloudHealthAwsAccountCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudhealth.Client)

	account, err := client.CreateAwsAccount(cloudhealth.AwsAccount{
		Name: d.Get("name").(string),
		Authentication: cloudhealth.AwsAccountAuthentication{
			Protocol:             d.Get("authentication.0.protocol").(string),
			AssumeRoleArn:        d.Get("authentication.0.assume_role_arn").(string),
			AssumeRoleExternalID: d.Get("authentication.0.assume_role_external_id").(string),
		},
	})
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(account.ID))

	return resourceCloudHealthAwsAccountUpdate(d, m)
}

func resourceCloudHealthAwsAccountRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudhealth.Client)

	id, _ := strconv.Atoi(d.Id())
	account, err := client.GetAwsAccount(id)
	if err == cloudhealth.ErrAwsAccountNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	d.Set("name", account.Name)
	auth := make(map[string]interface{})
	authList := make([]map[string]interface{}, 0, 1)
	auth["protocol"] = account.Authentication.Protocol
	auth["assume_role_arn"] = account.Authentication.AssumeRoleArn
	auth["assume_role_external_id"] = account.Authentication.AssumeRoleExternalID
	authList = append(authList, auth)
	d.Set("authentication", authList)

	return nil
}

func resourceCloudHealthAwsAccountUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudhealth.Client)

	id, _ := strconv.Atoi(d.Id())
	account := cloudhealth.AwsAccount{
		ID:   id,
		Name: d.Get("name").(string),
		Authentication: cloudhealth.AwsAccountAuthentication{
			Protocol:             d.Get("authentication.0.protocol").(string),
			AssumeRoleArn:        d.Get("authentication.0.assume_role_arn").(string),
			AssumeRoleExternalID: d.Get("authentication.0.assume_role_external_id").(string),
		},
	}

	updatedAccount, err := client.UpdateAwsAccount(account)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(updatedAccount.ID))

	return resourceCloudHealthAwsAccountRead(d, m)
}

func resourceCloudHealthAwsAccountDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudhealth.Client)

	id, _ := strconv.Atoi(d.Id())
	err := client.DeleteAwsAccount(id)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
