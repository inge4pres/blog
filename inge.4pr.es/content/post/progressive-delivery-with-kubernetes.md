---
title: "Progressive Delivery with Kubernetes"
description: "Shipping fast and keeping most of your users happy"
date: "2019-05-05T15:00:51+02:00"
author: "inge4pres"
emoji: true
comment: true

categories:
  - tech
  - work

tags:
  - containers
  - kubernetes
  - platform
  - continuous-delivery

draft: false
---

I'm more and more fond of finding the perfect solution to manage application delivery: dev teams want to be fast but their ops counterpart is not happy to loose control over the growing number of deployments that could cause an outage. We as an industry need to find the right balance to have features delivered in time and keep the service up and running for our users! And that's where progressive delivery can help!

What is progressive delivery? It's the evolution of continuous integration and continuous delivery practices, taken to the extreme - if this is the first time you hear about it, read [this excellent post by James Governon](https://redmonk.com/jgovernor/2018/08/06/towards-progressive-delivery/). But as of today what are the tools that embed this practice in deployment pipelines? None that I could find ☹️... hence I started this post to share some of the techniques that you can use to achieve progresive delivery today on Kubernetes! I'd be really glad to have any comments and discuss on the matter.

## The goal

The progressive delivery manifesto (if there will ever be one) should explain the rationale why delivering feature in parts is better than all at once: feedback. 
In this Agile world feedback is everything, and the only feedback that matters is your users'; as Jez Humble puts it

> “Users don’t know what they want. Users know what they don’t want once you’ve built it for them.”

You are not going to build anything useful if you don't collect your users opinion _while_ building the product, that is why having a system that is able to change quickly and that can collect this feedback is so vital to success.

## To mesh or not to mesh

The first approach to progressive delivery is via infrastructure components.
I cheated a bit in the post intro: there actually is a tool combination that is able to implement a feature close to progressive delivery right now: it's [Istio](https://istio.io/) plus [Spinnaker](https://www.spinnaker.io/); the network mesh in this scenario is a router for connections originating by clients between multiple versions of the same backend, whose deployed versions are managed by Spinnaker releases. The mesh could be another product (Envoy, Linkerd, Consul Connect...) but it is the necessary component that contains the logic to serve the user a specific version of the application, based on goegraphic location rules, latency or even application rules (layer 7).

If you want to avoid the burden of installing and maintaining a mesh network you need to manage custom tooling to have the traffic routed for a subset of users to a specific version, [Skipper](https://github.com/bookingcom/shipper) is a good example but comes with the restriction of not being able to manage percentage of traffic, so the percentage of user served is based only on the number of pods configured from service to service (so not ideal for small sized deployments).

The other way I see right now is creating a [Kubernetes operator and a CustomResourceDefinition](/2018/05/05/cloud-native-software-delivery/) that can interact with the Ingress resource: this is hypothetical and I am not aware of any tool that is doing this but it could be posible to configure the ingresses to serve part of the requests by a specific Service (e.g. `v1.2.3` backed by a Deployment with a proper `selector`). As far as I know the current ingress controler based on nginx does not have such feature, but I just discovered writing this line that [Traefik does support this](https://docs.traefik.io/user-guide/kubernetes/#traffic-splitting)! It would be great to understand if Traefik can manage multiple rules at once and if it can be managed via API so that the traffic is gradually moved from service to service.

## Feature flags

Of course if you move to the application things get easier in terms of programmability, problem is they tend to be more difficult to manage at scale. If you use one of the [multitude of available](http://featureflags.io/feature-flags/) feature-flag products (also referred as [feature toggles](https://www.martinfowler.com/articles/feature-toggles.html)) you are soon going to be able to experiment with progressive delivery capabilities; your application will most likely contain the logic required to show a specific user a feature or another. While this is intriguing, if you have more than 2 product teams this can easily become a nightmare because if each team implements its own solution of feature toggle the company as a whole can really struggle to get what type of experience is serving to its users. Change management, for as light as it can be, should still be accounting for features enabled and disabled that may cause a service disuption, even if for a small percentage of users, and when the logic for serving different versions of your system is scattered around multiple applications, this goes quickly out of control.

One approach I've seen succeed in using software-defined toggle is adopt a centralized, company-wide solution around an existing product: this simplifies greatly the management around features that are delivered passing through multiple services while being able to keep track of changes consitently.

## Delivering it

Once you've established an infrastructure for serving users based on some policies you should also have in place automation to be able to push your features out. In case you went for the infrastructure/network path you'll need a deployment tool that can sit between your CI artifacts and the platform running your services; on the contrary for a software-driven solution you will just need an application build and deployed regularly.
For the former I am really struggling to find a product that suites my need, I've poked around with Spinnaker, ArgoCD and Tekton Pipelines but none of them seems to have the adequate primitives to address my progressive delivery needs.
I'd be happy to hear from the community how this is being addressed: I'd like to have a descriptive way of defining multiple versions of artifacts and configurations paired together (maybe via commit hash?) and have all of them deployed at the same time; I'd also like to update the configurations of a given version while it's running.

Seems fair right? I might need to tweak my service here and there, but I'd like to tune it only for a specific set of features that I know are in version ABC. Now I could not find on the internet a single product that works with Kubernetes able to satisfy this requirement, so please if you happen to have something in mind leave a comment!

## Conclusion

As always Kubernetes is a great enabler for delivery techniques based on software, it's an extensible platform and the multiple uses that can lead to achieve progressive delivery just confirm it. Personally I see a lot of space for progressive delivery in the upcoming future, especially for IoT. Let's see what's next!