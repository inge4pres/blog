---
author: "inge4pres"
date: 2017-05-28T17:04:10+02:00
description: "My 2 cents on what is happening with DevOps in the Enterprise"

tags:
  - devops
  - culture
  - enterprise

categories:
  - social
  - work

title: "DevOps: you're doing it wrong"
---
Recently I received a mail pointing me to a post about DevOps culture and some anti-patterns and misconception on how to build and grow a DevOps culture in a company. Whoever like me works in the Enterprise ("the one with the big E" - Kelsey Hightower) knows that applying DevOps practices often is limited to the adoption of some tools or the creation of a "DevOps team" responsible of managing some continuous delivery pipeline. I would like to share my personal experience and possibly explain what to show to your boss when he thinks they are doing DevOps right.

Luckily for me I had the chance to work 5 years in a small company before DevOps was a thing: we were a small team and were forced to handle development and operations together, scaling our services in feature and size was becoming impossible without the use of some practices that a few years later we discovered was called "DevOps".
When I changed job for an Enterprise company I was happy because I thought to go ad find a better imlementation of the same practices, with people trained better than me to handle complex build and deployment systems; I actually found that a deployment system was not there and I was hired to help build one. My delusion was high because there was nothing for me to learn in there, I was stepping back into 3 years of work done automating JVM delivery with Jenkins.
At the very same time the higher management wanted desperately to implement a DevOps organization because studies was talking about how it could increase team productivity and reduce the time to market of features. So this decisions were made:

* a dedicated DevOps team was created to take care of the automation of the build/deployment process (100% manual)
* a set of tools was chosen by the management and dictated to be used for monitoring, configuration management, etc...
* no change to the rest of organization was made: developers department and operations department kept the very same structure
* no change on how the software was developed: waterfall methodology was not removed

Then once the automation of the release is almost complete and some of the hot-sounding technologies on the market were used (Ansible, Packer, Terraform) the company declares that "through DevOps" they can make the best out of any client requirement in no time. The DevOps button was pressed!

There is enough literature (see bottom of page) around on DevOps culture, how it originated and the practices that make it powerful for modern application delivery; I think most of Enterprise companies struggle impementing a DevOps model because the higher management do not understand what the DevOps transformation requires: a complete rethinking of people, processes and organization. The technological tools that end up implemented in DevOps successful organizations are only a consequence of such change, not the fule of the change itself.
You cannot switch from SVN to Git, install Jenkins and Gradle and call yourself a DevOps company for doing that. 

There is no silver bullet, there is no one-size-fits-all solution; DevOps is about empowering people to understand continuous improvement and the value of collaboration. There is so much confusion in the people that should laed the organization of the Enterprise that [just changing your job title into "DevOps Engineer" will grant you a salary increase](https://puppet.com/resources/whitepaper/2016-state-of-devops-report) (happened to me).

##### Suggested reads 
[Cloud System Adminstration Volume 2](http://the-cloud-book.com/)

[Effective DevOps](https://effectivedevops.net/)

[The Phoenix Project](https://itrevolution.com/book/the-phoenix-project/)