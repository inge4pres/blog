---
title: "gRPC Traffic Mirroring Using NGINX"
subtitle: "Unveiling the power of mirror, grpc and proxy modules"
date: 2021-05-30T17:23:35+02:00
lastmod: 2021-05-30T17:23:35+02:00
draft: false
author: "inge4pres"
description: "An NGINX configuration trick to mirror traffic to gRPC backends"

page:
    theme: "wide"

upd: ""
authorComment: ""

tags:
    - SRE
    - infrastructure
categories:
    - tech

hiddenFromHomePage: false
hiddenFromSearch: false

resources:
- name: "featured-image"
  src: "nginx-grpc.png"
- name: featured-image-preview
  src: featured-image-preview.jpg

featuredImage: ""
featuredImagePreview: ""
images: [""]

toc:
  enable: false
math:
  enable: false

license: "MIT"
---

Recently at work with the [Optimyze](https://optimyze.cloud) team we faced the necessity of copying traffic from our current customer-facing environment to a new environment.
We have assumptions and ideas about architectural changes that cannot be validated only with synthetic tests and require cloning traffic to a separate, internal testing environment.

There is no better test than the one performed with real-world data: when you hear speaking about _testing in production_, a deployment of a new feature to "see what happens" is not what I have in mind.
I think of techniques that allow testing _with_ production such as traffic mirroring (or _shadowing_), canary releases, A/B testing and feature flags.

These techniques will allow gathering data directly from your customers, in the form of raw input data (web/API requests) or feedback recorded from customers' interaction with a new product version.

In the case of traffic mirroring we don't introduce any changes to a production system, we simply copy the traffic (entirely or a portion of it) into a different system, for internal analysis.
When applying this technique, security considerations apply: you do not want to leak customers data into a non-production system,
so ensure you apply the same security policies to the shadow infrastructure receiving a copy of production traffic.

### TL;DR

**If you're in a rush and don't care about the details, the solution is [here](#the-solution).**

### The problem

We ♥ NGINX!
We knew it has the ability to serve as [reverse proxy for gRPC services](https://www.nginx.com/blog/nginx-1-13-10-grpc/),
and it has [built-in support for mirroring](https://nginx.org/en/docs/http/ngx_http_mirror_module.html).

What we didn't know is that mirroring gRPC traffic has a caveat: URI rewrite is not possible within the `grpc_pass` directive.
This means that traffic copied over to a new location (as described in the mirroring module docs) will have an URI set with the mirror location, making the service unable to map it with any RPC when received. 

### Debugging NGINX capabilities  

NGINX documentation is not very rich of examples or references to complex configurations, it is often left to the user 
to build the desired outcome using the various building blocks offered by the modules.

In the case of traffic mirroring, the documentation states:

```
Sets the URI to which an original request will be mirrored. 
Several mirrors can be specified on the same configuration level.
```

In first instance we thought we could simply copy directly the traffic together with `grpc_pass` in a naive configuration like the following (**not working as intended!**):

```nginx configuration
server {
	listen 443 ssl http2;
	server_name grpc-backend.company.net;
    ssl_certificate /etc/letsencrypt/live/company.net/fullchain.pem;
	ssl_certificate_key /etc/letsencrypt/live/company.net/privkey.pem;
	root /var/www/default;
	
	location / {
        mirror /grpc-mirror;
        grpc_pass grpc://production-upstream:12345;
	}

    location = /grpc-mirror {
        internal;
        grpc_pass grpcs://mirror-upstream:443;
	}
}
```

Note 2 things here:

* we are using NGINX to do TLS termination and we have a production service listening at `production-upstream:12345` for unencrypted traffic
* we are re-encrypting traffic when copying over to the mirror because we want to emulate the full end-to-end traffic that a client would generate, and that includes TLS termination

This may not apply to you but I prefer to sponsor a secure-by-default configuration, you can get free and automated TLS certificates with [LetsEncrypt](https://letsencrypt.org/).

We noticed the shadow environment was not processing data, so it was time to enable some debugging and see what was happening.

We add some logging to help us troubleshoot what is going on.

```nginx configuration
	location / {
        mirror /grpc-mirror;
        grpc_pass grpc://production-upstream:12345;
	}

    location = /grpc-mirror {
        internal;
        grpc_pass grpcs://mirror-upstream:443;
	}
	error_log debug;
```

We can now find these messages in the logs:

```
2021/05/27 20:15:51 [debug] 15765#15765: *1 http upstream request: "/grpc-mirror?"
2021/05/27 20:15:51 [debug] 15765#15765: *1 grpc header: "grpc-message: malformed method name: "/grpc-mirror""
```

As documented, the `mirror-upstream` server receives a request with URI `/mirror`, making it useless on the receiving end and creating a path mismatch.

At this point we thought native gRPC mirroring was not possible with NGINX and we tried to copy the traffic via the kernel network stack.

So we introduced a new server on a different port, willing to direct there the traffic _outside_ of NGINX.

```nginx configuration
server {
	listen 9443 ssl http2;

	ssl_certificate /etc/letsencrypt/live/company.net/fullchain.pem;
	ssl_certificate_key /etc/letsencrypt/live/company.net/privkey.pem;
	root /var/www/default;
	server_name grpc-backend.company.net;

	location / {
		grpc_pass grpcs://mirror-upstream:443;
	}
}

server {
	listen 443 ssl http2;

	ssl_certificate /etc/letsencrypt/live/company.net/fullchain.pem;
	ssl_certificate_key /etc/letsencrypt/live/company.net/privkey.pem;
	root /var/www/default;
    server_name grpc-backend.company.net;

	location / {
        grpc_pass grpc://production-upstream:12345;
	}
}
```

When we thought about cloning raw TCP traffic we also had to include TLS termination for the cloned traffic.

With this we are thinking we'll keep the same desired testing properties because the traffic will be cloned without NGINX knowing it,
so we tried with iptables rules using the `TEE` target only to realize minutes later that cloning TCP packets simply does not work because the protocol
is stateful and there is no way for the kernel to handle duplicated responses from a second client that has never been tracked.

Using a separate server is key to the correct solution presented below (thanks [@mejofi](https://twitter.com/mejofi) for the brilliant idea).

### The solution

The trick is to combine the gRPC module with the proxy module to achieve the URI rewrite that is needed.

Here is a working gRPC traffic mirroring configuration:

```nginx configuration
server {
	listen 127.0.0.1:9443 http2;

	root /var/www/default;
	server_name grpc-backend.company.net;

	location / {
		grpc_pass grpcs://mirror-upstream:443;
	}
}

server {
	listen 443 ssl http2;

	ssl_certificate /etc/letsencrypt/live/company.net/fullchain.pem;
	ssl_certificate_key /etc/letsencrypt/live/company.net/privkey.pem;
	root /var/www/default;
    server_name grpc-backend.company.net;

	location / {
        mirror /grpc-mirror;
        grpc_pass grpc://production-upstream:12345;
	}
	
    location = /grpc-mirror {
        internal;
        proxy_pass http://127.0.0.1:9443$request_uri;
	}
}
```

**In the snippets, the definition of upstream has been omitted.**

What we have here:

* the second server, now listens on the loopback interface to avoid traffic leaving the NIC; no longer we need to apply TLS termination again
* the un-encrypted traffic is mirrored to an `internal` location
* a `proxy_pass` directive that will set the original URI in the next upstream, proxying the traffic the loopback server
* the `grpc_pass` in the loopback server will re-encrypt and pass the request to the shadow environment, but with the right gRPC path URI

**Pay attention**: the previous configuration will copy _all traffic_ from a production gRPC service into a shadow environment.

It is left as an exercise to the reader to find the proper way of shadowing only a portion of the traffic (hint: check the 
[split_clients_module](https://nginx.org/en/docs/http/ngx_http_split_clients_module.html)).

A further improvement could be to use a UNIX socket for the server that will proxy requests to `mirror-upstream`,
on some OSes it might be faster than TCP proxying. I didn't test it yet but it looks like a possible combination offered
by `listen` and `proxy_pass` (the URI rewrite _should_ be possible).

### Conclusions

Traffic mirroring is a key enabler for so many scenarios that it should be a first class citizen in any major webserver, for every protocol.

NGINX configuration can be treated as an art, given the infinite possibilities the server offers with its modularity.

I hope this post will save you some hours of debugging and will increase your testing capabilities,
I decided to write it because I could not find a complete answer online do the question: "can NGXIN mirror gRPC traffic?".
And the answer is "yes"!
