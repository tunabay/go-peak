// Copyright (c) 2022 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package peak_test

import (
	"testing"
	"time"

	"github.com/tunabay/go-peak"
)

func veri(t *testing.T, v *peak.Value[int32], last, min, max int32) {
	t.Helper()

	la, mi, ma := v.Get()
	if la != last {
		t.Errorf("unexpected last value: got %v, want %v.", la, last)
	}
	if mi != min {
		t.Errorf("unexpected min value: got %v, want %v.", mi, min)
	}
	if ma != max {
		t.Errorf("unexpected max value: got %v, want %v.", ma, max)
	}
}

func TestValue_Add(t *testing.T) {
	v := peak.New[int32](time.Second, 100)

	v.Add(10)
	v.Add(10)
	v.Add(10)
	v.Add(10)
	v.Add(15)
	veri(t, v, 155, 100, 155)
	v.Add(-15)
	veri(t, v, 140, 100, 155)
	v.Add(-10)
	veri(t, v, 130, 100, 155)
	v.Add(-10)
	veri(t, v, 120, 100, 155)
	v.Add(-10)
	veri(t, v, 110, 100, 155)
	v.Add(-10)
	veri(t, v, 100, 100, 155)
	v.Add(-10)
	veri(t, v, 90, 90, 155)
	v.Add(-10)
	veri(t, v, 80, 80, 155)
	v.Add(15)
	veri(t, v, 95, 80, 155)
}

func TestValue_Get_sparse(t *testing.T) {
	v := peak.New[int32](time.Second, 100)
	v.Set(50)
	time.Sleep(time.Millisecond * 900)
	veri(t, v, 50, 50, 100)
	v.Set(500)
	veri(t, v, 500, 50, 500)
	time.Sleep(time.Millisecond * 200)
	veri(t, v, 500, 500, 500)
}

func TestValue_Get_1(t *testing.T) {
	//                                        [00:00.00]  cur,  min,  max
	v := peak.New[int32](time.Second, 100) // [00:00.00]  100,  100,  100
	v.Add(10)                              // [00:00.00]  110,  100,  110
	v.Add(10)                              // [00:00.00]  120,  100,  120
	v.Add(10)                              // [00:00.00]  130,  100,  130
	time.Sleep(time.Millisecond * 200)     // [00:00.20]  130,  100,  130
	v.Add(10)                              // [00:00.20]  140,  100,  140
	time.Sleep(time.Millisecond * 200)     // [00:00.40]  140,  100,  140
	v.Add(15)                              // [00:00.40]  155,  100,  155
	v.Add(-15)                             // [00:00.40]  140,  100,  155
	veri(t, v, 140, 100, 155)              //
	time.Sleep(time.Millisecond * 200)     // [00:00.60]  140,  100,  155
	v.Sub(10)                              // [00:00.60]  130,  100,  155
	veri(t, v, 130, 100, 155)              //
	time.Sleep(time.Millisecond * 200)     // [00:00.80]  130,  100,  155
	v.Sub(10)                              // [00:00.80]  120,  100,  155
	veri(t, v, 120, 100, 155)              //
	time.Sleep(time.Millisecond * 190)     // [00:00.99]  120,  100,  155
	v.Sub(10)                              // [00:00.99]  110,  100,  155
	veri(t, v, 110, 100, 155)              //
	v.Sub(10)                              // [00:00.99]  100,  100,  155
	v.Sub(10)                              // [00:00.99]   90,   90,  155
	veri(t, v, 90, 90, 155)                //
	time.Sleep(time.Millisecond * 200)     // [00:01.19]   90,   90,  155, since 00:00.19
	v.Sub(50)                              // [00:01.19]   40,   40,  155
	v.Add(12)                              // [00:01.19]   52,   40,  155
	veri(t, v, 52, 40, 155)                //
	time.Sleep(time.Millisecond * 200)     // [00:01.39]   52,   40,  155, since 00:00.39
	veri(t, v, 52, 40, 155)                //
	v.Add(10)                              // [00:01.39]   62,   40,  155
	veri(t, v, 62, 40, 155)                //
	time.Sleep(time.Millisecond * 200)     // [00:01.59]   62,   40,  130, since 00:00.59
	veri(t, v, 62, 40, 130)                //
	v.Add(10)                              // [00:01.59]   72,   40,  130
	veri(t, v, 72, 40, 130)                //
	time.Sleep(time.Millisecond * 200)     // [00:01.79]   72,   40,  120, since 00:00.79
	veri(t, v, 72, 40, 120)                //
	v.Add(10)                              // [00:01.79]   82,   40,  120
	veri(t, v, 82, 40, 120)                //
	time.Sleep(time.Millisecond * 190)     // [00:01.98]   82,   40,  100, since 00:00.98
	veri(t, v, 82, 40, 110)                //
	time.Sleep(time.Millisecond * 200)     // [00:02.18]   82,   40,   82, since 00:01.18
	veri(t, v, 82, 40, 82)                 //
	time.Sleep(time.Millisecond * 200)     // [00:02.38]   82,   62,   82, since 00:01.38
	veri(t, v, 82, 62, 82)                 //
}

func TestValue_Get_2(t *testing.T) {
	//                                        [00:00.00]  cur,  min,  max
	v := peak.New[int32](time.Second, 100) // [00:00.00]  100,  100,  100
	v.Set(10)                              // [00:00.00]   10,   10,  100
	v.Set(1000)                            // [00:00.00] 1000,   10, 1000
	v.Add(20)                              // [00:00.00] 1020,   10, 1020
	v.Sub(120)                             // [00:00.00]  900,   10, 1020
	veri(t, v, 900, 10, 1020)              //
	time.Sleep(time.Millisecond * 200)     // [00:00.20]  900,   10, 1020
	v.Set(999)                             // [00:00.20]  999,   10, 1020
	v.Set(899)                             // [00:00.20]  899,   10, 1020
	veri(t, v, 899, 10, 1020)              //
	time.Sleep(time.Millisecond * 200)     // [00:00.40]  899,   10, 1020
	veri(t, v, 899, 10, 1020)              //
	v.Set(500)                             // [00:00.40]  500,   10, 1020
	veri(t, v, 500, 10, 1020)              //
	v.Set(-500)                            // [00:00.40] -500, -500, 1020
	time.Sleep(time.Millisecond * 200)     // [00:00.60] -500, -500, 1020
	veri(t, v, -500, -500, 1020)           //
	time.Sleep(time.Millisecond * 200)     // [00:00.80]
	veri(t, v, -500, -500, 1020)           //
	time.Sleep(time.Millisecond * 190)     // [00:00.99] -500, -500, 1020
	veri(t, v, -500, -500, 1020)           //
	time.Sleep(time.Millisecond * 200)     // [00:01.19] -500, -500, 999
	veri(t, v, -500, -500, 999)            //
	time.Sleep(time.Millisecond * 400)     // [00:01.59] -500, -500, -500
	veri(t, v, -500, -500, -500)           //
	time.Sleep(time.Millisecond * 410)     // [00:02.00] -500, -500, -500
	veri(t, v, -500, -500, -500)
}
