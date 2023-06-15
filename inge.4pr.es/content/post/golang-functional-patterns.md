---
date: "2019-08-07T15:04:52+02:00"
title: "Golang-Functional-Patterns"
author: "inge4pres"
emoji: true
comment: true

categories:
  - tech
  - work

tags:
  - golang
  - functional
  - programming

draft: true
---

Often I get asked by coworkers why Go does not have language builtins for filtering a slice or map/reducing a collection of structs; my instinct at first is to reply with **"because you don't need them!"** but then I try to be more gentle and softly go with a "you can implement your own".
One of the core principles in the Go philosophy is _simplicity_ and most of the functional patterns that come from working with Streams in the JVM world are far from simple.
Plus, for the average software craftsman, approaching functional programming means having to deal with complex category theory concepts and shifting away from the comfortable way of modeling software like real-world objects interacting one with each other.
Go does not have this dichotomy baked in the language, it gives you enough freedom to allow both object-oriented and functional programming and comes with enough power to still have good performance from both styles; furthermore you are also given all the tooling to understand which solution works better (`testing` package with benchmark functions).

Here are some example of rewriting object-oriented software into more functional style with Go!

### Filter a slice

This is by far the one I get asked the most... "Why is not built into the language?". 
Here's the object-oriented version (with a little irony on GoF patterns)

```go
type Obj struct {
	Strings  []string
	criteria *criteria
}

type criteria struct {
	filter func(interface{}) bool
}

func (c *criteria) WithLengthBiggerThan(number int) *criteria {
	c.filter = func(in interface{}) bool {
		n := in.(int)
		return n > number
	}
	return c
}

func (o *Obj) Filter() {
	for i, s := range o.Strings {
		if o.criteria.filter(len(s)) {
			o.Strings = append(o.Strings[:i],o.Strings[i+1:]...)
		}
	}
}

func (o *Obj) CriteriaFactory() *criteria{
	c := &criteria{}
	o.criteria = c
	return c

}
```

Because the criteria for filtering is an object too, it requires a factory-builder to create it! (I skipped the irony on FilterRepositoryInterface for testing 😅) 
Here's the functional version:

```go
type Check func(interface{}) bool

func Strings(in []string, f Check) []string {
	r := make([]string, 0)
	for _, s := range in {
		if f(len(s)) {
			continue
		}
		r = append(r, s)
	}
	return r
}

func BiggerThan(b int) Check {
	return func(i interface{}) bool {
		return i.(int) > b
	}
}
```

Main difference: the functional version uses immutability and makes the check a higher order function; you can do that in Go because [functions are types](https://tour.golang.org/moretypes/24)!
You can argue which is the best version, the most readable and that is conveying the message in the proper way, but I prefer to look at how it performs.

```bash
=== RUN   Test_filterObject
--- PASS: Test_filterObject (0.00s)
=== RUN   Test_filterFunctional
--- PASS: Test_filterFunctional (0.00s)
goos: darwin
goarch: amd64
pkg: github.com/inge4pres/golang-functional-patterns
Benchmark_filterObject-8       	10000000	       188 ns/op	      80 B/op	       9 allocs/op
Benchmark_filterFunctional-8   	10000000	       450 ns/op	     312 B/op	      12 allocs/op
PASS
ok  	github.com/inge4pres/golang-functional-patterns	6.406s
``` 

Ouch!
From the benchmarks, we can see that the object-oriented version performs better:
it's faster, lighter and performs fewer allocations! But why?
