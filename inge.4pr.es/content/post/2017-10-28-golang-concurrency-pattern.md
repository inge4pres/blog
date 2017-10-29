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

#### WaitGroups
[WaitGroups](https://golang.org/pkg/sync/#WaitGroup) are part of the `sync` package from the Go standard lib: they are a way for waiting the execution of goroutines to end properly and ensure all the work done in the background is completed. WaitGroups are often used with `defer` to fill in the wait queue when the goroutine exits.

#### Some examples
For me the most difficult thing to understand when approaching concurrency was how to ensure all of my goroutines completed execution: to do this the easiest way is using WaitGroups as in [waitgroup_test.go](https://github.com/inge4pres/blog/blob/master/golang-concurrency-patterns/waitgroup_test.go): `wg.Add(1)` adds one item in th wiat queue and `wg.Done()` remove one item; using `wg.Wait()` in the main process makes the process wait until the wait group is emptied.
If you run the tests with `go test -v . -race  -run ^TestWaitGroup` you can see the execution time when using concurrency or not. Changing the value of `ops` variable in [functions_test.go](https://github.com/inge4pres/blog/blob/master/golang-concurrency-patterns/functions_test.go#L3) will make the tests process less or more items.

With channels there are more features and gotchas that need to be taken into account:
* a read from a closed channel returns the type's zero-value
* a send to a closed channel will `panic`
* a read and a send alone to an unbuffered channel are blocking: they will generate a deadlock if there is not a corresponding send/read operation on the other side of the channel
* a send on a buffered channel will block when the buffer is full and no other read is happening on the other side

That being said, there are a couple of notable usages that I like to include in my concurrency-enabled Go software: the [fan-out pattern](https://github.com/inge4pres/blog/blob/master/golang-concurrency-patterns/channels_test.go#L22) where an input generates multiple goroutines that perform operations in the background and the output of the concurrent goroutines is fetched by a channel in the main thread. Another pattern is [fan-in](https://github.com/inge4pres/blog/blob/master/golang-concurrency-patterns/channels_test.go#L36): multiple functions can return values to a channel as long as the type is consistent. Run tests with `go test -timeout 30s -run ^TestChannelBuffered` to see fan-out/fan-in patterns in action.

Another interesting feature is powered by the `select` statement: you can read from multiple channels in the same function and define behavior for any given channel message, it is another sample of fan-in pattern. Using `select` will block until one of the send/receive operation is available, the operation gets chosen randomly if multiple are available at the same time. `select` has a similar syntax to `switch` so `case` and `default` are the scenario selector. Running [multiple channels](https://github.com/inge4pres/blog/blob/master/golang-concurrency-patterns/channels_test.go#L55) lets you manage multiple types in a single point: run the test with `go test -timeout 30s -run ^TestMultipleChannelsSelect$` to check the execution.

#### Conclusions
My experience with Go concurrency primitives is still forming, I hope I can read and experiment more on the topic as it's one of Go's most powerful and at the same time less documented features! I'd really love to hear feedback from the read so if you get up to this point, take a step forward and leave a comment below, I'd really love to discuss.

#### References
[Worker pools](https://gobyexample.com/worker-pools)
[Channels](https://golangbot.com/channels/)
[Concurrency made easy - Dave Cheney](https://www.youtube.com/watch?v=yKQOunhhf4A)
[Share memory by communicating - Codewalk](https://golang.org/doc/codewalk/sharemem/)
[Go channels are bad and you should feel bad](http://www.jtolds.com/writing/2016/03/go-channels-are-bad-and-you-should-feel-bad/)
