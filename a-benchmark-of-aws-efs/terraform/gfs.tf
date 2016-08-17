resource "aws_instance" "server_a" {
    ami = "${var.amzn_linux_ami["${var.aws_region}"]}"
    instance_type = "c4.large"
	availability_zone = "${var.efs_vpc_az["a"]}"
	subnet_id = "${aws_subnet.public_a.id}"
	ebs_optimized = true
	monitoring = false
	vpc_security_group_ids  = ["${aws_security_group.gluster_server.id}","${aws_security_group.sg_nfs.id}"]
	user_data = "${file("file/bootstrap_server.sh")}"
	key_name = "${var.key_name}"
	
	ebs_block_device {
		device_name = "/dev/sdf"
		volume_type = "gp2"
		volume_size = "10"
		delete_on_termination = true
	}

	tags {
		Name = "server_1"
    	TFmanaged = "true"
		Server_id = 1
	}
}

resource "aws_instance" "server_b" {
    ami = "${var.amzn_linux_ami["${var.aws_region}"]}"
    instance_type = "c4.large"
	availability_zone = "${var.efs_vpc_az["b"]}"
	subnet_id = "${aws_subnet.public_b.id}"
	ebs_optimized = true
	monitoring = false
	vpc_security_group_ids  = ["${aws_security_group.gluster_server.id}","${aws_security_group.sg_nfs.id}"]
	user_data = "${file("file/bootstrap_server.sh")}"
	key_name = "${var.key_name}"
	
	ebs_block_device {
		device_name = "/dev/sdf"
		volume_type = "gp2"
		volume_size = "10"
		delete_on_termination = true
	}

	tags {
		Name = "server_2"
    	TFmanaged = "true"
		Server_id = 2
	}
}