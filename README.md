## go-stream

[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/chinalhr/go-stream)](https://img.shields.io/github/go-mod/go-version/chinalhr/go-stream)
[![Release](https://img.shields.io/github/v/release/chinalhr/go-stream.svg?style=flat-square)](https://github.com/chinalhr/go-stream)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/ChinaLHR/go-stream/Build)](https://github.com/ChinaLHR/go-stream/actions/workflows/build.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/chinalhr/go-stream)](https://goreportcard.com/report/github.com/chinalhr/go-stream)
[![codecov](https://codecov.io/gh/chinalhr/go-stream/branch/main/graph/badge.svg?token=ZHMPMQP0CP)](https://codecov.io/gh/chinalhr/go-stream)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/chinalhr/go-stream/blob/main/LICENSE)

## Overview
Go-Stream is a stream processing library to implement the Java Stream API with Go.

## Features
- `non-storage`: Stream is not a data structure, but the view of the data source. The data source can come from slice, map or supplier function.
- `pipeline`: Operate the elements through pipeline, support short-circuit and parallel execution.
- `lazy-evaluation`: The intermediate operation on the stream is lazy, and it will be truly executed only when the terminal operation is performed.

Go-Stream supports the following operations

| Stream Operation            |                      |                                                              |
| --------------------------- | -------------------- | ------------------------------------------------------------ |
| **Intermediate operations** | Stateless            | Filter、Map、Peek、FlatMap                                   |
|                             | Stateful             | Distinct、Sorted、Skip、Limit、TakeWhile、DropWhile          |
| **Terminal operations**     | non short-circuiting | ForEach、Reduce、ReduceFromIdentity、Count、Max、Min、FindLast、ToSlice、ToMap、GroupingBy |
|                             | short-circuiting     | AllMatch、AnyMatch、NoneMatch、FindFirst                     |

## Quick Start
1. installation go-stream library

```
go get github.com/chinalhr/go-stream
```

2. import

```
import "github.com/chinalhr/go-stream"
```

3. Use Stream to process data

```go
	type widget struct {
		color  string
		weight int
	}

	widgets := []widget{
		{
			color:  "yellow",
			weight: 4,
		},
		{
			color:  "red",
			weight: 3,
		},
		{
			color:  "yellow",
			weight: 2,
		},
		{
			color:  "blue",
			weight: 1,
		},
	}

	sum := stream.OfSlice(widgets).
		Filter(func(e types.T) bool {
			return e.(widget).color == "yellow"
		}).
		Map(func(e types.T) (r types.R) {
			return e.(widget).weight
		}).
		Reduce(func(e1 types.T, e2 types.T) types.T {
			return e1.(int) + e2.(int)
		})
```