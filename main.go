// Copyright 2026 The Smith Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math/rand"

	"github.com/pointlander/smith/pagerank"
)

// Size is the size of the universe
const Size = 8

func main() {
	rng := rand.New(rand.NewSource(1))
	var u [Size][Size]byte
	rank := func() []float64 {
		g := pagerank.NewGraph(Size, rng)
		for i := range u {
			for j := range u {
				if j > i {
					break
				}
				if u[i][j] == 1 && u[j][i] == 1 {
					g.Link(uint32(i), uint32(j), 1)
					g.Link(uint32(j), uint32(i), 1)
				}
			}
		}
		ranks := make([]float64, Size)
		g.Rank(.85, 0.0000001, func(node int, rank float64) {
			ranks[node] = rank
		})
		return ranks
	}
	for i := range u {
		for j := range u {
			if j > i {
				break
			}
			linked := rng.Intn(2)
			if linked == 0 {
				u[i][j] = 1
				u[j][i] = 1
			}
		}
	}
	for {
		hd := true
	search:
		for a := range Size {
			for b := range Size {
				current := rank()
				if u[a][b] == 1 && u[b][a] == 1 {
					u[a][b] = 0
					u[b][a] = 0
					next := rank()
					if current[a] > next[a] || current[b] > next[b] {
						u[a][b] = 1
						u[b][a] = 1
					} else {
						hd = false
						break search
					}
				} else {
					u[a][b] = 1
					u[b][a] = 1
					next := rank()
					if current[a] > next[a] || current[b] > next[b] {
						u[a][b] = 0
						u[b][a] = 0
					} else {
						hd = false
						break search
					}
				}
			}
		}
		for i := range u {
			for j := range u {
				if u[i][j] != u[j][i] {
					panic("not symmetric")
				}
			}
		}
		for i := range u {
			fmt.Println(u[i])
		}
		fmt.Println()
		if hd {
			break
		}
	}
}
