resource "aws_efs_file_system" "efs" {
  performance_mode = "generalPurpose"
  tags {
    Name = "efs-benchmark"
	TFmanaged = "true"
  }
}

resource "aws_efs_mount_target" "efs_mount" {
  file_system_id = "${aws_efs_file_system.efs.id}"
  subnet_id = "${aws_subnet.public_a.id}"
  security_groups = ["${aws_security_group.sg_nfs.id}"]
}