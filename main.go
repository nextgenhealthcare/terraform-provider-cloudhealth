package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/nextgenhealthcare/terraform-provider-cloudhealth/cloudhealth"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cloudhealth.Provider,
	})
}
