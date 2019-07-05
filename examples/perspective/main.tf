variable "api_key" {}

provider "cloudhealth" {
  api_key = "${var.api_key}"
}

resource "cloudhealth_perspective" "my_perspective" {
  name = "My Perspective"
  include_in_reports = false

  group {
    name = "My Team"
    type = "filter"

    rule {
      asset = "AwsAsset"
      condition {
        tag_field = ["team"]
        val = "my_team"
      }
    }

    rule {
      asset = "AwsAsset"
      condition {
        tag_field = ["team"]
        val = "my_team@corp.com"
      }
    }
  }

  group {
    name = "redshift"
    type = "categorize"

    rule {
      asset = "AwsRedshiftCluster"
      field = ["Cluster Identifier"]
    }
  }
}
