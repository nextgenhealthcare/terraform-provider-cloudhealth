package cloudhealth

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/nextgenhealthcare/cloudhealth-sdk-go"
)

func TestAccCloudHealthPerspective_basic(t *testing.T) {
	perspectiveName := fmt.Sprintf("perspective-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudHealthPerspectiveDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudHealthPerspectiveWithDefaults(perspectiveName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudHealthPerspectiveExists("cloudhealth_perspective.acc_test_perspective"),
				),
			},
		},
	})
}

func testAccCheckCloudHealthPerspectiveExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*cloudhealth.Client)
		for _, r := range s.RootModule().Resources {
			i := r.Primary.ID
			if _, err := client.GetPerspective(i); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckCloudHealthPerspectiveDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cloudhealth.Client)

	for _, r := range s.RootModule().Resources {
		i := r.Primary.ID
		if _, err := client.GetPerspective(i); err != nil {
			if err == cloudhealth.ErrPerspectiveNotFound {
				continue
			}
			return err
		}
		return fmt.Errorf("Perspective still exists")
	}
	return nil
}

func testAccCloudHealthPerspectiveWithDefaults(r string) string {
	return fmt.Sprintf(`
resource "cloudhealth_perspective" "acc_test_perspective" {
  name               = "%s"
  include_in_reports = false

  group {
    name = "OwnerAccTest"
    type = "categorize"

    rule {
      asset     = "AwsAsset"
      tag_field = ["owner"]
    }
  }
}
`, r)
}
