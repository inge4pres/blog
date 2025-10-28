---
title: "Zig is great for Observability"
author: "inge4pres"
type: ""
date: 2025-10-11T09:00:00+02:00
subtitle: Why Zig is the most suited language for Observability

categories:

- tech

tags:

- software-engineering
- observability
- zig
- performance
---

After years of building observability tools and working with various languages, I believe I am i the position to express a strong opinion: [Zig](https://ziglang.org) is the most suited language for building Observability systems.

This is a bold claim, especially in a field dominated mostly by Go (hello OpenTelemetry) and more recently Rust. But after exploring Zig's design philosophy and capabilities, I have developed a sense of confidence that Zig is the right language and ecosystem for me to build observability tools.

## The observability tax

If you've ever worked on observability infrastructure, you may have heard the concept of "observability tax": this is the overhead that monitoring, tracing, and profiling tools add to your systems. Ironically, the very tools designed to help you understand performance can become performance bottlenecks themselves.

Traditional observability agents often consume a lot of the monitored hosts' resources, CPU, network and memory.
This is acceptable when you're debugging a critical issue on a test environment, but problematic when you're running distributed tracing, profiling hudreds of processes or producing high-cardinality metrics in production at scale.

With regards to backends (middlewares and storage systems), the "tax" is paid in real money spent on Cloud hardware to accumulate massive amounts of data that have ever-diminishing value.

After 2 years working with Zig, I feel the urge to transpose some of the Zig Zen to the Observability landscape, in order to improve  users' and developers' life with better software.

### Precise control of resources

Zig gives you explicit control over memory allocations, allowing you to minimize te cost of allocations significantly. Unlike languages with garbage collectors that introduce unpredictable pauses, or languages that hide allocations behind abstractions, Zig makes every allocation visible and controllable.

```zig
const std = @import("std");

pub fn processMetrics(allocator: std.mem.Allocator, data: []const u8) !void {
    // Explicit allocation - no hidden costs
    var metrics = try std.ArrayList(Metric).initCapacity(allocator, 1000);
    defer metrics.deinit();

    // Process with precise control over resources
    for (data) |byte| {
        // Every allocation is explicit and controllable
    }
}
```

In my experience working on continuous profiling, reducing allocator overhead can translate to significant cost savings at scale. Zig's allocator-passing pattern makes this optimization straightforward and explicit.

## Performance-first standard library

Zig's standard library is designed with performance as a first-class concern, and it includes some truly remarkable data structures for observability workloads.

### MultiArrayList

One of Zig's hidden gems is [`std.MultiArrayList`](https://ziglang.org/documentation/master/std/#std.MultiArrayList), a data structure that stores struct fields in separate arrays, essentially an in-memory columnar storage format built into the standard library.
Its peculiarity is that it can reduce memory usage both when allocating the data and when retrieveing them.

Why does this matter for observability? Because most observability data is columnar by nature:

- Time-series metrics: timestamps, values, labels, etc...
- Distributed traces: span IDs, timestamps, durations, service names, etc...
- Logs: timestamps, severity levels, messages, metadata, etc..

See the following code snippet for an example of why MultiArrayList is king.

```zig
const Span = struct {
    trace_id: u128,
    span_id: u64,
    timestamp: i64,
    duration: u64,
};

// Traditional approach: array of structs (AoS)
// Memory layout: [trace_id|span_id|ts|dur][trace_id|span_id|ts|dur]...
var traditional_spans: []Span = ...;

// Zig's MultiArrayList: struct of arrays (SoA)
var spans = std.MultiArrayList(Span){};
try spans.append(allocator, .{ .trace_id = 123, .span_id = 1, .timestamp = 1000, .duration = 50 });
try spans.append(allocator, .{ .trace_id = 124, .span_id = 2, .timestamp = 1100, .duration = 75 });

// Memory layout is columnar - only loads what you need:
// trace_ids:  [123][124]...
// span_ids:   [1][2]...
// timestamps: [1000][1100]...  <- Only this loaded for timestamp queries
// durations:  [50][75]...      <- And this for duration queries

// Query only timestamps - touches only the timestamp array!
var slow_spans: usize = 0;
for (spans.items(.timestamp), spans.items(.duration)) |ts, dur| {
    if (ts > cutoff_time and dur > 100) {
        slow_spans += 1;
    }
}
```

The structure-of-arrays layout provides the developers a tool to design carefully CPU cache utilization when we're aggregating metrics, searching traces by timestamp, or performing other column-oriented operations that are common in magling observability data.

The performance difference can be dramatic. In columnar workloads, `MultiArrayList` can be 2/3 times faster than traditional array-of-structs approaches due to cache efficiency.

## Complete C/C++ interoperability

Observability often requires deep integration with system internals. You might need to hook into:

- Kernel tracing mechanisms (perf, ftrace)
- System libraries (libc, libpthread)
- Low-level interfaces

Zig's C interoperability is not just good, it's seamless.
You can import C headers directly, call C functions without FFI wrappers, and even improve upon C libraries by extending them with Zig code.

```zig
// Direct C header import - no bindings needed
const c = @cImport({
    @cInclude("linux/perf_event.h");
    @cInclude("sys/ioctl.h");
});

pub fn setupPerfEvent() !c_int {
    var pe: c.struct_perf_event_attr = undefined;
    @memset(@as([*]u8, @ptrCast(&pe))[0..@sizeOf(c.struct_perf_event_attr)], 0);

    pe.type = c.PERF_TYPE_HARDWARE;
    pe.config = c.PERF_COUNT_HW_CPU_CYCLES;
    pe.disabled = 1;

    const fd = c.perf_event_open(&pe, 0, -1, -1, 0);
    if (fd < 0) return error.PerfEventOpenFailed;

    return fd;
}
```
This matters because observability tools need to interact with operating system primitives frequently. 
The ability to call C functions directly, without FFI overhead or binding generation, means you can build high-performance tools that integrate seamlessly with existing system infrastructure.

There are plans in Zig to actually _remove_ `@cInclude` and only rely on the build system to link C symbols into the exeutables.
I trust the Zig core team and the community will take the steps in the most sensible direction to keep working in Zig with existing C codebases as pleasing as it is today.

The Zig compiler is also able to compile (and cross-compile!) C/C++ code without too much boilerplate, and by configuring compilation properties directly in the build system.

## eBPF

Here's where it gets really exciting: one of the many Zig compilation target is the BPF bytecode.

[eBPF](https://ebpf.io) has become the foundation of modern observability in Linux.
It allows you to run sandboxed programs in the kernel without modifying kernel code or loading kernel modules.

Traditionally, writing eBPF programs means writing C/C++ code with restricted features, having to deal with LLVM, and struggling in debugging the programs.

I believe Zig provides a great developer experience for eBPF, as the core logic of programs can often be tested like a regular Zig program, relying on the contract from the kernel in tests and actually using the kernel interface only at runtime.

The snippet below is an example of a Zig BPF program that traces the execution of a syscall. 

```zig
const std = @import("std");
const BPF = std.os.linux.BPF;

// GPL License is mandatory to install in kernel
export const _license linksection("license") = "GPL".*;

// Kprobe that traces when processes read from an FD
export fn trace_ksysread(ctx: *anyopaque) linksection("kprobe/ksys_read") c_int {
    _ = ctx;
    const pid = BPF.kern.helpers.get_current_pid_tgid() >> 32;
    const fmt = "execve called by PID: %d\n";
    _ = BPF.kern.helpers.trace_printk(fmt, fmt.len, pid, 0, 0);
    return 0;
}
```

When saved as a file `trace.zig` you can compile it to BPF with:

```bash
zig build-obj -target bpfel-freestanding -O ReleaseFast trace.zig
```

This produces a `trace.o` object that can be loaded into a Linux kernel.

This is a silly example (basically the eBPF "Hello world"), nothing close to something you would deploy in production.
Yet it shows the ability to write kernel-space eBPF programs in Zig, with the same language you'll write the userpsace part, which I believe is a game-changer for the eBPF-based observability agents developer experience.

Like in other programming environments, you can build a complete agent where eBPF programs collect data in kernel space and a userspace program loads the object, manages the maps and processes and aggregate data.
With Zig though, you can do that in a single language, with a build system also written in the same language, having consistent tooling and overall less cognitive load.

## What's wrong with Go or Rust?

You might be wondering: "What about Go, the current favorite language to build Observabilty tools in?"
and "How dare you not mention Rust as excellent tool for this!".

Following here a list of **personal opinions** that I hope will resonate with you.

Go is my go-to language for many projects, and it's widely used in observability (Elastic, Prometheus, Grafana, OpenTelemetry). But:
- The garbage collector introduces unpredictable latencies; when trying to optimize them you go with `sync.Pool` and end up [like this](https://github.com/open-telemetry/opentelemetry-collector/issues/13954)
- Runtime overhead makes it less suitable for high-frequency data collection
- Limited control over memory layout (no equivalent to `MultiArrayList`)
- Cannot compile directly to eBPF, neither interoperate with C; DX is decent with libraries like [`cilium/ebpf`](https://github.com/cilium/ebpf) or [acquasecurity/libbpfgo](https://github.com/aquasecurity/libbpfgo)

Rust is excellent and has strong memory safety guarantees. But:
- Its complexity can slow down development: personally I feel very, very dumb when I have to interpret other people's Rust code
- Compile times are significantly longer, due to its added ref counting checks at compile time
- The borrow checker, while powerful, can be restrictive for certain observability patterns where you want to share memory freely
- eBPF support is still meh IMO, probably [aya-rs](https://github.com/aya-rs/aya) still the most (only?) interesting project

I feel like for the eBPF use case Zig finds the sweetest spot: the modern tooling of Rust, the simplicity and development speed of Go, and the performance and control of C.

## Zig is not perfect either

Before you rush to rewrite all your observability tools in Zig, it's important to acknowledge the real challenges and trade-offs.

### Pre-1.0 instability

Zig is still pre-1.0, and that's not just a version number. Breaking changes are likely to happen with every minor release.
The recent 0.15 release introduced significant API changes, and if you're building a long-term project, you need to be comfortable with occasional rewrites or be ready to pin to a specific Zig version (which I don't recommend).

This instability makes it harder to maintain libraries and can be frustrating when you upgrade and find your code no longer compiles. The trade-off is that the language is actively and rapidly improving, but it comes at a cost to stability.

### Manual memory management

While I've praised Zig's explicit memory management, let's be honest: it's also a source of bugs.
Zig provides *spatial* memory safety (buffer overflow protection) but not *temporal* memory safety (use-after-free protection). In practice, [experienced developers report](https://lobste.rs/s/xnyrve/memory_safety_features_zig) that it's easier to introduce memory safety bugs or create resources' leaks in Zig codebases, compared to rare occurrences in Go or Rust.

Zig's `std.heap.DebugAllocator` can catch memory leaks, double frees and use-after-frees, but you would not use such allocator in production due to its overhead, so you'd be running without those safety nets.
This means memory bugs can slip through to production if you're not extremely disciplined with your allocator usage, and use development practices that improve detection of such issues in the IDE.

For observability tools that run in production 24/7, this is a real risk. A memory leak in a long-running agent or a use-after-free in a high-throughput collector can be catastrophic.

### Immature ecosystem

The Zig package ecosystem is young. While the language has excellent fundamentals, you'll often find yourself:
- Writing bindings to C libraries yourself
- Finding that the library you need doesn't exist yet
- Encountering incomplete or abandoned packages
- Struggling with sparse documentation

Compared to Go's massive standard library and mature ecosystem (including well-supported libraries like cilium/ebpf), or Rust's crates.io with thousands of production-ready crates, Zig is a little behind.

This is particularly relevant for observability, where you might want to integrate with various protocols, storage backends, or cloud services. In Go, there's likely already a well-maintained library; in Zig, you might be the first.

### eBPF support is still maturing

While I've highlighted Zig's ability to compile to BPF bytecode, the reality is that the tooling and ecosystem are still immature compared to the battle-tested C + libbpf workflow.

The challenges include:
- Limited high-level abstractions for common eBPF patterns (compared to libraries like cilium/ebpf or aya-rs)
- Debugging eBPF programs is already difficult; doing it in a pre-1.0 language adds another layer of complexity
- Most eBPF examples, documentation, and community knowledge are in C
- CO-RE (Compile Once - Run Everywhere) support in Zig is less mature than in C

Projects like [zbpf](https://github.com/tw4452852/zbpf) show promise, but they're not yet at the maturity level of the established C ecosystem. If you're building production eBPF tools today, you'll likely hit rough edges.

### Being an early adopter

Choosing Zig for production observability tools in 2025 means accepting that you're an early adopter. You might encounter bugs in the stanrd library, missing features, and moments of "is this a bug in my code or in the compiler?"

You'll also face challenges in hiring and onboarding: Zig developers are rare, and training people on a pre-1.0 language with frequent breaking changes is not cost effective.

Despite these challenges, projects like [Bun](https://bun.sh) and [TigerBeetle](https://tigerbeetle.com) have bet on Zig for production systems and are seeing success. But they've also invested heavily in working around these limitations.

## Reducing the observability tax

Let me bring this back to where we started: the observability tax.

I think Zig is really suited to build observability tools that:
- have better visibility on the allocated memory than equivalent Go tools
- process more data per CPU second (no runtime overhead, better cache utilization, explicit control)
- run in both kernel and userspace with the same language
- integrate seamlessly with existing C/C++ infrastructure

This isn't theoretical. I remember when at Optimyze and we were building Prodfiler, the mandate was that the agent should _never_ take more than 1% the host CPU, and possibly work using less than 250 MB memory.
If we had done the agent in Zig instead of Go, we wouldn't had to use sync.Pool and caches here and there to avoid CPU-costly garbage collection cycles; we could've used an arena allocator backed by a fixed size of memory, and every time we would have filled it up we could've handled allocation error.

I think that's the spirit that we should have when producing observability tools: un-obtrusive, reliable and predictable tools.

Note: **I don't regret a single second of my worktime at Optimyze**, even if Go was the tool of choice at the company. It was the best choice at that time. Optimyze was the most professionally-enriching and fun work environment I've ever been in.

## Conclusion

Observability is fundamentally about understanding systems with minimal overhead. Zig's design philosophy of expressive control flow, performance by default, seamless interop with the C ABI, and systems-level capabilities align perfectly with these requirements.

As we push Observability toward large, higher-cardinalty structured events, continuous profiling, always-on distributed tracing, and more interconnected signals, the tools we use need to be more efficient. Zig is positioned to be the language that powers the next generation of the Observability infrastructure.

Is Zig perfect? No. It's still pre-1.0, the ecosystem is young, and there are rough edges. But the fundamentals are sound, and for observability workloads, those fundamentals matter immensely.

I startedy my Zig journey a while ago and along the way I decided to create an OpenTelemetry SDK for Zig.
The [project](https://github.com/zig-o11y/opentelemetry-sdk) is still incomplete after 1.5 years in the making, and if you'd like to contribute we would love to get newcomers onboard. 

Even outside of OpenTelemetry, if you're building some observability tools, I encourage you to give Zig a serious look.

You might be surprised at what it enables.
