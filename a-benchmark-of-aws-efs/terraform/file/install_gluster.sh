#!/bin/bash
yum update -y
yum install -y ntp
/sbin/service ntpd restart
#install xfs tools
yum install -y xfsprogs
#format new disk
mkfs.xfs -i size=512 /dev/sdf
#setup mount point 
mkdir -p /export/gluster
#setup fstab for brick
echo "/dev/sdf /export/gluster xfs defaults 1 2" >> /etc/fstab
#mount
mount -a