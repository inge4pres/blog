---
title: A benchmark of AWS EFS
author: "inge4pres"
layout: post
date: 2016-09-18T17:26:17+00:00

categories:
  - social
  - tech
tags:
  - aws
  - benchmark
  - cloud
  - EFS
  - filesystem

---
Amazon Web Services <a href="https://aws.amazon.com/efs/" target="_blank">Elastic File System</a> has been to my knowledge the service to have the longest beta testing period: reason for this may be that not as many client as expected tested it and AWS received too few feedback on it or that there were issues not to release GA. I don&#8217;t want to speculate on which one is correct but now that it has been <a href="https://aws.amazon.com/about-aws/whats-new/2016/06/amazon-elastic-file-system-efs-is-now-generally-available/" target="_blank">officially released</a> I decided to give it a try and of course compare it to a self-managed solution on the same platform.

If you followed AWS evolution you may agree that EFS has been introduced to fill the gap between EBS storage and S3: before EFS was live there was no &#8220;easy&#8221; way of having a distributed file system in AWS, you could only <a href="http://inge.4pr.es/blog/?p=318" target="_blank">set up your own</a> using  a combination of EC2 instances mounting Elastic Block Storage volumes and S3. Now with EFS you can have a AWS-managed distributed file system to be used in your cloud environment or even across the internet (will try that on a public subnet) with all the benefits of offloading the high-availability and replication burden to Amazon, and at a reasonable price. Will performance be enough compared to a self-managed solution?

##### Playground

I use <a href="https://www.terraform.io/" target="_blank">terraform</a> to create an infrastructure template to run the tests, you can see it <a href="https://github.com/inge4pres/blog/tree/master/a-benchmark-of-aws-efs/terraform" target="_blank">here</a>. Once

<pre class="">terraform apply</pre>

has finished, you&#8217;ll end up with:

  * An EFS with General Purpose performance mode
  * An EFS mount target for 1 Availability Zone
  * 1 EC2 instance named &#8220;client&#8221; to mount remote file systems
  * 2 EC2 instances named &#8220;server_X&#8221; each one with a 10 GB General Purpose EBS, they will serve a self-managed distributed, replicated file system

This is the terraform output and the steps to run on the 2 server nodes to have a running GlusterFS replicated volume; to configure the NFS export on server1, I used <a href="http://www.tldp.org/HOWTO/NFS-HOWTO/server.html" target="_blank">this guide</a>.

<pre class="theme:github lang:sh decode:true " title="Terraform output">Apply complete! Resources: 15 added, 0 changed, 0 destroyed.

The state of your infrastructure has been saved to the path
below. This state is required to modify and destroy your
infrastructure, so keep it safe. To inspect the complete state
use the `terraform show` command.

State path: terraform.tfstate

Outputs:

client-public-ip = ec2-54-89-124-238.compute-1.amazonaws.com
efs-mount-target-a = us-east-1b.fs-a6418fef.efs.us-east-1.amazonaws.com
server-a-ip = 172.30.10.122
server-b-ip = 172.30.11.31</pre>

<pre class="theme:github lang:sh decode:true " title="Create GlusterFS vol">ssh ec2-user@ec2-54-89-124-238.compute-1.amazonaws.com
# ssh to first server, need agent forwarding setup
ssh 172.30.10.122
sudo service glusterd restart
# probe the node 2
sudo gluster peer probe 172.30.11.31
sudo gluster peer status
** 
Number of Peers: 1

Hostname: 172.30.11.31
Uuid: 237bc59a-20b6-4b30-b133-abab34e36720
State: Peer in Cluster (Connected)
** 

sudo gluster volume create efs-bench replica 2 transport tcp \
   172.30.10.122:/export/gluster/brick 172.30.11.31:/export/gluster/brick force
# force required to use the root disk as export brick
** volume create: efs-bench: success: please start the volume to access data

sudo gluster volume start efs-bench
** volume start: efs-bench: success</pre>

On the client I mount the EFS target with NFS4.1, the GlusterFS volume from the server _in the same subnet_ via the GlusterFS native client_ _and the NFS export on the client&#8217;s designated mount points. I use the server on the same subnet as the client is, because the EFS target exposes a mount point in the same subnet and latency is a key factor in remote file system.

<pre class="theme:github lang:sh decode:true" title="Client mounts">sudo mount -t nfs4 -o vers=4.1 \
  us-east-1b.fs-a6418fef.efs.us-east-1.amazonaws.com:/ /mnt/efs

sudo mount -t glusterfs 172.30.10.122:/efs-bench /mnt/gluster

sudo mount -t nfs 172.30.10.122:/export/nfs /mnt/nfs -o user=ec2-user</pre>

##### Benchmark Tests

I used <a href="https://github.com/axboe/fio" target="_blank">fio</a> installed on the client box with a command suggested by this <a href="https://www.binarylane.com.au/support/solutions/articles/1000055889-how-to-benchmark-disk-i-o" target="_blank">BinaryLane post</a> and run it against the mount point for EFS, GlusterFS and NFS v4.0 with the following command

<pre class="">fio --randrepeat=1 --ioengine=libaio --direct=1 --gtod_reduce=1 --name=test \
   --filename=test --bs=4k --iodepth=64 --size=1G --readwrite=randrw --rwmixread=75</pre>

changing the target directory each time I tun the test; each test is run isolated.

I did not customize any storage option for GlusterFS or NFS, so I&#8217;m using the default options.

##### Benchmark Results

###### EFS

