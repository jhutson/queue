# Queue Package

[![Build Status](https://github.com/jhutson/queue/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/jhutson/queue/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jhutson/queue)](https://goreportcard.com/report/github.com/jhutson/queue)
[![Go Reference](https://pkg.go.dev/badge/github.com/jhutson/queue.svg)](https://pkg.go.dev/github.com/jhutson/queue)

The `queue` package provides a ring-buffer based implementation of the queue data structure.

The provided `Queue` interface is generic, to allow for varying the element type without resorting to using `interface{}`.

There are unbounded and bounded queue variations available. The unbounded variation resizes its internal storage as needed. The bounded version keeps to a maximum specified capacity.

## Installation
```shell
go get github.com/jhutson/queue
```

## Example
### Creating and using an unbounded queue

```go
package main

import (
	"fmt"

	"github.com/jhutson/queue"
)

func main() {
	// unbounded queue of integers with initial capacity of five.
	q := queue.NewUnboundedQueue[int](5)

	// add elements
	_ = q.Push(1)
	_ = q.Push(2)

	fmt.Println(q.Length()) // prints 2

	// remove first
	top, err := q.Pop()
	fmt.Println(top, err) // prints: 1 <nil>

	// peek next
	newTop, err := q.Peek()
	fmt.Println(newTop, err) // prints: 2 <nil>
	
	newTop, err = q.Pop()
	fmt.Println(newTop, err) // prints: 2 <nil>
	
	newTop, err = q.Pop()
	fmt.Println(newTop, err) // prints: 0 cannot take element from empty queue
}
```

The bounded queue works similarly. The difference is that `Push` will return an error if the queue is full.

