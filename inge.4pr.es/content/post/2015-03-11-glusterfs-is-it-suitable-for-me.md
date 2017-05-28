---
title: 'GlusterFS: is it suitable for me?'
author: inge4pres
layout: post
date: 2015-03-11T13:00:46+00:00
categories:
  - tech
tags:
  - cloud
  - gluster
  - storage

---
During the last years I&#8217;ve been experimenting with <a title="GlusterFS" href="http://www.gluster.org/" target="_blank">GlusterFS</a> and his functionalities as distributed object store; a lot has changed in the software, overall since <a title="Red Hat acquires Gluster" href="http://www.redhat.com/promo/storage/press-release.html" target="_blank">Red Hat acquired it</a>. I have been using it and find it useful for many projects but not for others: what I love is the community oriented approach with a very responsive team and support for any kind of users (meaning from the 2 nodes web server to a RAID10 Infiniband cluster for high end storage).

My personal story with Gluster starts with a porting of a on-premise architecture in the cloud: moving an existing application to the cloud, instead of redesigning it from scratch, involves a lot of engineering to adapt the current system settings to a scalable infrastructure. Gluster comes handy when talking about scaling: the latest milestone has a very simple and efficient way of reconfiguring the underlying hardware, adding and removing nodes in the storage pool is as simple as inputting a couple of commands from any of the peers in the cluster.

If you&#8217;re unfamiliar with Gluster concepts (storage pool, peers, etc&#8230;) I suggest you RTFM on <a title="Gluster Docs" href="http://www.gluster.org/documentation/" target="_blank">Gluster&#8217;s website</a>; in this post I will detail a few points you won&#8217;t find on documentation and you should definetely know before starting to evaluate Gluster adoption.

<li style="padding-left: 30px;">
  <span style="text-decoration: underline;">Gluster is not a replacement for disaster recovery and backups</span>
</li>

<p style="padding-left: 60px;">
  If you think that data redundancy mechanisms built in Gluster (replication and georeplication) are substitutes for backups you&#8217;re doing it wrong: in Gluster there is no way of recovering data present only in failed drives or unavailable portions of the pool. There is no SPOF free implementation that will avoid you regular backups, unless you can tollerate loss of data.
</p>

<li style="padding-left: 30px;">
  <span style="text-decoration: underline;">Gluster has limited configuration options</span>
</li>

<p style="padding-left: 60px;">
  Gluster has been developed to &#8220;take common hardware and turn it into scalable high performance storage solution&#8221;. Gluster is great when availability and durability are performance indicators because it has been thinked for horizontal scalability, but scaling vertically in system resources will not have the desired outcome. There are phisycal thresholds in a node&#8217;s configuration that make huge hardware resources useless (eg, limit to the number of CPU threads for the transaltor).
</p>

<li style="padding-left: 30px;">
  <span style="text-decoration: underline;">Do consider the application scenario</span>
</li>

<p style="padding-left: 60px;">
  If you are uncertain about Gluster capabilities, <a title="Gluster getting started" href="http://www.gluster.org/documentation/quickstart/" target="_blank">try it out yourself</a> installing the software on at least two virtual machines and test if your application works well with the native FUSE module. As storage layer for I/O intensive applications Gluster is useful when average file size is bigger than the minimum size of the read cache (4MB). Currently the Gluster community is <a title="Gluster small files performance" href="http://www.gluster.org/community/documentation/index.php/Features/Feature_Smallfile_Perf#remove_io-threads_translator" target="_blank">discussing how (relatively) small files</a> should be handled in the next major release of milestone 3 (release 3.7) scheduled for the end of April 2015, but for now if your application scenario has lots of small files written frequently, Gluster may not be the right chioce.
</p>

If you find Gluster is not suitable for your application, consider analizying  a different solution like <a title="DRBD" href="http://drbd.linbit.com/" target="_blank">DRBD</a>: it may not be as cutting edge as Gluster or <a title="Ceph" href="http://ceph.com/" target="_blank">Ceph</a> but may be the right solution for the job.