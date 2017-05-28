---
title: A smooth migration to the cloud
author: inge4pres
layout: post
date: 2014-10-15T18:15:28+00:00
categories:
  - tech

---
The last weekend my colleagues and I had a nice time moving an existing application from a bare-metal infrastructure to AWS. I would like to share some of the focal points involved in such process, in case you&#8217;d go through it and would like to know:

  * **don&#8217;t expect everything to work as usual**: you are changing the underlying hardware, moving to a virtualized environment. You can test every single part of the application but infrastructural side effects may occur in a second time
  * **relying on the provider**: consider well which functionalities should be delegated to the cloud provider (AWS, in this case, offers a lot) or should be managed internally; for example S3 is not a distributed filesystem, and in some cases an RDS instance won&#8217;t have the same performance as database installed on an EC2 instance
  * **test application compliance, not hardware failure**: instead of focusing on stress tests, you should focus first on functionality tests to ensure every part of te application is behaving as expected; hardware failure are easily handled in the cloud, that&#8217;s the primary purpose of IaaS. What is not handled by the cloud is that the application&#8217;s features will work on it!
  * **use a checklist**: this may seem obvious, but having a clear and well written to-do list with a time table and activities&#8217; details will help you analyze if anything is missing or needs to be done in advance

Aside of this technical considerations, having the support of your coworkers and managers it is what really makes the difference: it keeps you focused in every step and at the same time helps if any problem comes up. That&#8217;s why my boss decided to take several videos with his phone and produced a &#8220;movie&#8221;,  hope you like it 😀