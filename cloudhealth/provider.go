package cloudhealth

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/nextgenhealthcare/cloudhealth-sdk-go"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "API Key from CloudHealth.",
				DefaultFunc: schema.EnvDefaultFunc("CLOUDHEALTH_API_KEY", nil),
			},
			"timeout": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "API Timeout (in seconds).",
				DefaultFunc: schema.EnvDefaultFunc("CLOUDHEALTH_API_TIMEOUT", 15),
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "API URL",
				DefaultFunc: schema.EnvDefaultFunc("CLOUDHEALTH_API_URL", "https://chapi.cloudhealthtech.com/v1/"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"cloudhealth_aws_external_id": dataSourceCloudHealthAwsExternalId(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"cloudhealth_aws_account": resourceCloudHealthAwsAccount(),
			"cloudhealth_perspective": resourceCloudHealthPerspective(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return cloudhealth.NewClient(
		d.Get("api_key").(string),
		d.Get("url").(string),
		d.Get("timeout").(int),
	)
}
