---
layout: default
title: Build
permalink: /build/build
---

By utilizing our [Toolchain](/concepts/toolchain), there are currently two options to choose from for the GC implementation.
Modify the `GC` environment variable to switch between them.

#### Extalloc GC

It works with the host's external allocator as per specification.

```bash
GC="" make build
```

#### Conservative GC 

It is used only for **development** and **testing** and works by using a different heap base offset from the allocator's one (as a workaround) so that the GC can use a separate heap region for its allocations and not interfere with the allocator's region.

```bash
GC="conservative" make build
```