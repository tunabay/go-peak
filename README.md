[![Go Reference](https://pkg.go.dev/badge/github.com/tunabay/go-peak.svg)](https://pkg.go.dev/github.com/tunabay/go-peak)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

# go-peak

## Overview

go-peak is a Go package providing a generic data type that tracks the maximum
and minimum peak values within a specific period of time.

### Background

go-peak was originally designed for Prometheus monitoring to catch the peaks
between scraping. For example, suppose the scraping interval is configured to
one minute, and a Gauge metric spikes immediately after a scraping. If this
Gauge value had dropped by the next scraping, that spike value would not be
detected. If there are secondary metrics that report the maximum and/or minimum
values for the last 1 minute, these peaks will be able to be detected.

## Usage

```go
import (
	"fmt"
	"time"

	"github.com/tunabay/go-peak"
)

func main() {
	// uint32 Value that tracks the maximum/minimum in last 1 sec, with the
	// initial value 1000.
	v := peak.New[uint32](time.Second, 1000)

	// Add, Sub, and Set change the value.
	time.Sleep(time.Second / 4) // [0.25s]
	v.Add(500)                  // [0.25s] 1500
	time.Sleep(time.Second / 4) // [0.50s]
	v.Sub(700)                  // [0.50s] 800
	time.Sleep(time.Second / 4) // [0.75s]
	v.Set(300)                  // [0.75s] 300
	time.Sleep(time.Second / 4) // [1.00s]

	// Get the current value and maximum/minimum values in last 1 sec.
	cur, min, max := v.Get()
	fmt.Printf("cur: %3v\n", cur)
	fmt.Printf("min: %3v\n", min)
	fmt.Printf("max: %3v\n", max)
}
```
[Run in Go Playground](https://go.dev/play/p/MPCpKkrw_UG?v=gotip)

## Limitation

- It uses the new Go generics syntaxes and requires Go v1.18 or higher.
- The detection period is not accurate. For example, a report for the last one
  second may contain a maximum value at 1.001 seconds ago. This is because not
  all individual changes are recorded, but at a resolution of about 1/128 to
  1/256 of the reporting period.
