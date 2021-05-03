---
title: 'Golang Message Queue: a simple TCP message bus'
author: "inge4pres"
layout: post
date: 2015-08-02T16:25:54+00:00
categories:
  - social
  - tech
tags:
  - cloud
  - development
  - golang
  - software

---
**[TL;DR]**

I wrote a Pub/Sub message queue in <a href="http://golang.org" target="_blank">Go</a>, branch &#8220;master&#8221; is stable but missing some interesting feature like distributed memory synchronization between nodes in a cluster and encryption. Code at

<a href="https://github.com/inge4pres/gmq" target="_blank">https://github.com/inge4pres/gmq</a>

Being a cloud system engineer, my work is to design and implement distributed systems: one of the key principles in designing such architectures is _decoupling_, which means ensuring the many parts composing the system are able to share informations and complete a sequence of operations without being tied together. You can read more about cloud architectures and decoupling here.

One of the most common scenario in a cloud application is a series of asynchronous operations executed by many nodes on different layers: for example a front end server tier receiving  files and a backend server tier doing analysis on them; a good practice is to have a message queue between the two serving as an orchestration component. Each web server node will post a message in the queue for every files received, each backend node will consume a message from the queue to complete his operations on the files. In this way the two tiers are independent one from each other: in case of backend failure or over-capacity, the web servers will keep receiving files and storing message in the queue. If the two operations where done synchronously, the backend failure would stop the whole system to work.

A lot of _off-the-shelf_ message queue software is already available, but I felt like writing my own would give me a good point of view on system programming with Go, so I wrote it, and the result is pretty awesome too. In a few days I was able to have a configurable message queue storing messages in memory, on filesystem or database (MySQL); communication is based on JSON via TCP, and the server can be configured to support a maximum number of queues, a maximum message length and queue capacity: combination of the configured parameters will have performance effects on the single node.

The roadmap of &#8220;develoment&#8221; branch is:

  * adding cluster mode
  * adding memory synchronization in cluster mode
  * adding encryption: TLS over TCP
  * adding client authentication

As you may have guessed from the above list, security of GMQ is not implemented at the moment, be careful!

Feel free to try it out and give suggestions!

Cheers 😀
