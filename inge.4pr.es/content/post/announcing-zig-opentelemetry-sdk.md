---
title: "Announcing Zig OpenTelemetry SDK"
author: "inge4pres"
type: ""
date: 2025-10-28T12:00:00+01:00
subtitle: The first alpha release is out!

categories:

- tech

tags:

- software-engineering
- observability
- zig
- opentelemetry
- open-source
---

Last Sunday I shared on social media my excitement for the [first alpha release](https://github.com/zig-o11y/opentelemetry-sdk/releases/tag/v0.1.0) of the Zig OpenTelemetry SDK!
This is the culmination of an investment of my personal time over the past 15 months, and I am so happy that it is finally usable.

## Why OpenTelemetry for Zig?

If you've read my [previous post about why Zig is great for observability](../zig-is-great-for-observability), you already know my point of view on Zig's suitability for building observability tools.
The combination of explicit memory management, performance-first design, and seamless C interoperability makes Zig an ideal choice for implementing observability systems.

Aside from Zig's technical merits, there are three fundamental reasons why I wanted to build this SDK:

### Cloud-native needs Zig 

I wish more CNCF projects appreciate Zig's design philosophy and values: handling edge cases, be as efficient as possible, ultimately better serving users.
Specifically for OpenTelemetry, Zig aligns perfectly with the requirements of observability tools.

The language gives you:
- Precise control over memory allocations to minimize overhead
- Performance-first standard library with data structures like `MultiArrayList` that excel at columnar data
- Complete C/C++ interoperability for deep system integration
- Native eBPF compilation for kernel-space instrumentation

### OpenTelemetry is the future

[OpenTelemetry](https://opentelemetry.io/) has emerged as the industry standard for observability.
It provides vendor-neutral APIs and SDKs for generating, collecting, and exporting telemetry data (traces, metrics, logs and profiles).
It also has a vast community.

Having a Zig implementation means:
- Zig developers can instrument their applications with industry-standard primitives
- C developers can use the SDK through Zig's C ABI compatibility (take this with a grain of salt, we still have [some hurdles](https://github.com/zig-o11y/opentelemetry-sdk/issues/56))
- The observability ecosystem becomes more accessible to systems programmers

### Zig helps developers reason about resources

One of Zig's greatest strengths is how it makes resource management explicit.
When we're building observability instrumentation that runs in production 24/7, it is vital to understand exactly what resources are consumed.

Zig's allocator-passing pattern and explicit error handling make it clear:
- Where allocations happen and with what allocator
- What can fail and why
- What resources need to be cleaned up

This is invaluable when building tools that must have minimal and predictable overhead.

Most importantly, when running in Kuberentes, we might even cap the total allocated memory by the application with the configured container memory limits. 

Imagine a container that is never OOM'd by the kernel. This is not only possible, but also practical with a binary compiled with Zig, as we can decide how to interpret the allocation error.

## A long journey

This project has been a labor of love that started back in July 2024. I began working on the metrics implementation by myself, exploring how to best represent OpenTelemetry's data model in Zig while taking advantage of the language's unique features.
At first, it was an excuse to learn more about tagged unions, mutability and memory management; it turned into a pleasant codebase that I wanted to share.

By January 2025, I had put together a couple of people that was sharing the same interest and felt confident enough to create a [community card](https://github.com/open-telemetry/community/issues/2514) in the OpenTelemetry GitHub organization to gauge interest.
I wanted to see if others in the OpenTelemetry community would be interested in maintaining an official Zig implementation. The response was encouraging, and it helped attract contributors to the project.

Around that time, [@Drumato](https://github.com/Drumato) (Yamato) and [@kakkoyun](https://github.com/kakkoyun) (Kemal) joined the effort.
Yamato had already been working on a traces implementation on his own, and Kemal brought deep OpenTelemetry expertise from his work on other OpenTelemetry projects. Together, we pushed forward on implementing all three signals: traces, metrics, and logs.

Shortly after I asked my ex-coworker [@kmos](https://github.com/kmos) (Giovanni) to join us and when we onboarded benchmarks [@hendriknielaender](https://github.com/hendriknielaender) also joined the effort.

**To all the people that contributed, a deep thank you. You helped realizing my dream of a working Zig OpenTelemetry SDK**.

The journey wasn't always smooth. We hit challenges with:
- Zig's pre-1.0 instability and breaking changes between releases
- Figuring out the right abstractions for OpenTelemetry's specification in Zig's type system
- Designing APIs that felt natural to Zig developers while staying true to the specification
- Building comprehensive tests and integration with the OpenTelemetry Collector

## What's in this first release?

This alpha release includes:

- **Complete Stable Signal Support**: Full implementations of Traces, Metrics, and Logs APIs and SDKs
- **OTLP Protocol Exporters**: HTTP/Protobuf and HTTP/JSON exporters with gzip compression support
- **Production-Ready Examples**: 10+ working examples demonstrating real-world usage patterns
- **Integration Tested**: Verified against the official OpenTelemetry Collector
- **Performance Benchmarks**: Comprehensive benchmark suite for performance tracking

You can find the complete [documentation here](https://zig-o11y.github.io/opentelemetry-sdk/).

## We need you!

This is just the beginning. To make this SDK production-ready, we need **users** to:

- Try it in your Zig projects
- Report bugs and edge cases we haven't encountered
- Suggest improvements to the API design
- Contribute missing features or optimizations
- Help us understand real-world usage patterns

You can:
- Report issues on [GitHub Issues](https://github.com/zig-o11y/opentelemetry-sdk/issues)
- Contribute improvements via Pull Requests
- Join discussions in the Zig Discord

## Looking Forward

This alpha release is a milestone, but there's much more work ahead.

We're working toward a stable 1.0 release that will include:
- gRPC protocol support
- Profile signal implementation
- Additional exporters and integrations
- Performance optimizations
- More comprehensive documentation and examples

If you're interested in Zig, Observability, or OpenTelemetry, I'd love to have you join us on this journey. 
Together, we can build observability tools that embody Zig's philosophy: simple, readable, and efficient.

Let's reduce that observability tax, one allocation at a time.
