/*
 Network structure defined in single VPC
*/

# ---- Base VPC named "app"
resource "aws_vpc" "app" {
	cidr_block = "${var.efs_vpc["cidr"]}"
	enable_dns_support = true
	enable_dns_hostnames = true
	
	tags {
		Name = "${var.efs_vpc["name"]}"
	}
}

# ---- Internet Gateway
resource "aws_internet_gateway" "igw" {
    vpc_id = "${aws_vpc.app.id}"

    tags {
        Name = "${var.efs_vpc["name"]}_igw"
		Scope = "public"
    }
}

# ---- Subnets
resource "aws_subnet" "public_a" {
	vpc_id = "${aws_vpc.app.id}"
	cidr_block = "${var.efs_vpc_subnet_cidr["public_a"]}"
    map_public_ip_on_launch = true
	tags {
		Name = "subnet_public_a"
		Scope = "public"
		TF_managed = "true"
	}
}

resource "aws_subnet" "public_b" {
	vpc_id = "${aws_vpc.app.id}"
	cidr_block = "${var.efs_vpc_subnet_cidr["public_b"]}"
	map_public_ip_on_launch = true
	tags {
		Name = "subnet_public_b"
		Scope = "public"
		TF_managed = "true"
	}
}
# ---- Routes - only one main public route will be used
resource "aws_route_table" "rt_main" {
	vpc_id = "${aws_vpc.app.id}"
    route {
        cidr_block = "0.0.0.0/0"
        gateway_id = "${aws_internet_gateway.igw.id}"
    }

    tags {
        Name = "route_main"
		TF_managed = "true"
    }
}
resource "aws_route_table_association" "assoc_main_a" {
    subnet_id = "${aws_subnet.public_a.id}"
    route_table_id = "${aws_route_table.rt_main.id}"
}
resource "aws_route_table_association" "assoc_main_b" {
    subnet_id = "${aws_subnet.public_b.id}"
    route_table_id = "${aws_route_table.rt_main.id}"
}

# ---- Security Groups
resource "aws_security_group" "gluster_server" {
	description  = "Gluster FS access"
	name = "glusterfs"
	vpc_id = "${aws_vpc.app.id}"
	
	tags {
		Name = "sg_glusterfs"
		Role = "efs-test"
		TF_managed = "true"
	}
	
	ingress {
		from_port = "22"
		to_port = "22"
		protocol = "tcp"
		cidr_blocks = ["${var.efs_vpc["cidr"]}"]
	}
	
	ingress {
		from_port = "111"
		to_port = "111"
		protocol = "udp"
		cidr_blocks = ["${var.efs_vpc["cidr"]}"]
	}
	
	ingress {
		from_port = "111"
		to_port = "111"
		protocol = "tcp"
		cidr_blocks = ["${var.efs_vpc["cidr"]}"]
	}
	
	ingress {
		from_port = "24007"
		to_port = "24010"
		protocol = "tcp"
		cidr_blocks = ["${var.efs_vpc["cidr"]}"]
	}
	
	ingress {
		from_port = "38465"
		to_port = "38467"
		protocol = "tcp"
		cidr_blocks = ["${var.efs_vpc["cidr"]}"]
	}
	
	ingress {
		from_port = "49152"
		to_port = "49153"
		protocol = "tcp"
		cidr_blocks = ["${var.efs_vpc["cidr"]}"]
	}

	egress {
		from_port = 0
		to_port = 0
		protocol = "-1"
		cidr_blocks = ["0.0.0.0/0"]
	}
}

resource "aws_security_group" "client" {
	description  = "Client access"
	name = "client"
	vpc_id = "${aws_vpc.app.id}"
	
	tags {
		Name = "sg_glusterfs"
		Role = "efs-test"
		TF_managed = "true"
	}
	
	ingress {
		from_port = "22"
		to_port = "22"
		protocol = "tcp"
		cidr_blocks = ["0.0.0.0/0"]
	}

	egress {
		from_port = 0
		to_port = 0
		protocol = "-1"
		cidr_blocks = ["0.0.0.0/0"]
	}
}

resource "aws_security_group" "sg_nfs" {
	description  = "EFS access"
	name = "sg_efs"
	vpc_id = "${aws_vpc.app.id}"
	
	tags {
		Name = "sg_efs"
		Role = "efs-test"
		TF_managed = "true"
	}
	
	ingress {
		from_port = "2049"
		to_port = "2049"
		protocol = "tcp"
		cidr_blocks = ["${var.efs_vpc["cidr"]}"]
	}

	egress {
		from_port = 0
		to_port = 0
		protocol = "-1"
		cidr_blocks = ["0.0.0.0/0"]
	}
}