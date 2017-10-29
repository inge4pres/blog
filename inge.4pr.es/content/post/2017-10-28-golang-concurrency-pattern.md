---
title: "Golang Concurrency Patterns"
date: 2017-10-28T19:40:09+02:00
author: "inge4pres"
description: "Some interesting insights on concurrency primitives"
tags:
  - golang
  - development

categories:
  - tech

draft: true
---

In the early days of Go the language was often tailored towards "system programming" due to its C-stlye syntax and ability to write high-performance applications. Not long time after Go adoption was starting to gain traction for distributed systems development and projects like etcd, docker and kubernetes revealed the power of the networking capabilities offered by the internals. Along the way a lot of libraries have been built around the powerful primitives offered by Go but in my opinion there is not enough use literature around the _[Communicating_sequential_processes](https://en.wikipedia.org/wiki/Communicating_sequential_processes)_ implementation available through channels and goroutines, they are not even widely used in the standard library. I'll detail here some concurrency patterns that I found useful and hopefully they'll be idiomatic enough to represent a good use case for any use case.

##### A premise
CSP it's kind of a similar feature to threading which I used in Java before using Go, but there are some differences; to know more on CSP I really recommend watching Rob Pike's [excellent talk on the topic](https://vimeo.com/49718712).

#### My experience
Personally it took me a while to find my way out of the issues I ran into when first using concurrency features in Go: they are definitely the most complicated part of Go usage, which is on average simpler that any other language I tried. So for me, the biggest problem was to understand what it means to have a goroutine spawned and how to control its execution or get data out of it, so I put together a list of examples on how concurrent programs flow can be controlled with the primitives built in the language.

#### Channel, channels, channels everywhere
A channel in Go is a way to pass messages between functions and goroutines, the official definition from [A Tour of Go](https://tour.golang.org/concurrency/2) is:

_Channels are a typed conduit through which you can send and receive values with the channel operator, `<-`_

So what are they good for? They are actually not very helpful without goroutines: a goroutine is a lightweigth thread managed by the Go runtime ([definition](https://tour.golang.org/concurrency/1)), think like a background process that can be spawned and does not need to be managed directly by you, I like the concept of _"run and forget"_.
The easiest concurrency pattern available is thinking of a goroutine processing some data in the background and returning them through a channel to the main thread executing our code; this can be very powerful and scale well to multiple functions and channels.

#### Synchronization via WaitGroups
[WaitGroups](https://golang.org/pkg/sync/#WaitGroup) are part of the `sync` package from the Go standard lib: they are a way for waiting the execution of goroutines to ends properly and ensure all the work done in the background is completed. WaitGroups are often used with `defer` 

#### References
[Concurrency made easy - Dave Cheney](https://www.youtube.com/watch?v=yKQOunhhf4A)
[Share memory by communicating - Codewalk](https://golang.org/doc/codewalk/sharemem/)
[Go channels are bad and you should feel bad](http://www.jtolds.com/writing/2016/03/go-channels-are-bad-and-you-should-feel-bad/)
