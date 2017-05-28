---
title: Set up a private master/slave DNS using BIND
author: inge4pres
layout: post
date: 2014-12-19T17:20:13+00:00
categories:
  - tech
tags:
  - BIND
  - CentOS
  - DNS
  - Linux

---
One of the very basic need of any startup is setting up a LAN in the workspace and configuring the Internet most used service: DNS. Relying on a public DNS may give you full functionality towards WAN connectivity, but when you need to address some hosts inside your LAN it can be handy to use names instead of IPs (especially with IPv6).

Here&#8217;s a straight forward guide to get you started with your private DNS in a few minutes.

OS filesystem&#8217;s path and package management utility may vary with the flavour of your distro, here I use CentOS.

* * *

**Requirements**

  * a router with DHCP ad WAN connectivity already setup
  * 2 CentOS 6.x servers (physical or virtual with due availability concerns), enable to network with each others
  * a desktop PC
  
    * * *

**Some general info**

I use this data as example, change them to your needs

<li style="text-align: left;">
  <code>startup.me</code>             the domain name you want to use in your LAN
</li>
  * `10.20.30.0/24`        the subnet of your LAN
  * `10.20.30.40`            the static IP address assigned to CentOS server 1, with hostname  `centos1.startup.me`
  * `10.20.30.50`            the static IP address assigned to CentOS server 2, with hostname  `centos2.startup.me`

* * *

**Configuring the primary DNS server **

On server `centos1`

`[root@centos1 ~]# yum update && yum -y install bind bind-libs bind-utils`

The BIND daemon is now installed; the base dir for the service is `/var/named` and the configuration file is `/etc/named.conf` ; modify the configuration file with your favourite editor

`[root@centos1 ~]# vim /etc/named.conf`

In the `options` section adjust the settings to your LAN configurations, changing the example values

    options {
           listen-on port 53 { 10.20.30.40 };             # inet address of centos1
           listen-on-v6 port 53 { ::1; };                 # comment this out to use IPv4 only
           directory "/var/named";
           recursion yes;
           allow-recursion { 10.20.30.0/24; };            # recursion only in LAN, change this with your subnet
           allow-transfer { localhost; 10.20.30.50; };    # enable zone transfers only to secondary DNS sever
           forwarders {
                      208.67.222.222;
                      208.67.220.220;
           };                                             # OpenDNS used here, Google 8.8.8.8, 8.8.4.4 can be used
           dump-file "/var/named/data/cache_dump.db";
           statistics-file "/var/named/data/named_stats.txt";
           memstatistics-file "/var/named/data/named_mem_stats.txt";
           allow-query { 10.20.30.0/24; };                # accept queries only from LAN, change this with your subnet
    };

Now comment the `include` lines following the `options` section and comment out all the rest. The end of the file should be (modify the domain)

    # zone "." IN {
    #       type hint;
    #       file "named.ca";
    #};
    
    #include "/etc/named.rfc1912.zones";
    #include "/etc/named.root.key";
    
    zone "startup.me" IN {
          type master;
          file "/var/named/startup.me.zone";
          allow-update { none; };
    };

The latter lines create the definition of the domain and specify that DNS records should be looked up in the `/var/named/startup.me.zone `file which we are about to write. If you&#8217;re unfamiliar with DNS records and basic concepts, have a look at <a title="DNS IANA's guide" href="http://www.iana.org/assignments/dns-parameters/dns-parameters.xhtml" target="_blank">this</a>. Now create and edit a new file

`[root@centos1 ~]# vim /var/named/startup.me.zone`

Insert this

    $TTL 8H
    @      IN      SOA      centos1.startup.me. root.startup.me. (
                   1        ; serial
                   1D       ; refresh
                   1H       ; retry
                   1W       ; expire
                   1H )     ; minimum TTL
    ; Name servers
            IN      NS      centos1.startup.me.
            IN      NS      centos2.startup.me.
    ; Resolvers
    @       IN      A       10.50.30.10  ; default IP for your domain root startup.me
    centos1 IN      A       10.20.30.40  ; the primary DNS server centos1
    centos2 IN      A       10.20.30.50  ; the secondary DNS server centos2
    www     IN      A       10.20.30.20  ; a web server in your LAN
    ftp     IN      A       10.20.30.30  ; an FTP server in your LAN
    git     IN      A       10.20.30.5   ; a GIT server in your LAN

Be careful when copying the above snippet into the config file: indentation must be respected!

You have configured the primary DNS with some DNS A records; the provided settings are not for a heavy load DNS server, nor they should be used in networks wih frequent IP address change: the time-to-live settings are high, therefore any IP mapping change should be followed by a clients&#8217; resolver cache flush, which may be inconvenient for most users.

Now it&#8217;s time to enable the service; verify files permission

`[root@centos1 ~]# chown named /var/named/startup.me.zone`

and start the service

`[root@centos1 ~]# service named start`

If you want to have the service enabled at boot

`[root@centos1 ~]# chkconfig named on && chkconfig save`

To test your primary DNS server you should first be sure to use it; check the resolver settings

`[root@centos1 ~]# vim /etc/resolv.conf`

it should contain only one line

`nameserver 10.20.30.40`

Verify the DNS is working with

`[root@centos1 ~]# dig www.startup.me`

if the output contains &#8220;AUTHORITY SECTION&#8221; you&#8217;re done.

* * *

**Configuring the secondary DNS server **

On server `centos2` install BIND as you did for `centos1`.

`[root@centos2 ~]# yum update && yum -y install bind bind-libs bind-utils`

Secure-CoPy the configuration files from `centos1` to `centos2`

`[root@centos2 ~]# scp centos1.startup.me:/etc/named.conf /etc/named.conf`

`[root@centos2 ~]# scp centos1.startup.me:/var/named/startup.me.zone /var/named/startup.me.zone`

Edit the `/etc/named.conf` adjusting the IPv4 address of `centos2` in the `options` section, change `10.20.30.40` with `10.20.30.50`, and at the end of the file modify the zone settings to have a secondary (slave) DNS; as usual, change the IP with your actual primary DNS IP.

    zone "startup.me" IN {
          type slave;
          masters { 10.20.30.40; }; 
          file "/var/named/startup.me.zone";
    };
    

Start the BIND daemon on `centos2` and enjoy!

`[root@centos2 ~]# chown named /var/named/startup.me.zone`

`[root@centos2 ~]# service named start`

There it is! You are ready to test it from your desktop.

In your router settings change the primary and secondary DNS servers for the DHCP server, renew all adresses and try browsing the domain with any software. So name, much fast, wow!

_Note: having a script doing an rsync to your secondary DNS will ease the pain when the primary DNS server goes down._

Cheers 😀