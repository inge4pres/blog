#!/bin/bash
yum update -y
yum install -y ntp
/sbin/service ntpd restart
#install xfs tools
yum install -y xfsprogs
#create one partion from disk
printf "o\nn\np\n1\n\n\nw\n" | fdisk /dev/sdf
#format new partition
mkfs.xfs -i size=512 /dev/sdf1
#setup mount point 
mkdir -p /export/gluster
#mount
mount -t xfs /dev/sdf1 /export/gluster
mkdir -p /export/gluster/brick
#install glusterfs repo for stable version
wget -P /etc/yum.repos.d https://download.gluster.org/pub/gluster/glusterfs/3.7/LATEST/EPEL.repo/glusterfs-epel.repo
#fix it for amazon linux
sed -i 's/$releasever/6/g' /etc/yum.repos.d/glusterfs-epel.repo
#install
yum install -y glusterfs{-fuse,-server}