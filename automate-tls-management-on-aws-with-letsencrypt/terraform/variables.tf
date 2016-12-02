variable "ssh-key-name" {
	default = "inge4pres"
}

variable "domain" {
	default = "4pr.es"
}

variable "san" {
	default = "inge.4pr.es"
}

variable "email" {
	default = "fgualazzi@gmail.com"
}

variable "lego-local-path" {
	default = "%userprofile%"
}

variable "cf_distribution_id" {
	default = "123"
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