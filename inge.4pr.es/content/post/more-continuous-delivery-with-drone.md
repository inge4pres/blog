---
date: "2018-04-14T11:32:51+02:00"
description: "Pipeline with multi-environment stage and approval"
title: "More Continuous Delivery with Drone"
authors: ["inge4pres"]
categories:
  - tech
tags:
  - devops
  - docker
  - containers
  - platform
draft: true
---

In a [previous post](https://inge.4pr.es/2018/02/25/continuous-delivery-with-drone/) I started using _drone_ to see what its CD capabilities are. I really enjoyed it so I thought it could be good to try a more complex workflow and see how the tool can handle it.

#### Expectations
I'd like to setup a build-test-deploy pipeline where the deployment in the _"production"_ environment happens only after an approval; the approval will be manual for now but can surely be automated to achieve Continuous Deployment.
Let's dive in.

#### Setting up the playground
I have already a docker-compose stack built for the previous post so I'll use that. Only thing I need to dig a little deeper into is what kind of authorization is used for approval. Reading the [gated builds doc](http://readme.drone.io/releases/0.6.0-rc.1/gating/) I am unhappily presented the decision of the creatot not to support any kind of approval in the core product engine, suggesting to use a plugin to do it.
Everyone who worked with me for the past 5 years knows how I hate Jenkins for its plugin architecture implementation, but here plugins are containers so let's give it a try...
I will write an HTTP REST endpoint 