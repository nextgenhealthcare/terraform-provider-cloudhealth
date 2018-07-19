## Enabling an AWS Account in CloudHealth

This example provides sample configuration for enabling an AWS Account. It configures the IAM Role in AWS for CloudHealth and enables the AWS Account in CloudHealth. In addition to having a CloudHealth account, you'll also need an AWS Account and credentials. The example does not specify AWS credential configuration and will either use the _default_ profile or environment variables.

Once ready run `terraform plan -out example.plan` to review.

You will be prompted to provide input for the following variables:

* api_key: API Key from CloudHealth. For more information, see [Getting Your API Key](http://apidocs.cloudhealthtech.com/#documentation_getting-your-api-key).
* aws_account_name: Provide a name for the AWS Account.

Once satisfied with plan, run `terraform apply example.plan`
