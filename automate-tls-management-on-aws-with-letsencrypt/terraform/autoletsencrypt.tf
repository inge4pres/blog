provider "aws" {
    access_key = ""
    secret_key = ""
    region = "eu-central-1"
}

resource "aws_iam_role" "updatetls" {
	name = "tls-letsencrypt-update"
	assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": "AllowEc2"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "updatetls_policy" {
    name = "tls-letsencrypt-update-policy"
    role = "${aws_iam_role.updatetls.id}"
    policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
	    "cloudfront:GetDistributionConfig",
        "cloudfront:UpdateDistribution",
		"iam:UploadServerCertificate",
		"iam:DeleteServerCertificate",
		"apigateway:*",
		"route53:ChangeResourceRecordSets",
		"route53:Get*",
		"route53:List*",
		"ses:Send*Email",
		"ses:List*",
		"ses:Get*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_autoscaling_group" "updatetls_as" {
  availability_zones = ["eu-central-1"]
  name = "4pres-auto-tls"
  max_size = 1
  min_size = 1
  health_check_grace_period = 240
  health_check_type = "EC2"
  desired_capacity = 1
  force_delete = true
  launch_configuration = "${aws_launch_configuration.updatetls_lc.name}"

  lifecycle {
    create_before_destroy = true
  }
	
  tag {
    key = "Domain"
    value = "4pres"
    propagate_at_launch = true
  }
  tag {
    key = "TLS"
    value = "Letsencrypt"
    propagate_at_launch = false
  }
}

resource "aws_launch_configuration" "updatetls_lc" {
    name_prefix = "tf-autotls-"
    name = "auto-tls"
    image_id = "${var.packer-ami}"
    instance_type = "t2.small"
	iam_instance_profile = "${aws_iam_role.updatetls.name}"
	key_name = "${var.ssh-key-name}"
	associate_public_ip_address = true
	
	lifecycle {
      create_before_destroy = true
    }
}


















