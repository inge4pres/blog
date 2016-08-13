/*
	General Naming Convention
	variables typical of AWS will start with aws_
	variables typical of AVS will start with avs_
	resources will be addressed with resource_tpye.type_subtype_attribute[_attribute]+
	
*/

# Configure the AWS Provider
provider "aws" {
    access_key = "${var.aws_access_key}"
    secret_key = "${var.aws_secret_key}"
    region = "${var.aws_region}"
}