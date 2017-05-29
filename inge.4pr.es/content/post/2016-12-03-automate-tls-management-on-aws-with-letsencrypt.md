---
author: "inge4pres"
date: 2016-12-03T15:34:48+01:00
description: "TLS management on AWS"
tags:
  - aws
  - tls
  - letsencrypt

categories:
  - tech

draft: false
title: "Automate TLS management on AWS with LetsEncrypt"
---

Letsencrypt is cool: automated, free TLS certificates for everybody! They are sponsored mainly by internet corps and recently they started a crowd-funding campaign to avoid the influence of this corps in the future of the project. I recently moved the blog to hugo on AWS and I'm now porting the TLS management scripts I wrote a while go on AWS: this is a nice exercise to give a proper TLS automation valid for everyone on AWS.

I used [Terraform](https://terraform.io "terraform"), the AWS CLI, the amazing (although still in beta) [lego](https://github.com/xenolf/lego "lego") and [jq](https://stedolan.github.io/jq/ "jq"); the idea is to spin up an EC2 instance every 15 day to query Letsencrypt to renew a  certificate I already provisioned; the instance will have permissions to update the API Gateway custom domain names and Cloudfront distribution powering the website. I made the Terraform code as abstract and reusable as possible: cloudfront distribution id, domain name and issuer mail are all variable configurable via config file or at command line.

I run `terraform plan -var 'cf_distribution_id=EX4MPL3D15TR0'` and once completed I can see a new Autoscaling group with scheduled policies, a laucn configuration starting from the base Amazon AMI and with user-data script to install lego and sync the TLS account with an isolated S3 bucket. It all works well until the scripts tries to update the Cloudfront distribution, risulting in an error.