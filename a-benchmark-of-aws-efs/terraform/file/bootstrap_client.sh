#install glusterfs repo for stable version
sudo mkdir -p /mnt/{gluster,efs}
sudo wget -P /etc/yum.repos.d https://download.gluster.org/pub/gluster/glusterfs/3.6/LATEST/EPEL.repo/glusterfs-epel.repo
#fix it for amazon linux
sudo sed -i 's/$releasever/6/g' /etc/yum.repos.d/glusterfs-epel.repo
#install
sudo yum install -y fio glusterfs-fuse