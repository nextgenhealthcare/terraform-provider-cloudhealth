package cloudhealth

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCloudHealthAwsExternalID_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudHealthAwsExternalIDConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudHealthAwsExternalIDExists("data.cloudhealth_aws_external_id.selected"),
				),
			},
		},
	})
}

func testAccCheckCloudHealthAwsExternalIDExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Can't find AWS External ID data source: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("AWS External ID data source ID not set")
		}

		return nil
	}
}

const testAccCloudHealthAwsExternalIDConfig = `
data "cloudhealth_aws_external_id" "selected" {}
`
