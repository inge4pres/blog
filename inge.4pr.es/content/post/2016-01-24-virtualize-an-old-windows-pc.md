---
title: Virtualize an old Windows PC
author: "inge4pres"
layout: post
date: 2016-01-24T15:42:29+00:00

categories:
  - culture
  - tech
tags:
  - geek
  - guide
  - oldies
  - virtualization
  - windows

---
A friend asked me if I was able to get back working a Windows 98 PC he had in his house; I have never done it so I said &#8220;sure I can!&#8221; just to have the opportunity to learn something new, and of course do a friend a favour.
  
My idea was to copy the whole PC and get it running on a virtual machine thus doing what I later discovered is called a &#8220;P2V&#8221; (Physical to Virtual); the result of which would have been a portable VM which I could then install in my friend laptop to have his old PC back. The idea was indeed good and I am writing this blog just after finishing the job in the hope I will save someone the headache I&#8217;ve had in the last 4 days to get this done.
  
My tools for this job are :

  * an ATA/SATA/IDE to USB adapter (<a href="http://www.amazon.it/gp/product/B001EOO660?psc=1&redirect=true&ref_=oh_aui_detailpage_o00_s00" target="_blank">this</a> exactly) to mount the old drive
  * <a href="http://www.winimage.com/download.htm" target="_blank">winImage</a>: a fantastic freeware by Gilles Vollant (runs on Windows only)
  * Windows 98 ISO (get one <a href="https://winworldpc.com" target="_blank">here</a>)
  * VMWare Workstation: to run winImage on a pre-installed Windows VM, configure and test the new Windows 98 VM

Why VMWare? It is by my knowledge the most reliable and portable hypervisor, with a vast documentation and huge community support; I also had it already installed in my Fedora box running a virtual Windows 10. Don&#8217;t worry if you don&#8217;t have VMWare: VirtualBox or Virtual-PC will do the same&#8230; So let&#8217;s get it started!

Verify the drive is working connecting it to the USB adapter: connect the USB to your PC first, connect the power adapter to te disk then finally plug the power cable; if a green light is on and you hear the disk spinning, move ahead. If the disk has no vitals, consider calling a data recovery expert. If you&#8217;re on Windows the disk should be visibile in your drives: press the Windows logo and E key to open explorer and the drive should be visible and browseable. Otherwise check the Disk Utility to verify it has been recognized and why is not mounted. On Linux, use `fdisk -l` to know the device and `mount -t filesystemtype /dev/deviceid /your/mount/point` to mount on a folder specifying the filesystem type. To know the file system type use `file -sL /dev/deviceid` using the device identifier.

Now that the disk is operational I recommend to take a backup with a raw image: this will help if any trouble happens with the disk during the P2V process. On windows you can use ShadowCopy or Ghost or Acronis Image, it all depends on your budget. On Linux I will use

`dd if=/dev/deviceid of=/my/backup/destination/disk.img`

Now I feel more confident because any action can be reverted with the original disk image backup: the next step is to virtualize the disk and all of its content with WinImage. From the menu select &#8220;Convert physical to virtual&#8221;, choose the input drive, a destination folder and the type of virtual disk (.vmdk if you want to use VMWare); WinImage will ask if you want to make a backup of your drive and if you are paranoic like me, answer yes and save another backup. Depending on the disk size the convertion process may take about 30 minutes long, so relax and wait; once WinImage is done you will notice it because the disk will lower down the noise.

Open VMWare (or any other hypervisor you like) and create a new virtual machine: for Windows 98 I choose 1 vCPU with 512MB RAM; attach the virtual drive created with WinImage to the guest and start it up. If you&#8217;re lucky enough you should end up with a Windows boot screen and the operating system loading.

Now the fun part! As we took a raw image of the old disk the drivers the hypervisor will use won&#8217;t be recognized! So in my case I had to click through a lot of driver reinstall, but they were all already present. A couple of reboots and the Windows 98 system is back up and running. Don&#8217;t be discouraged by what seems an infinite loop of driver reinstall! Continue installing the missing drivers and you willl succeed! This is definitely the most important thing to do: never give up!

&nbsp;
