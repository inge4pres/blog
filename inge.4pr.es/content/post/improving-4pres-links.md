---
title: "Improving 4pr.es Links"
author: "inge4pres"
type: ""
draft: false
date: 2023-04-15T23:03:05+02:00
subtitle: Revamping a toy project and making it a product

page:
    theme: "wide"

tags:
  - product
  - software-engineering
  - serverless

toc:
    enable: true
math:
    enable: false


---

In the past few days I have been doing maintenance on a pet project of mine: a URL shortener.

I initially built "[4pres](https://go.4pr.es)" in 2013 as a way of experimenting with the Go programming language.
It is a URL shortener, it does one simple thing: it generates a random short URL (prefixed with the domain `4pr.es`)
for a longer one.

Despite its "dumbness", it served the greater purpose of exposing me to some software engineering problems I had never
faced before. It was a practical way of starting my SWE journey, finding my own answers to little complexities that
until then I had marked as "someone else's job".

The website was initially running on DigitalOcean, with NGINX fronting the Go application, 
storing the redirects into a MySQL database.
It was not meant for world-scale, but it was in production and available to anyone.
When I say "in production" I mean it: it had monitoring and alerts (via SMS!) to help me recover it in case the 
service was degraded.

[In 2016](/2016-11-17-4pres-goes-serverless), I used the service as guinea pig to start adopting serverless technologies,
specifically porting it to AWS. I had been using AWS for 4 years and it was evident that they were on to something.
I substituted MySQL with DynamoDB, virtual machines with Lambda, and NGINX with API Gateway.

This second challenge was mostly technical as I had to go through the unsupported Go runtime for Lambda (only NodeJS and 
Python were available at that time), changing the storage layer to be document-oriented rather than RDBMS, and discover that
to have support for TLS-encrypted connections I needed a CloudFront distribution on top of API Gateway.

Jumping forward to 2023: many things have improved in AWS since 2016! The serverless hype did bring goodness all around!
AWS Lambda now supports Go natively, DynamoDB has better query controls, and API Gateway has finally TLS endpoints!

I could not resist the temptation of doing some maintenance to my beloved shortener, but this time I felt like I wanted a
different challenge.
_I wanted to make something that others would love_, a _product_ that might be helpful to someone else.
While still needing some technical tweaks, my intent was to focus mostly on the user experience.

So I took the initiative and reworked the website UI, adapted the Go code to the new official Lambda runtime, improved
the storage layer and the frontend! It was fun. A good start, but not what I was after.
I sat down for a while, and pondered what was missing, then I wrote it down as a [roadmap](https://go.4pr.es/features/).

As first improvement, I planned mostly to achieve two results:

1. good user experience: it should be trivial to get started using the service
2. sustainable performance: this could be a blog post on its own, but long story short it means reaching latency budgets
   with the least possible resource usage

In all humbleness, I think I have reached a pretty good result on both so far, and I am ready to share the project
with the world once again.
Be aware: I am not done! As it is all software: never finished, always evolving.

In the next section of the post, I will highlight some of the things I found interesting along the way.

### Pricing

Finding the right pricing model for something that is highly abundant in the market is no easy task.
While I would love to let anyone access 4pr.es for _free_, this is not sustainable at scale.

I figured out the best choice would be to allow for a free tier, and lower the price to the minimum to buy uniquely
personalized links.

The free tier allows trying out the service, without any charges: the shortened link expires after 6 months, it supports
up to 50_000 redirects, but only HTTPS links can be shortened.

Paid tiers, on the contrary, allow shortening HTTP _and_ HTTPS links; the paid short links never expire.
There are 3 tiers, supporting different number of maximum redirects: 5, 50 and 500 million redirects, respectively.

I am planning to add a way to configure the maximum number of redirects and adjust the price accordingly in the future.
Let me know in the comments if you would like to have this feature.

I have implemented many AWS wizardries across the stack, in order to lower the cost of the product, and reduce prices 
for users. More details about it in the next section.

### Sustainable performance

This concept is relatively new: its inception comes from the environmental attention placed in the latest years on 
software. Adopting it means to have the product and its features delivered with as little computational resources as 
possible, while still respecting certain latency budgets.

The product was already entirely running on serverless technologies, which by design "scale to 0" when not used at all.
I had to introduce new APIs to the existing code to serve the paid tiers, and write completely new software to count for
redirects.

With the good spirit of measuring the performance of every test and calculating its costs, I applied several tweaks to
my AWS configurations.
Here below 2 tricks among the many that you can use as well in your AWS deployments.


#### 1. Add a CloudFront distribution on top of API gateway

AWS API Gateway supports TLS connections for a while, so adding CloudFront on top is not needed anymore
_just to have encryption_.

Caching your API calls (where possible) can save a lot of resources!
In the first iteration of serverless `4pr.es`, every redirect request was going through the API Gateway and executing
a Lambda function calling DynamoDB, just to return an HTTP response with the right `Location: ...` header.

These redirects can be cached very efficiently by CloudFront, reducing user latency _and_ incurring into fewer costs, 
as we don't need to run Lambda or query DynamoDB.

#### 2. Use the `provided.al2` Lambda runtime and Rust for batch jobs

In order to keep track of visits to redirect links, I initially had implemented an "epilogue" function that would 
execute after the redirect was returned to the user (a `defer` in the main goroutine).
When introducing cached redirects, this would no longer be possible.

I explored the use of CloudFront Functions and Lambda@Edge, both practically allowing to run some code on the
CloudFront point of presence: my idea was to run an AWS client call to DynamoDB to update the stats for the link.

I could not make use of neither because of [their limitations](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/edge-functions-restrictions.html):

* [no network access](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/functions-javascript-runtime-features.html#writing-functions-javascript-features-restricted-features)
* only Python or NodeJS are supported (un-sustainable performance!)

I then figure out I could use CloudFront access logs and Lambda to keep track of redirects.
This method is cheap and effective, as the logs are stored compressed in S3, and they can be deleted once processed.
A Lambda function can do it, processing each file in its own thread and using as many threads as available in the runtime.

This Lambda function runs as a batch job multiple times a day, and it's the most cost-effective way I could find to store
statistics for the visited links.

Along the way, I found out that:

1. CloudFront logs are very verbose, but their tabular, one request-per-line format is dead easy to parse - and it can be easily parallelized 
2. Lambda on ARM is **much faster than x86**, we're talking of 15-20% latency reduction for the same function just by recompiling it - always
   prefer the Graviton runtime if you can compile your code to ARM64
3. Lambda supports Rust natively now! There is even a `cargo lambda` plugin that facilitates onboarding, [check it out here](https://cargo-lambda.info/)
4. The [AWS Rust SDK](https://github.com/awslabs/aws-sdk-rust) is experimental, but it looks like the API surface is going in the right direction

### Wrapping up

[4pr.es](https://go.4pr.es) is waiting for you!

I'll keep refining the product when I have time to do so, if you have feature requests or proposals, feel free to drop
a comment below or send feedback on GitHub via the website.

**Happy shortening!**
