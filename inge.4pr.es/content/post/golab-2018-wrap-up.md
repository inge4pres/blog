---
date: "2018-10-23T18:07:07+02:00"
title: "GoLab 2018: Wrap Up"
author: "inge4pres"
description: "A recap of the best Go conference in Italy"

draft: false

categories:
  - tech
  - culture

tags:
  - conferences
  - golang
---

I'm in the train back from [GoLab 2018](https://www.golab.io/) and I am so happy that I attended this conference!
It's been definitely one of most beautiful con I have attended in Italy, with tremendous speakers from all over the globe like Filippo Valsorda, Eleanor McHugh, Ron Evans and Bill Kennedy among many others; I have to say the organizers were just perfect in everything from the venue setup to the workshop organization, as if the quality of the talks ware not enough.
I was so delighted I want to write a wrap-up immediately with my head full of ideas for the upcoming future.

### Speakers layout
It was 2 dense days with 45 minutes talks and 4 tracks, 2 each day: Patterns, Embedded (day 1); Web, DX (day 2). Before the talks a keynote each day (Eleanor McHugh and Bill Kennedy respectively), 30 minutes for lightning talks at the end of day 1, cocktails and networking at the end of day 2.
    
### The community
Poeple in the Gophers community is just awesome, and the environment was welcoming and warm in every part of the con; in the afternoon of day 1 a panel on diversity and inclusion was hoste by Cassandra Salisbury and it was really good in the sense that was not the regular D&I talk with practices and models to adopt to "pretend"
 that you are inclusive, it was actually a discussion and sharing of stories that enable a good community. Because a good community is inclusive by design, and the Go community is very good. During breaks I had the chance to shake hands and chat with many people I only had virtually met on Twitter before.
 
### Things I learnt
Here's some of the thing I got to know thanks to the amazing talks I saw:

- [CERN](https://it.wikipedia.org/wiki/CERN) uses Go too! They rebuild their DNS service, see [cernops/golbd](https://github.com/cernops/golbd)
- Google created a project called [Flutter](https://github.com/flutter/flutter/) to build native mobile apps for multiple platforms \[thanks [@edoardo849](https://twitter.com/edoardo849)\]
- the go runtime is powerful but also can be tricky in some occasions, especially with closures; when using them =, make sure to clean up the resources they use and ensure local variable are scoped correctly to avoid data races and bugs. Always run tests with `-race` option and use channels and mutexes to orchestrate your pipelines \[thanks [@empijei](https://twitter.com/empijei)\]
- finite state machines are not only academic lecture, you can make them solve actual problems like decoding different variants of base64 \[thanks [@annaopss](https://twitter.com/annaopss)\]
- there are 2 packages for auto-generating Go structs code from an arbitrary input in JSON and XML format: [JSONGen](https://github.com/bemasher/JSONGen) and [XMLGen](https://github.com/dutchcoders/XMLGen)
- because Go is written in Go, you can parse `.go` files before they are compiled and create your own language extensions (macro) to be then re-compiled in a final binary, crazy! \[thanks Max Ghilardi, see [cosmos72/gomacro](http://github.com/cosmos72/gomacro)\]
- [fnProject](https://github.com/fnproject) is a serverless platform to run function-as-a-service on Oracle cloud and on your private cloud on Kubernetes, as I covered it in this [previous post](https://inge.4pr.es/2018/01/30/serverless-on-kubernetes/) 
- you can manipulate network packets directly with Go! The Linux kernel has a feature to let programs in userspace filter packets based on custom logic, so Telefonica develop (and open sourced) a Go library to manipulate packets because C++ with `libevent` hadn't enough throughput, so they moved from connection-based filtering to packet-based filtering (DPI) doubling throughput with 8 times less CPU! \[See [telefonica/nfqueue](https://github.com/telefonica/nfqueue)\]
- you can close a buffered channel immediately after sending all the items into it, and you'll still be able to read from it all the values; only then the `close(channel)` instruction will be magically executed by the runtime. So [this is valid](https://play.golang.org/p/3YYQo2WK37R) \[thanks [@goinggodotnet](https://twitter.com/goinggodotnet)\] 
- instrumenting and monitoring applications in Kubernetes can be fun! When building an operator you can interact with the exported metrics directly in the controller, you can enrich your systems with distributed traces thanks to OpenTracing/OpenCensus
- using [gRPC](https://grpc.io/) you can auto-generate Swagger documentation from protobuf service definitions, thus auto-generating web-browseable documentation of you service \[thanks [@pawel_slomka](https://twitter.com/pawel_slomka)\]  
- `io.Copy()` from the standard library is the most efficient way of sending network data from 2 `net.Conn` instances; this thanks to [interface upgrade](http://avtok.com/2014/11/05/interface-upgrades.html) in the Copy() implementation, ending in leveraging the kernel packet passing in the most efficient way - without involving the application at all! \[thanks [@filosottile](https://twitter.com/filosottile)\]

And these are just a few (IMO the most important) things I got to know and discover these 2 days... What a ride! Can't wait for next year to be back in Florence!

<a class="twitter-moment" href="https://twitter.com/i/moments/1055353623541112832?ref_src=twsrc%5Etfw">GoLab 2018</a> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>
