---
title: "Composition via embedding in Go"
author: "inge4pres"
type: ""
date: 2024-03-03T16:38:28+02:00
subtitle: Unlock simpler design via a lesser known feature

categories:

- tech

tags:

- software-engineering
- golang
- design-patterns
---

I was introduced to type embedding (or embedding, in short) by Filippo Valsorda in one of his talks.

Apparently this is a lesser known (and lesser used) feature of the Go language, although it's been around since the beginning.
I hope I can spread with this post how it works and how it can be pragmatically helpful to create more expressive and reusable code.

### What is embedding?

[Embedding](https://go.dev/doc/effective_go#embedding) is a powerful feature of the Go language that allows to compose types together,
"bundling" new types that inherit the properties or behavior of the embedded types.
These new types can have their own methods, or they can override the methods of the embedded types: 
this feature is called _method overriding_ and it's existence is not often known by many Go programmers.

From whom approaches Go coming from Java, embedding may resemble inheritance, but it's actually not the same thing in terms of syntax, semantics and memory layout.
When doing inheritance in Java, inherited classes are always allocated in memory as a standalone object, and referenced by the parent class.
In Go, the embedded type is not allocated as a new field in the memory representation of the embedding type: the fields of the embedded struct are rather "inlined" in the embedding type as if they were part of the struct.

### Why is it helpful?

On a design perspective, embedding is a way to express the _is-a_ relationship between types, and it's a way to reuse code
and provide consistent behavior across different types.

One of the most common cases where I've found embedding to be very effective is to create [decorators](https://en.wikipedia.org/wiki/Decorator_pattern) 
for types, and [adapters](https://en.wikipedia.org/wiki/Adapter_pattern) for interfaces.
The Go standard library is using embedding as well, for example in the `testing` package, the `T` and `B` types embed the 
[`common` type](https://github.com/golang/go/blob/6f5d77454e31be8af11a7e2bcda36d200fda07c5/src/testing/testing.go#L917-L918)
to be able to provide shared functionalities between tests and benchmarks without duplicating code.

### Example use cases

There may be many more, but I've found embedding to be particularly helpful in the following scenarios.

#### Composition 

In this scenario multiple types share a set of fields and methods, and you want to avoid code duplication.
You can define a single type that embeds the common properties and methods, and then create new types that embed the common type (just like the `testing` package above).

This is especially helpful for serialization and deserialization of data structures that share the same fields.

```go
type Document struct {
    ID int `json:"_id"`
    Filename string `json:"filename"`
}

type WebsitePage struct {
    Document
    URL string `json:"url"`
}

type Image struct {
    Document
    Width int `json:"width"`
    Height int `json:"height"`
}
```

The `Document` type is re-used across multiple entities from different domains. 

Defining interfaces that embed other interfaces is also a way to create a set of methods that are shared across different types.

```go
type Reader interface {
    Read([]byte) (int, error)
}
type Writer interface {
    Write([]byte) (int, error)
}

type ReadWriter interface {
    Reader
    Writer
}
```

This example from the Go standard library shows how the `ReadWriter` interface is composed by the `Reader` and `Writer` interfaces.
Composing interfaces in this way allows types to implement either the `Reader` or the `Writer` interfaces independently, and be consumed also as a `ReadWriter` if it implements both.

I guess composition is the principle backing the Go idiom of trying to have many small interfaces instead of a few big ones.

#### Blanket implementations

When you have a type that needs to implement a set of methods, and you want to provide a default implementation for some of them, you can embed a type that provides the default implementation.
This is helpful to avoid code duplication and to provide a consistent behavior across different types.

Imagine you are writing a database server that writes to disk through a set of helpers; each data structures that wants access to the filesystem will have to write to it.
If you want to keep the definition of the default behavior for writing to disk in a central module, you can define a `DiskWriter` type that embeds the `os.File` type, and then create new types that embed the `DiskWriter` type.

```go
type DiskWriter struct {
    *os.File
}

func (d *DiskWriter) Write(data []byte) error {
    // write to disk in a specific way, overriding the default behavior
	// inherited from *os.File
	// ...
}

type LSMTree struct {
    DiskWriter
	// ...
}

type BTree struct {
    DiskWriter
	// ...
}
```

In this way, the `LSMTree` and `BTree` types will inherit the `Write` method from the `DiskWriter` type, producing the same
default type of write to disk in their default implementation.

This is handy because when a better way to write to disk is found, it can be implemented in the `DiskWriter` type and all the types that embed it will inherit the new behavior.

The same principle applies to generic types that need to implement a set of methods, and you want to provide a default implementation for some of them.

#### Decorators

The decorator pattern is a way to augment the behavior of an object without changing its external dependency.
It is used to provide additional functionality to a type, without modifying how the callers interact with its methods.

When implementing Observability in our applications, we often need to augment the behavior of our types to add tracing, logging, metrics and other observability features.
We can use embedding to create a decorator type that wraps the original type and adds the observability features we need.

This is also helpful to de-couple the observability features from the business logic, and to make the original type more focused on its main responsibility.
Swapping the original type with the decorated type at call sites should be trivial then, as the decorator type implements the same interface as the original type.

In [this example](https://github.com/inge4pres/blog/commit/55cb83de46235c1ea6bc5b859c47f3a06a9b8ffa#diff-79587407adde647cde1cfa6e3848d32eb30381d38893d15de805488f1bec79ff) we have a `BusinessLogic` component that performs some application task by implementing a `BusinessTask` interface, and we want to add observability features to it:

```go
type BusinessTask interface {
	Execute() string
}

type BusinessLogic struct {
    // ...
}

func (o BusinessLogic) Execute() string {
	return fmt.Sprintf("Executing %s operation", o.Name)
}
```

One could be tempted to add observability features directly to the `BusinessLogic` type, but this would couple the observability features to the business logic,
introducing a dependency that is harmful to the evolution of the business domain.

Instead, we can create a `InstrumentedTask` type that embeds the `BusinessTask` interface and adds the observability features:

```go
// InstrumentedTask is a decorator that adds observability to an operation.
type InstrumentedTask struct {
	BusinessTask

	latencyUsec,
	operationsCount int64
}

// GetMetrics returns the number of operations and the cumulative latency in microseconds.
func (io *InstrumentedTask) GetMetrics() (int64, int64) {
	return io.operationsCount, io.latencyUsec
}

// Execute implements the BusinessTask interface
// and adds observability by measuring how many times we ran it
// and keep track of latency through a counter.
func (io *InstrumentedTask) Execute() string {
	start := time.Now()
	defer func() {
		io.operationsCount += 1
		io.latencyUsec += time.Since(start).Microseconds()
	}()

	return io.BusinessTask.Execute()
}
```

Embedding the interface has the advantage that the same instrumentation behavior can be composed with arbitrary types
that implement the `BusinessTask` interface, without the need to modify the original types.
The method override is taking place here to allow the modified behavior to be executed when the `InstrumentedTask` is used.

### Downsides and caveats

Similarly to other tools at our disposal, overusing causes more harm than good.
Too many levels of indirections across multiple types, referencing each other via embedded types will eventually hurt the readability of the codebase.

A complex type hierarchy may also confuse the compiler in case of generic types, leading the method set of the embedded types to be promoted to the embedding type, and potentially causing method name clashes at compile time.
Personally I've never seen this happening in practice, but it's something to be aware of.

### Conclusion

Despite being part of the language since the beginning, embedding is not often used in Go, and in my opinion it's not often talked enough about in Go communities.
The ability to create more lean designs thanks to embedding makes Go a better language for software engineering, from startups to enterprise users.

Happy hacking Gophers!
