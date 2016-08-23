#install glusterfs repo for stable version
mkdir -p /mnt/{gluster,efs}
wget -P /etc/yum.repos.d https://download.gluster.org/pub/gluster/glusterfs/3.7/LATEST/EPEL.repo/glusterfs-epel.repo
#fix it for amazon linux
sed -i 's/$releasever/6/g' /etc/yum.repos.d/glusterfs-epel.repo
#install
yum install -y glusterfs-fuse fio