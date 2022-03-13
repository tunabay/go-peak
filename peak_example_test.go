// Copyright (c) 2022 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package peak_test

import (
	"fmt"
	"time"

	"github.com/tunabay/go-peak"
)

func Example_usage() {
	// Value that tracks the maximum/minimum in last 1 sec, with the initial
	// value 1000.
	v := peak.New[uint32](time.Second, 1000)

	// Add, Sub, and Set change the value.
	time.Sleep(time.Millisecond * 250) // [0.25s]
	v.Add(300)                         // [0.25s] 1300
	time.Sleep(time.Millisecond * 250) // [0.50s]
	v.Sub(1100)                        // [0.50s] 200
	time.Sleep(time.Millisecond * 250) // [0.75s]
	v.Set(900)                         // [0.75s] 900
	time.Sleep(time.Millisecond * 740) // [1.49s]

	// Get the current value and maximum/minimum values in last 1 sec.
	cur, min, max := v.Get()
	fmt.Printf("cur: %3v\n", cur)
	fmt.Printf("min: %3v\n", min)
	fmt.Printf("max: %3v\n", max)

	// Output:
	// cur: 900
	// min: 200
	// max: 900
}

func ExampleValue() {
	v := peak.New[uint16](time.Millisecond*500, 100)

	v.Sub(50)                          // 0.0s: 100 -> 50
	time.Sleep(time.Millisecond * 200) // 0.0s -> 0.2s
	v.Add(30)                          // 0.2s: 50 -> 80
	time.Sleep(time.Millisecond * 200) // 0.2s -> 0.4s
	v.Set(250)                         // 0.4s: 80 -> 250
	v.Sub(50)                          // 0.4s: 250 -> 200
	time.Sleep(time.Millisecond * 200) // 0.4s -> 0.6s

	cur, min, max := v.Get() // min/max within 0.1s .. 0.5s
	fmt.Printf("cur: %3v\n", cur)
	fmt.Printf("min: %3v\n", min)
	fmt.Printf("max: %3v\n", max)

	// Output:
	// cur: 200
	// min:  80
	// max: 250
}
