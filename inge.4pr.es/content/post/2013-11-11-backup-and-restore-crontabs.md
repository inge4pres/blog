---
title: Backup and restore crontabs
author: inge4pres
layout: post
date: 2013-11-11T10:27:38+00:00
categories:
  - tech
tags:
  - bash
  - cron
  - scripting
  - shell

---
Here&#8217;s a thing I came up with: you&#8217;re administering a Linux system with 100 users circa and you&#8217;re moving to a new server, you can save crontabs  per user with this:
  
``mkdir crontabz && cd crontabz; for user in `cat /etc/cron.allow`; do crontab -l -u $user > cron_$user; done``
  
you will end up with a list of files _cron_xxx_, each one has the users&#8217; cron. Hopefully your cron version will use the _/etc/cron.allow_ file to control  user based access to crontab.

Now, once users in the new system are all in place you can copy the directory _crontabz_ and then restore their cron with:
  
``cd crontabz;  for file in `ls -m1`; do echo `basename $file`|sed -s 's/cron_//' >> temp ; done ; \ ``
  
``for user in `cat temp`; do  cat cron_$user | crontab -u $user -; done;``
  
Optionally you can include an &#8220;rm -f temp&#8221; at the end to delete the file used to store user names.

I put the script on github <a title="crontab_manager" href="https://github.com/inge4pres/bash_tools/blob/master/crontab_manager" target="_blank">here</a>,

Cheers