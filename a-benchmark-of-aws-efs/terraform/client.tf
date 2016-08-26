resource "aws_instance" "client" {
    ami = "${var.amzn_linux_ami["${var.aws_region}"]}"
    instance_type = "m3.medium"
	availability_zone = "${var.efs_vpc_az["a"]}"
	subnet_id = "${aws_subnet.public_a.id}"
	vpc_security_group_ids  = ["${aws_security_group.client.id}"]
	user_data = "${file("file/bootstrap_client.sh")}"
	key_name = "${var.key_name}"

	tags {
		Name = "fs_client"
    	TFmanaged = "true"
	}
}