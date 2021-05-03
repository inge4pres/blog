# DevOps: you're doing it wrong

Recently I received a mail pointing me to a post about DevOps culture and some anti-patterns and misconception on how to build and grow a DevOps culture in a company. Whoever like me works in the Enterprise ("the one with the big E" - Kelsey Hightower) knows that applying DevOps practices often is limited to the adoption of some tools or the creation of a "DevOps team" responsible of managing some continuous delivery pipeline. I would like to share my personal experience and possibly explain what to show to your boss when he thinks they are doing DevOps right.

Luckily for me I had the chance to work 5 years in a small company before DevOps was a thing: we were a small team and were forced to handle development and operations together, scaling our services in feature and size was becoming impossible without the use of some practices that a few years later we discovered was called "DevOps".
When I changed job for an Enterprise company I was happy because I thought I would go and find a better implementation of the same practices, with people trained better than me to handle complex build and deployment systems; I actually found that a deployment system was not there and I was hired to help build one. My delusion was high because there was nothing for me to learn in there, I was stepping back into 3 years of work done automating JVM delivery with Jenkins.
At the very same time the higher management wanted desperately to implement a DevOps organization because studies was talking about how it could increase team productivity and reduce the time to market of features. So this decisions were made:

* a dedicated DevOps team was created to take care of the automation of the build/deployment process (at that time 100% manual)
* a set of tools was chosen by the management and dictated to be used for monitoring, configuration management, etc...
* no change to the rest of organization was made: developers department and operations department kept the very same separate structure
* no change on how the software was developed: waterfall methodology was not removed

Then once the automation of the release was almost complete and some of the hot-sounding technologies on the market (Ansible, Packer, Terraform) were used by some team members the company declared that "through DevOps" they could be the best in breed. The DevOps button was pressed!

There is enough literature around (see bottom of page) to have a common understanding on what a DevOps culture is, how it originated and the practices that make it powerful for modern application delivery; I think most of Enterprise companies struggle impementing a DevOps model because the higher management do not understand what the DevOps transformation requires: a complete rethinking of people, processes and organization. The technological tools that end up implemented in successful organizations implementing DevOps are only a consequence of such change, not the fuel of the change itself.
You cannot switch from SVN to Git, install Jenkins and Gradle and call yourself a DevOps company for doing that. 

There is no silver bullet, there is no one-size-fits-all solution; DevOps is about empowering people to understand continuous improvement and the value of collaboration. There is so much confusion in the higher management of the Enterprise companies that [just changing your job title into "DevOps Engineer" will grant you a salary increase](https://puppet.com/resources/whitepaper/2016-state-of-devops-report) (happened to me).

Take your stand for a truly cohesive organization that put "people over processes over tools" and find the courage to change, experiment and fail. This is what DevOps is about. Or you can keep doing it wrong.

##### Suggested reads 
[Cloud System Adminstration Volume 2](http://the-cloud-book.com/)

[Effective DevOps](https://effectivedevops.net/)

[The Phoenix Project](https://itrevolution.com/book/the-phoenix-project/)

