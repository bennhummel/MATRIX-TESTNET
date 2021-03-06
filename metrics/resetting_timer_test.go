// Copyright 2018 The MATRIX Authors as well as Copyright 2014-2017 The go-ethereum Authors
// This file is consisted of the MATRIX library and part of the go-ethereum library.
//
// The MATRIX-ethereum library is free software: you can redistribute it and/or modify it under the terms of the MIT License.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, 
//and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject tothe following conditions:
//
//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, 
//WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISINGFROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
//OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package metrics

import (
	"testing"
	"time"
)

func TestResettingTimer(t *testing.T) {
	tests := []struct {
		values   []int64
		start    int
		end      int
		wantP50  int64
		wantP95  int64
		wantP99  int64
		wantMean float64
		wantMin  int64
		wantMax  int64
	}{
		{
			values:  []int64{},
			start:   1,
			end:     11,
			wantP50: 5, wantP95: 10, wantP99: 10,
			wantMin: 1, wantMax: 10, wantMean: 5.5,
		},
		{
			values:  []int64{},
			start:   1,
			end:     101,
			wantP50: 50, wantP95: 95, wantP99: 99,
			wantMin: 1, wantMax: 100, wantMean: 50.5,
		},
		{
			values:  []int64{1},
			start:   0,
			end:     0,
			wantP50: 1, wantP95: 1, wantP99: 1,
			wantMin: 1, wantMax: 1, wantMean: 1,
		},
		{
			values:  []int64{0},
			start:   0,
			end:     0,
			wantP50: 0, wantP95: 0, wantP99: 0,
			wantMin: 0, wantMax: 0, wantMean: 0,
		},
		{
			values:  []int64{},
			start:   0,
			end:     0,
			wantP50: 0, wantP95: 0, wantP99: 0,
			wantMin: 0, wantMax: 0, wantMean: 0,
		},
		{
			values:  []int64{1, 10},
			start:   0,
			end:     0,
			wantP50: 1, wantP95: 10, wantP99: 10,
			wantMin: 1, wantMax: 10, wantMean: 5.5,
		},
	}
	for ind, tt := range tests {
		timer := NewResettingTimer()

		for i := tt.start; i < tt.end; i++ {
			tt.values = append(tt.values, int64(i))
		}

		for _, v := range tt.values {
			timer.Update(time.Duration(v))
		}

		snap := timer.Snapshot()

		ps := snap.Percentiles([]float64{50, 95, 99})

		val := snap.Values()

		if len(val) > 0 {
			if tt.wantMin != val[0] {
				t.Fatalf("%d: min: got %d, want %d", ind, val[0], tt.wantMin)
			}

			if tt.wantMax != val[len(val)-1] {
				t.Fatalf("%d: max: got %d, want %d", ind, val[len(val)-1], tt.wantMax)
			}
		}

		if tt.wantMean != snap.Mean() {
			t.Fatalf("%d: mean: got %.2f, want %.2f", ind, snap.Mean(), tt.wantMean)
		}

		if tt.wantP50 != ps[0] {
			t.Fatalf("%d: p50: got %d, want %d", ind, ps[0], tt.wantP50)
		}

		if tt.wantP95 != ps[1] {
			t.Fatalf("%d: p95: got %d, want %d", ind, ps[1], tt.wantP95)
		}

		if tt.wantP99 != ps[2] {
			t.Fatalf("%d: p99: got %d, want %d", ind, ps[2], tt.wantP99)
		}
	}
}
