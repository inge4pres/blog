output "client-ip" {
  value = "${aws_instance.client.dns_name}"
}

output "server-a-ip" {
  value = "${aws_instance.server_a.private_ip}"
}

output "server-b-ip" {
  value = "${aws_instance.server_b.private_ip}"
}

output "efs-mount-target-a" {
  value = "${aws_efs_mount_target.efs_mount.dns_name}"
}