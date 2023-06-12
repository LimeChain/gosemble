---
layout: default
title: Debug
permalink: /test/debug
---

To aid the debugging process, there is a set of imported functions that can be called within the Runtime to log any message via the Host.

```go
func Critical(message string) // logs and aborts the execution
func Warn(message string)
func Info(message string)
func Debug(message string)
func Trace(message string)
```