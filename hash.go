// Copyright © The MPS Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mps

import (
	"math/bits"

	"golang.org/x/exp/constraints"
)

// hashable is the set of hashable types.
type hashable interface {
	constraints.Integer
}

// hashmap is a hash table implementing
// a minimal perfect hashing function.
type hashmap[T hashable] struct {
	hf *hashFunc // hf is the hashmap's hashing function
	vals []T // vals is the hashmap's values index
}

// newHashmap constructs a new hash table.
// It configures it's minimal perfect hash function
// to the given keys and allocates the necessary size
// for the values.
func newHashmap[T hashable](keys []T) *hashmap[T] {
	hm := new(hashmap[T])

	keysBits := make(bitvec, len(keys))

	for i := 0; i < len(keys); i++ {
		keysBits[i] = uint64(keys[i])
	}

	hm.hf = newHashFunc(keysBits)
	hm.vals = make([]T, len(keys))

	return hm
}

// index hashes the given key and returns
// its index in the hashmap.
func (hm *hashmap[T]) index(key uint64) uint64 {
	return hm.hf.query(key)
}

// lookup returns the value at the given index
// from the hashmap.
func (hm *hashmap[T]) lookup(idx uint64) T {
	return hm.vals[idx]
}

// set sets the value at the given index
// from the hashmap to the given value.
func (hm *hashmap[T]) set(idx uint64, val T) {
	hm.vals[idx] = val
}

// Much of this code is a refactored version
// of the code from https://github.com/dgryski/go-boomphf

// bitvec holds the bit representation of a hash function.
type bitvec []uint64

// newBitvec initializes and returns a new
// bitvec of the given size.
func newBitvec(size uint32) bitvec {
	return make([]uint64, uint(size+63)/64)
}

// bit returns the n'th bit in the bitvec.
func (b bitvec) bit(n uint32) uint {
	shift := n % 64
	nbit := b[n/64]
	nbit &= (1 << shift)

	return uint(nbit >> shift)
}

// setBit switches the n'th bit in the bitvec.
func (b bitvec) setBit(n uint32) {
	b[n/64] |= (1 << (n % 64))
}

// hashFunc represents a hash function
// and its hashing data.
type hashFunc struct {
	bits []bitvec
	ranks []bitvec
}

// Gamma controls the hash functions bias between
// memory footprint and construction speed.
const gamma float64 = 2

// newHash constructs a perfect hash function
// specific to the given keys.
func newHashFunc(keys bitvec) *hashFunc {
	hf := new(hashFunc)

	var level uint32
	var redo bitvec

	size := uint32(gamma * float64(len(keys)))
	size = (size + 63) &^ 63
	a := newBitvec(size)
	collide := newBitvec(size)

	for len(keys) > 0 {
		for _, k := range keys {
			hash := xorshift(k)
			h1, h2 := uint32(hash), uint32(hash>>32)
			idx := (h1 ^ rotl(h2, level)) % size

			if collide.bit(idx) == 1 {
				continue
			}

			if a.bit(idx) == 1 {
				collide.setBit(idx)
				continue
			}

			a.setBit(idx)
		}

		bv := newBitvec(size)
		for _, k := range keys {
			hash := xorshift(k)
			h1, h2 := uint32(hash), uint32(hash>>32)
			idx := (h1 ^ rotl(h2, level)) % size

			if collide.bit(idx) == 1 {
				redo = append(redo, k)
				continue
			}

			bv.setBit(idx)
		}
		hf.bits = append(hf.bits, bv)

		keys = redo
		redo = redo[:0]
		size = uint32(gamma * float64(len(keys)))
		size = (size + 63) &^ 63
		a = newBitvec(uint32(len(a)))
		collide = newBitvec(uint32(len(a)))
		level++
	}

	hf.configRanks()

	return hf
}

// query returns a key's index in a hash table.
func (hf *hashFunc) query(key uint64) uint64 {

	hash := xorshift(key)
	h1, h2 := uint32(hash), uint32(hash>>32)

	for i, bv := range hf.bits {
		idx := (h1 ^ rotl(h2, uint32(i))) % uint32(len(bv)*64)

		if bv.bit(idx) == 0 {
			continue
		}

		rank := hf.ranks[i][idx/512]

		for j := (idx / 64) &^ 7; j < idx/64; j++ {
			rank += uint64(bits.OnesCount64(bv[j]))
		}

		w := bv[idx/64]

		rank += uint64(bits.OnesCount64(w << (64 - (idx % 64))))

		return rank + 1
	}

	return 0
}

// size returns the hash function's data size in bytes.
func (hf *hashFunc) size() int {
	var size int
	for _, v := range hf.bits {
		size += len(v) * 8
	}
	for _, v := range hf.ranks {
		size += len(v) * 8
	}
	return size
}

// configRanks configures the hash function's ranks
// according to its bits.
func (hf *hashFunc) configRanks() {
	var pop uint64
	for _, bv := range hf.bits {

		r := make([]uint64, 0, 1+(len(bv)/8))

		for i, v := range bv {
			if i%8 == 0 {
				r = append(r, pop)
			}
			pop += uint64(bits.OnesCount64(v))
		}
		hf.ranks = append(hf.ranks, r)
	}
}

// rotl rotates the bits to the left.
func rotl(v uint32, r uint32) uint32 {
	return (v << r) | (v >> (32 - r))
}

// 64-bit xorshift multiply range
// from http://vigna.di.unimi.it/ftp/papers/xorshift.pdf
func xorshift(x uint64) uint64 {
	x ^= x >> 12 // a
	x ^= x << 25 // b
	x ^= x >> 27 // c
	return x * 2685821657736338717
}