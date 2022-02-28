// Copyright © The MPS Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mps_test

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/vorduin/mps"
	"github.com/vorduin/slices"
)

func TestHashsort(t *testing.T) {
	s := []int{2, 1, -1, 0}
	sorted := mps.Hashsort(s)

	if !slices.Equal(sorted, []int{-1, 0, 1, 2}) {
		t.Fail()
	}

	s = []int{}
	sorted = mps.Hashsort(s)

	if sorted == nil || len(sorted) != 0 {
		t.Fail()
	}

	s = nil
	sorted = mps.Hashsort(s)

	if sorted != nil {
		t.Fail()
	}
}

func BenchmarkHashsort(b *testing.B) {
	benchmarkSort(b, 1e1)
	
	benchmarkSort(b, 1e2)

	benchmarkSort(b, 1e3)

	benchmarkSort(b, 1e4)

	benchmarkSort(b, 1e5)

	benchmarkSort(b, 1e6)

	benchmarkSort(b, 1e7)
}

func benchmarkMicro(b *testing.B, f func()) {
	b.ResetTimer()

	start := time.Now()
	for i := 0; i < b.N; i++ {
		f()
	}
	execTime := time.Since(start)

	b.ReportMetric(0, "ns/op")
	b.ReportMetric((1e6*execTime.Seconds())/float64(b.N), "μs/op")
}

func benchmarkMilli(b *testing.B, f func()) {
	b.ResetTimer()

	start := time.Now()
	for i := 0; i < b.N; i++ {
		f()
	}
	execTime := time.Since(start)

	b.ReportMetric(0, "ns/op")
	b.ReportMetric((1e3*execTime.Seconds())/float64(b.N), "ms/op")
}

func benchmarkSort(b *testing.B, size int) {
	b.Run(fmt.Sprintf("Hashsort%.0eInt", float64(size)), func(b *testing.B) {
		s := rand.Perm(size)
		
		if size <= 1e4 {
			benchmarkMicro(b, func() {
				mps.Hashsort(s)
			})
		} else {
			benchmarkMilli(b, func() {
				mps.Hashsort(s)
			})
		}
	})

	b.Run(fmt.Sprintf("sort.Ints%.0eInt", float64(size)), func(b *testing.B) {
		s := rand.Perm(size)

		if size <= 1e4 {
			benchmarkMicro(b, func() {
				sort.Ints(s)
			})
		} else {
			benchmarkMilli(b, func() {
				sort.Ints(s)
			})
		}
	})
}