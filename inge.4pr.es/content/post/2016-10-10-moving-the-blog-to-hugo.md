---
title: "Moving the blog to hugo"
date: 2016-10-10T20:07:35+02:00
categories:
  - tech
  - social
author: "inge4pres"
---
In times of experimenting, I am now having a lot of fun with docker, rkt, kubernetes and containers ecosystem in general. But one thing I never forget to play with is content editing and publishing! So here I am, trying to migrate all my blog and website to Hugo :)

So instead of a bare VPS I am moving my blog to AWS S3 + Cloudfront CDN. This will be more scalable and far less expensive. And Hugo generates static HTML so no more patching security issues.

### How I did it

##### Migrate contents
First I found that there is a page on Hugo site that explains how to migrate from any CMS know to human to Hugo. I used the wordpress export plugin to convert the site from Wordpress to markdown and make it work ith Hugo: it did the the job pretty well but still I had to convert some of the posts.

##### HTTPS
Then the hard part: HTTPS! I chose to use [letsencrypt](http://letsencrypt.org "lestencrypt") to automate the TLS certificate handling. The powerful [lego](https://github.com/xenolf/lego "lego") is way more multi-platform than the official certbot from EFF written in Python, and has built-in support for DNS challenge on Route53! Generating a new certificate of renewing is [a couple of commands away](https://gist.github.com/inge4pres/597bb9350ff3e9cc43ecb476a10e636b "gist").  
