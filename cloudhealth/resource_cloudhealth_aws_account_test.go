package cloudhealth

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/nextgenhealthcare/cloudhealth-sdk-go"
)

func TestAccCloudHealthAwsAccount_basic(t *testing.T) {
	accountName := fmt.Sprintf("account-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudHealthAwsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudHealthAwsAccountWithDefaults(accountName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudHealthAwsAccountExists("cloudhealth_aws_account.account"),
				),
			},
		},
	})
}

func testAccCheckCloudHealthAwsAccountExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*cloudhealth.Client)
		for _, r := range s.RootModule().Resources {
			i, _ := strconv.Atoi(r.Primary.ID)
			if _, err := client.GetAwsAccount(i); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckCloudHealthAwsAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cloudhealth.Client)

	for _, r := range s.RootModule().Resources {
		i, _ := strconv.Atoi(r.Primary.ID)
		if _, err := client.GetAwsAccount(i); err != nil {
			if err == cloudhealth.ErrAwsAccountNotFound {
				continue
			}
			return err
		}
		return fmt.Errorf("AWS Account still exists")
	}
	return nil
}

func testAccCloudHealthAwsAccountWithDefaults(r string) string {
	return fmt.Sprintf(`
resource "cloudhealth_aws_account" "account" {
  name = "%s"
  authentication {
    protocol = "access_key"
  }
}
`, r)
}
