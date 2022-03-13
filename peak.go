// Copyright (c) 2022 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package peak

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/exp/constraints"
)

// Number is a constraint that permits any integer and floating point types.
type Number interface {
	constraints.Integer | constraints.Float
}

// Value represents a value that tracks the maximum and minimum values over the
// last period of time. Safe for concurrent use from multiple goroutines. Do not
// copy the value. The zero value does not work, so always need to be created
// by New().
type Value[T Number] struct {
	period uint64
	tmask  uint64
	last   *rec[T]
	mu     sync.RWMutex
}

// String returns the string representation of Value.
func (v *Value[_]) String() string {
	cur, min, max := v.Get()
	return fmt.Sprintf("%v (min=%v, max=%v)", cur, min, max)
}

type rec[T Number] struct {
	start          uint64 // uint64(time.Now().UnixNano()) & tmask
	last, min, max T
	prev, next     *rec[T]
}

// New creates and returns a Value with the specified tracking period and the
// initial value. It panics if the period is 0 or negative. Also, if the period
// is less than 2^20 ns (about 1.05ms), it is rounded up to that value. Perhaps,
// it does not make sense to use such a very short period.
func New[T Number](period time.Duration, iniValue T) *Value[T] {
	const resShift = 8                       // 1/2^(n-1) .. 1/2^n = 1/128 .. 1/256
	const minPeriod = time.Duration(1) << 20 // about 1.05ms
	if period <= 0 {
		panic(fmt.Sprintf("peak.New: invalid period %v", period))
	}
	if period < minPeriod {
		period = minPeriod
	}

	// calculate resolution
	res := uint64(period) - 1
	for i := 0; i < 6; i++ { // 64 bit = 2^6
		res |= res >> (1 << i)
	}
	res++
	res >>= resShift

	// create record link
	head := &rec[T]{last: iniValue, min: iniValue, max: iniValue}
	tail := head
	for i := 0; i < 1<<resShift; i++ {
		tail.next = &rec[T]{prev: tail}
		tail = tail.next
	}
	tail.next, head.prev = head, tail

	v := &Value[T]{
		period: uint64(period),
		tmask:  ^(res - 1),
		last:   head,
	}
	v.last.start = v.tm()

	return v
}

func (v *Value[_]) tm() uint64 {
	return uint64(time.Now().UnixNano()) & v.tmask
}

// Get returns the current value and the maximum / minimum values within the
// period.
func (v *Value[T]) Get() (cur, min, max T) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	since := (uint64(time.Now().UnixNano()) - v.period) & v.tmask
	rec := v.last
	cur, min, max = rec.last, rec.last, rec.last
	for ; since <= rec.start; rec = rec.prev {
		if rec.min < min {
			min = rec.min
		}
		if max < rec.max {
			max = rec.max
		}
	}
	return
}

func (v *Value[T]) update(op func(*T, T)) T {
	v.mu.Lock()
	defer v.mu.Unlock()
	tm, rec := v.tm(), v.last
	if rec.start == tm {
		op(&rec.last, rec.last)
		switch {
		case rec.max < rec.last:
			rec.max = rec.last
		case rec.last < rec.min:
			rec.min = rec.last
		}
	} else {
		rec = rec.next
		v.last = rec
		rec.start = tm
		op(&rec.last, rec.prev.last)
		rec.min, rec.max = rec.last, rec.last
	}
	return rec.last
}

// Add adds delta to the Value v and returns the new value. It is legal to pass
// a negative delta for signed integers and floating point types. This has the
// same result as Sub. It is also possible to subtract an unsigned integer by
// adding ^(delta-1) in the same manner as for atomic.AddUint32/64.
func (v *Value[T]) Add(delta T) T {
	return v.update(func(p *T, last T) { *p = last + delta })
}

// Sub subtracts delta from the Value v and returns the new value.
func (v *Value[T]) Sub(delta T) T {
	return v.update(func(p *T, last T) { *p = last - delta })
}

// Set sets the Value to newValue.
func (v *Value[T]) Set(newValue T) T {
	return v.update(func(p *T, _ T) { *p = newValue })
}
