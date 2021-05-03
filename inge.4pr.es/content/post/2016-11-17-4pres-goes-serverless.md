---
author: "inge4pres"
date: 2016-11-17T18:28:23+01:00
description: "Moving 4pres URL shortener to a more modern architecture"
tags: []
categories:
  - social
  - tech
title: '4pres goes #serverless'

---
Last month I felt I was a little late for the [#serverless](https://charity.wtf/2016/05/31/wtf-is-operations-serverless/ "#serverless") party going on all over the internet and I started taking a look at what the pros and cons would be to actually not manage any server myself.
Shutting down my VPS hosting my apps I will lose my mail server, my MySQL instances and my Docker registry but: who cares? There are cloud services I can use with hundreds of times more availability and for a fraction of the cost.

# Why #serverless?
We are moving more and more rapidly to a developer friendly world: all cloud providers tend to relief companies from the burden of managing complex cluster architectures and it's not a coincidence that things like [SRE at Google](https://landing.google.com/sre/interview/ben-treynor.html "SRE at Google") exist. Systadmins are no longer required where they can be substituted with far more performance by declarative syntax clusters (Kubernetes, Docker Data Center) and reliabile, consistent configuration deployment (Ansible, Puppet) in managing high volumes of phisycal or virtual machines.

This shift to a developer-centric world forces who embraces it to "trust" the IaaS, PaaS and FaaS providers but in the same time let them focus on core and valuable development processes.

That's why.

Take the simplest app in the world: [4pres](https://4pr.es), a URL shortener. It needs a presentation, computation and data layer as most of apps. Traditionally and depending on your budget and needs, you would spin up one or more VPS and deploy software on them (containers or not, doesn't matter here). The setup of nginx as reverse proxy to your application already requires skills that most of developers don't have, but in first hand why should anybody have them when there is a service like AWS API Gateway that lets you deploy one in seconds?
Having the possiility to do so, you _may_ want to forget about everything not strictly related to your app, so focusing only on building and maintaining functionalities of your app.

# How you can migrate your app today     
In terms of cloud provider I have a lot of experience with AWS therefore my first thought is for them when trying to do something like this: they probably already have enough mature services supporting what I need. Don't take this as a sponsorship: you can do the same with any other provider.

#### Traditional Architecture
* NGINX reverse proxy
* Golang application 
* MySQL database

NGINX terminates SSL and proxy back to app every request. The [Golang app](https://github.com/inge4pres/4pr.es "4pres/master") finds out what to do from a combination of HTTP Method and URL Path:

* `GET /` renders the landing page template
* `GET /{url}` queries the DB and redirect to long url or render 404 template
* `POST /` create a short link from `form-data url=http://alongurul/` and display the result template

#### Serverless Architecture (AWS)
* API Gateway expose a request/response mapping endpoint with integration to other services
* S3 to store and serve the static content 
* Lambda executes Golang functions thanks to the amazing [eawsy/aws-go-lambda](https://github.com/eawsy/aws-lambda-go "eawsy/aws-go-lambda") framework
* DynamoDB stores data in schema-less fashion (JSON)

The API Gateway definition acts as proxy and:

* `GET /` serves static page `index.html` from S3 bucket with 'Website Hosting' option enabled
* `GET /{url}` runs a Lamdba function [4pres_get](https://github.com/inge4pres/4pr.es/blob/serverless/cmd/4pres_get/get.go "4pres_get") that fetches the URL Path parameter from DynamoDB and redirect the client or renders a 404 template 
* `GET /s?{urlencodedURL}` runs a Lambda function [4pres_post](https://github.com/inge4pres/4pr.es/blob/serverless/cmd/4pres_post/post.go "4pres_post") that creates a short URL and tries to store it in DynamoDB, returning the result template or the 500 template.

Not that big change in the overall design, but the code for the Golang app only shares a function to shorten the URL between the 2 implementations: that is understandable because we no longer manage HTTP requests attributes and delegate that to API Gateway, we don't display static content anymore and leave that to S3. At the core we only have 2 things to worry about: 

* store a URL (4pres_post)
* get  URL to redirect the client

and then we can focus on extending new features:

* URL expiration
* user registration
* whatever!

# Conclusion

Thanks to this development model our craftsmen effort can be 100% dedicated to building features and forget about the rest: this can be a great advantage at large scale. Recently I have come across another interesting framework for Function as a Service called [iron.io](https://www.iron.io/ "iron.io") and I am willing to try it as soon as possible so stay tuned!