<pre class="theme:github lang:sh decode:true " title="EFS fio result">test: (g=0): rw=randrw, bs=4K-4K/4K-4K/4K-4K, ioengine=libaio, iodepth=64
fio-2.1.5
Starting 1 process
test: Laying out IO file(s) (1 file(s) / 1024MB)
Jobs: 1 (f=1): [m] [100.0% done] [3263KB/1078KB/0KB /s] [815/269/0 iops] [eta 00m:00s]
test: (groupid=0, jobs=1): err= 0: pid=8336: Fri Aug 26 17:30:05 2016
  read : io=784996KB, bw=3185.4KB/s, iops=796, runt=246443msec
  write: io=263580KB, bw=1069.6KB/s, iops=267, runt=246443msec
  cpu          : usr=0.31%, sys=0.58%, ctx=383567, majf=0, minf=5
  IO depths    : 1=0.1%, 2=0.1%, 4=0.1%, 8=0.1%, 16=0.1%, 32=0.1%, &gt;=64=100.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, &gt;=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.1%, &gt;=64=0.0%
     issued    : total=r=196249/w=65895/d=0, short=r=0/w=0/d=0
     latency   : target=0, window=0, percentile=100.00%, depth=64

Run status group 0 (all jobs):
   READ: io=784996KB, aggrb=3185KB/s, minb=3185KB/s, maxb=3185KB/s, mint=246443msec, maxt=246443msec
  WRITE: io=263580KB, aggrb=1069KB/s, minb=1069KB/s, maxb=1069KB/s, mint=246443msec, maxt=246443msec</pre>

###### GlusterFS

<pre class="theme:github lang:sh decode:true " title="GlusterFS results">test: (g=0): rw=randrw, bs=4K-4K/4K-4K/4K-4K, ioengine=libaio, iodepth=64
fio-2.1.5
Starting 1 process
test: Laying out IO file(s) (1 file(s) / 1024MB)
Jobs: 1 (f=1): [m] [100.0% done] [1648KB/632KB/0KB /s] [412/158/0 iops] [eta 00m:00s]
test: (groupid=0, jobs=1): err= 0: pid=2857: Sun Sep 18 16:20:46 2016
  read : io=784996KB, bw=1929.2KB/s, iops=482, runt=406921msec
  write: io=263580KB, bw=663288B/s, iops=161, runt=406921msec
  cpu          : usr=0.32%, sys=0.39%, ctx=262172, majf=0, minf=5
  IO depths    : 1=0.1%, 2=0.1%, 4=0.1%, 8=0.1%, 16=0.1%, 32=0.1%, &gt;=64=100.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, &gt;=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.1%, &gt;=64=0.0%
     issued    : total=r=196249/w=65895/d=0, short=r=0/w=0/d=0
     latency   : target=0, window=0, percentile=100.00%, depth=64

Run status group 0 (all jobs):
   READ: io=784996KB, aggrb=1929KB/s, minb=1929KB/s, maxb=1929KB/s, mint=406921msec, maxt=406921msec
  WRITE: io=263580KB, aggrb=647KB/s, minb=647KB/s, maxb=647KB/s, mint=406921msec, maxt=406921msec</pre>

###### NFS 4.0 (not replicated)

<pre class="theme:github lang:sh decode:true" title="NFS4.0 results">test: (g=0): rw=randrw, bs=4K-4K/4K-4K/4K-4K, ioengine=libaio, iodepth=64
fio-2.1.5
Starting 1 process
test: Laying out IO file(s) (1 file(s) / 1024MB)
Jobs: 1 (f=1): [m] [100.0% done] [35088KB/11640KB/0KB /s] [8772/2910/0 iops] [eta 00m:00s]
test: (groupid=0, jobs=1): err= 0: pid=23020: Sun Sep 18 16:57:17 2016
  read : io=784996KB, bw=40713KB/s, iops=10178, runt= 19281msec
  write: io=263580KB, bw=13670KB/s, iops=3417, runt= 19281msec
  cpu          : usr=3.67%, sys=10.23%, ctx=356147, majf=0, minf=5
  IO depths    : 1=0.1%, 2=0.1%, 4=0.1%, 8=0.1%, 16=0.1%, 32=0.1%, &gt;=64=100.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, &gt;=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.1%, &gt;=64=0.0%
     issued    : total=r=196249/w=65895/d=0, short=r=0/w=0/d=0
     latency   : target=0, window=0, percentile=100.00%, depth=64

Run status group 0 (all jobs):
   READ: io=784996KB, aggrb=40713KB/s, minb=40713KB/s, maxb=40713KB/s, mint=19281msec, maxt=19281msec
  WRITE: io=263580KB, aggrb=13670KB/s, minb=13670KB/s, maxb=13670KB/s, mint=19281msec, maxt=19281msec</pre>

###### Pricing considerations

EFS pricing is linear, you get billed a fixed amount for GB/month; this is not true with a self-managed cluster where you can surely reach higher performance, but the TCO is increasing every time you add capacity. If you&#8217;re not satisfied with EFS throughput you need a dedicated team to  manage a distributed file system cluster and its operations and maintenance.

###### Conclusions

If you need a stable, realiable file system in AWS to be shared between EC2 instances, go and use EFS and don&#8217;t reinvent the wheel: GlusterFS is outperformed in both read and write IOPS and bandwith with the default options! The workload simulated here is showing poor write performance, so if your use case is a lot of concurrent writes on many files, consider another solution. A good workload could be a WORM (Write Once Read Many) share for permanent stored contents (images, archives?).

The performance are still low in absolute terms or compared to a non-replicated mount such as a stock NFS v4.0 server without the high-availability burden, but if you don&#8217;t care about 100% uptime you should definitely set up NFS with <a href="https://www.drbd.org/en/" target="_blank">DRDB</a> to a secondary node, and switch your mounts when the primary node fails. But still: why manage all of this if AWS can do it for you?
