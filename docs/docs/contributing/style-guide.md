---
layout: default
title: Style guide
permalink: /contributing/style-guide
---

# Style guide

We try to follow common practices described in [Effective Go](https://go.dev/doc/effective_go) and [Google Style Guide for Go](https://google.github.io/styleguide/go/) to maintain consistency across the codebase. Following a consistent coding style makes the code more readable and maintainable for everyone.

- **Custom errors**: New polkadot-related custom errors should implement the error interface. See [#271](https://github.com/LimeChain/gosemble/issues/271) and linked PRs.
- **Error handling**: We use critical logging for resolving errors(panicing). We add logger to modules through dependency injection and we only resolve errors(critical log) in the api modules. See [#315](https://github.com/LimeChain/gosemble/pull/315).
- **SCALE types**: Certain [SCALE types](../overview/runtime-architecture.md#scale-codec) like the Result types are common in Rust which the original Substrate implementation is based on, but add unneeded complexity in Go. When we have these types as return types it's preferred to encode the data into the required type as late as possible - just when you need to return it in the api module, instead of building your logic around them. See [#322](https://github.com/LimeChain/gosemble/pull/322).