# Zabbix: a powerful yet simple monitoring software

It may come in mind to any IT system engineer to know what is the status of the network, server by server, instance by instance; it happened to me when I was given the responsibility to manage my company&#8217;s infrastructure and I was wondering which tool could have helped to do the job.

I chose Zabbix to monitor my infrastructure because:

  1. despite it&#8217;s a bit difficult to install (you need a PHP enabled web server, a database and a C compiler), you will benefit a very user-friendly web interface with lots of functionalities
  2. native agents for major OS release are already complied: FreeBSD, Linux, Windows, etc&#8230; Compiling to other OS just requires a &#8220;configure && make && make install&#8221;
  3. it offers many monitoring methods via a unique interface: you can group SNMP, JMX, HTTP monitoring in one shot
  4. it has multi-step HTTP/HTTPS monitoring, simulating different browsers and clients
  5. you can build nice infographics bundling all kind of monitored datas
  6. you can manage users and roles to give access to the web interface at your company&#8217;s employees
  7. you can build custom monitoring scripts to your needs

Well let&#8217;s see some action now: I would like to post a short tutorial on how to build a custom script to monitor resources used by a Glassfish application server. You can use this methodology for other application servers or services.

**Requirements:** you have installed Zabbix server, deployed an agent to a host, set up the necessary networking stuff

On the host to be monitored, you will have a directory where the agent configuration file is located (usually /usr/local/etc or C:\Zabbix)

**Step 1: enable custom parameters **parsing****

Edit the file zabbix\_agent.conf or zabbix\_agentd.conf (depending if you&#8217;re usgin the daemon or not) and uncommment/add the following line:

_Include=/usr/local/etc/zabbix_agentd.conf.d/  _or_ _Include=/usr/local/etc/zabbix_agent.conf.d/__

**Step 2: write the script**

Create a file, name it as you please and insert the script you want to be executed by the agent: I needed a script that would inetract with Glassfish, so I used th following:

> \# Flexible parameter to grab global variables. On the frontend side, use keys like glassfish.status[server.jvm.heapsize-current].
  
> \# Key syntax is glassfish.status[monitoring-key].
  
> UserParameter=glassfish.status[*],/opt/glassfish/bin/asadmin get &#8211;user admin &#8211;passwordfile /opt/glassfish/bin/.pwd -m $1 | cut -d &#8220;=&#8221; -f 2 | tr -d &#8216; &#8216; | bc

The script syntax is always UsrParameter=name.of.script[*] followed by the code to be executed. This one uses the Glassfish utility &#8220;_asadmin_&#8221; and a couple of shell commands to trim the string output and translate it into an integer value. You can see the arguments array can be retrieved using $ and index of argument. In this example you will call the script with one argument only (the monitoring data you want from Glassfish).

**Step 3: start harvesting datas!**

Once finished editing the script, go back to the Zabbix monitoring console and add an Item to the host you are monitoring. You will add the key as shown in the picture below:

&nbsp;

[<img class="aligncenter wp-image-189 size-full" title="zabbix_custom_script" src="http://inge.4pr.es/blog/wp-content/uploads/2013/10/zabbix_custom_script1.png" alt="Custom script insert" width="663" height="438" />][1]

&nbsp;

&nbsp;

Then go back to the dashboard and verify that the script just created is returning datas as expected. In the section _Monitoring->Latest Data_ check if the item is giving the expected values. In this exemple I chose to monitor the current heap size used by the server. One cool thing is: once the script is done you can call it with all the parameters Glassfish has, and then combine datas in an infograph like the following:

&nbsp;

[<img class="aligncenter wp-image-190 size-full" title="zabbix_graphs" src="http://inge.4pr.es/blog/wp-content/uploads/2013/10/zabbix_graphs1.png" alt="Zabbix graph" width="536" height="392" />][2]

&nbsp;

Here I put together two Glassfish parameters (JVM upper bound and current heap used) and a system parameter (free memory).

To get a list of all parameters you can monitor via Glassfish _asadmin_ command see Glassfish documentation <a title="Glassfish documentation" href="https://glassfish.java.net/documentation.html" target="_blank">here</a>.

Cheers,

inge4pres

&nbsp;

 [1]: https://inge.4pr.es/blog/wp-content/uploads/2013/10/zabbix_custom_script1.png
 [2]: https://inge.4pr.es/blog/wp-content/uploads/2013/10/zabbix_graphs1.png

