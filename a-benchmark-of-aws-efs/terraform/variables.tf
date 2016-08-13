variable "aws_access_key" {
  default = ""
}

variable "aws_secret_key" {
  default = ""
}

variable "aws_region" {
	type = "string"
	description = "AWS region"
	default = "us-east-1"
}

variable "efs_vpc" {
	description = "VPC-CIDR"
	default = {
		cidr = "172.20.0.0/16"
		name = "efs-benchmark"
	}
}

/* Cannot interpolate in variables.tf */

variable "efs_vpc_az" {
  default = {
    a = "us-east-1b"
    b = "us-east-1c"
    c = "us-east-1d"
  }
}

variable "efs_vpc_subnet_cidr" {
  type = "map"
  description = "Map of CIDR for private/public subnets"
  default = {
	public_a = "172.20.10.0/24"
	public_b = "172.20.11.0/24"
	public_c = "172.20.100.0/24"
  }
}

variable "key_name" {
  default = "inge4pres"
}

variable "amzn_linux_ami" {
  type = "map"
  default = {
    us-east-1 = "ami-6869aa05"
	us-west-1 = "ami-31490d51"
	us-west-2 = "ami-7172b611"
	eu-west-1 = "ami-f9dd458a"
	eu-central-1 = "ami-ea26ce85"
  }
}