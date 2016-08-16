output "client-ip" {
  value = "${aws_instance.client.dns_name}"
}

output "gluster-a-ip" {
  value = "${aws_instance.gluster_a.private_ip}"
}

output "gluster-b-ip" {
  value = "${aws_instance.gluster_b.private_ip}"
}

output "efs-mount-target-a" {
  value = "${aws_efs_mount_target.efs_mount.dns_name}"
}