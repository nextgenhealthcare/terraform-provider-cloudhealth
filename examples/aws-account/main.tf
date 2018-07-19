variable "api_key" {}
variable "aws_account_name" {}

provider "cloudhealth" {
  api_key = "${var.api_key}"
}

data "cloudhealth_aws_external_id" "nextgen" {}

module "cloudhealth_iam_role" {
  source = "github.com/CloudHealth/terraform-cloudhealth-iam/role"

  role-name   = "CloudHealth"
  external-id = "${data.cloudhealth_aws_external_id.nextgen.id}"
}

resource "cloudhealth_aws_account" "main" {
  name = "${var.aws_account_name}"

  authentication {
    protocol                = "assume_role"
    assume_role_arn         = "${module.cloudhealth_iam_role.cloudhealth-role-arn}"
    assume_role_external_id = "${data.cloudhealth_aws_external_id.nextgen.id}"
  }
}
