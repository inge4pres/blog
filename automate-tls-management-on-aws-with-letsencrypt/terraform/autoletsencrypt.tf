provider "aws" {
    access_key = ""
    secret_key = ""
    region = "eu-central-1"
}

data "template_file" "ud" {
    template = "${file("files/user-data.tmpl")}"

    vars {
		email = "${var.email}"
        domain  = "${var.domain}"
        san = "${var.san}"
		cf_id = "${var.cf_distribution_id}"
		region = "${data.aws_region.awsreg.name}"
    }
}

data "aws_region" "awsreg" {
  current = true
}

resource "aws_iam_role" "tlslego" {
	name = "autotls-lego"
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

resource "aws_iam_role_policy" "tlslego_policy" {
    name = "autotls-lego"
    role = "${aws_iam_role.tlslego.id}"
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
		"iam:ListServerCertificates",
		"apigateway:*",
		"route53:ChangeResourceRecordSets",
		"route53:Get*",
		"route53:List*",
		"s3:Get*",
		"s3:Put*",
		"s3:List*",
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

resource "aws_iam_instance_profile" "tlslego_profile" {
    name = "autotls-lego"
    roles = ["${aws_iam_role.tlslego.name}"]
}

resource "aws_s3_bucket" "legostore" {
	bucket = "${var.domain}-lego-account"
    acl = "private"

    tags {
        Domain = "${var.domain}"
        Environment = "prod"
    }
	
	provisioner "local-exec" {
        command = "aws s3 sync ${var.lego-local-path}/.lego s3://${var.domain}-lego-account/"
    }
}

resource "aws_launch_configuration" "updatetls_lc" {
    name_prefix = "tf-autotls-lego"
    image_id = "${lookup(var.amzn_linux_ami,data.aws_region.awsreg.name)}"
    instance_type = "t2.small"
	iam_instance_profile = "${aws_iam_instance_profile.tlslego_profile.name}"
	key_name = "${var.ssh-key-name}"
	user_data = "${data.template_file.ud.rendered}"
	
	lifecycle {
      create_before_destroy = true
    }
}

resource "aws_autoscaling_group" "updatetls_as" {
  availability_zones = ["${data.aws_region.awsreg.name}a"]
  name = "${var.domain}-auto-tls"
  max_size = 1
  min_size = 0
  health_check_grace_period = 240
  health_check_type = "EC2"
  desired_capacity = 0
  force_delete = true
  launch_configuration = "${aws_launch_configuration.updatetls_lc.name}"

  lifecycle {
    create_before_destroy = true
  }
  tag {
    key = "Name"
    value = "${var.domain}-autotls-lego"
    propagate_at_launch = true
  }
  tag {
    key = "Domain"
    value = "${var.domain}"
    propagate_at_launch = true
  }
  tag {
    key = "SAN"
    value = "${var.san}"
    propagate_at_launch = true
  }
  tag {
    key = "TLS"
    value = "letsencrypt"
    propagate_at_launch = false
  }
}

resource "aws_autoscaling_schedule" "autotls_up" {
    scheduled_action_name = "autotls-up"
    desired_capacity = 1
	max_size = 1
	min_size = 1
	//run every 15th day of month
	recurrence = "0 18 15 * *"
    autoscaling_group_name = "${aws_autoscaling_group.updatetls_as.name}"
}

resource "aws_autoscaling_schedule" "autotls_down" {
    scheduled_action_name = "autotls-down"
	max_size = 0
	min_size = 0
    desired_capacity = 0
	//run every 15th day of month
	recurrence = "24 18 15 * *"
    autoscaling_group_name = "${aws_autoscaling_group.updatetls_as.name}"
}

