// Copyright © The MPS Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mps

// Hashsort creates a sorted copy of the given
// slice and returns it.
func Hashsort[T hashable](s []T) []T {
	if s == nil {
		return nil
	} else if len(s) == 0 {
		return []T{}
	}
	
	sorted := make([]T, len(s))
	hm := newHashmap(s)

	var min, max T
	for i := 0; i < len(s); i++ {
		keyIdx := hm.index(uint64(s[i]))

		// increase the key's count in the hashmap
		hm.set(keyIdx, hm.lookup(keyIdx)+1)

		if s[i] < min {
			min = s[i]
		} else if s[i] > max {
			max = s[i]
		}
	}

	var idx int
	for i := min; i <= max; i++ {
		keyIdx := hm.index(uint64(i))

		if keyIdx != 0 {
			n := hm.lookup(keyIdx)

			for j := 0; j < int(n); j++ {
				sorted[idx] = i
				idx++
			}
		}
	}

	return sorted
}
