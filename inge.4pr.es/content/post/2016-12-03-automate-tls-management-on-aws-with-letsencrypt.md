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

Letsencrypt is cool: automated, free TLS certificates for everybody! They are sponsored mainly by internet corps and they started a crowd-funding campaign to avoid the influence of this corps in the future of the project. I recently moved the blog to hugo on AWS and I'm now porting the TLS management scripts I wrote a while ago on AWS: this is a nice exercise to give a proper TLS automation valid for everyone on AWS.

I used [Terraform](https://terraform.io "terraform"), the AWS CLI and the amazing (although still in beta) [lego](https://github.com/xenolf/lego "lego"); the idea is to spin up an EC2 instance every month to query Letsencrypt to renew a  certificate I already provisioned; the instance will have permissions to update the certificate using [AWS Certificate Manager](https://aws.amazon.com/certificate-manager/ "aws acm") and other services like API Gateway custom domain names and Cloudfront are already bound to use ACM certificate of the domain. I tried and make the Terraform code as abstract and reusable as possible: cloudfront distribution id, domain name and issuer mail are all variable configurable via config file or at command line.

#### Setting up the environment
I prepared the TLS certificate and key running first lego once on my local box, using the DNS challenge against AWS Route53, lego will use Letsencrypt ACME secret and put it into a TXT record for the domain to be validated so that Letsencrypt will know I am the owner of the domain. You can read more on Letsencrypt domain ownership validation [here](https://letsencrypt.org/how-it-works/ "Letsencrypt how it works").

After the lego account is created, generally in `$HOME/.lego`, you will find private key and certificate is created for the domain(s); the certificate is a bundle of all the SAN submitted in the request. Then how to automate all of this on AWS?

#### Designing automation
The idea is that every service that needs TLS encryption will be able to fetch the certificate from an AWS service. This was not possible until AWS Certificate Manager was released a while ago: ACM will create or store a certificate to be used in other AWS services via the AWS API. No more magic to spread certificates via S3 or other tricks, you now have an `awscli` command!
So here I chose to use Terraform to bootstrap all the infrastructure to run lego and run Letsencrypt automatically via AWS Autoscaling group.

I first imported the certificate created with lego in ACM in us-east-1 region (N. Virginia): this is due to API Gateway limitations to be able to integrate natively with ACM only in us-east-1. I get the certificate ARN for the certificate just uploaded and I run `terraform plan -var 'cf_distribution_id=EX4MPL3D15TR0' -var 'certificate_arn=...'`.
Once completed I can see a new Autoscaling group with scheduled policies, a launch configuration that starting from the base Amazon AMI with a user-data script at each boot will:
*  download and install lego binary
* sync the TLS account with the S3 bucket created
* split the certificate from the chain, as they need to be uploaded separately
* call the `acm` subcommand of `awscli` to update the certificate

This is amazing, I can now forget about TLS management for all of my websites!
