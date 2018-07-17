package cloudhealth

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nextgenhealthcare/cloudhealth-sdk-go"
)

type ExternalID struct {
	ExternalID string `json:"generated_external_id,omitempty"`
}

func dataSourceCloudHealthAwsExternalId() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsOrganizationsOrganizationRead,

		Schema: map[string]*schema.Schema{
			"external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAwsOrganizationsOrganizationRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudhealth.Client)

	id, err := client.GetAwsExternalID()
	if err != nil {
		return err
	}

	d.Set("external_id", id)
	d.SetId(id)

	return nil
}
