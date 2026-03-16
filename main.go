// Copyright 2026 The Smith Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"math/rand"

	"github.com/pointlander/smith/pagerank"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// Size is the size of the universe
const Size = 16

var (
	// FlagIterations number of iterations
	FlagIterations = flag.Int("i", 8, "number of iterations")
)

func main() {
	flag.Parse()

	rng := rand.New(rand.NewSource(1))
	var u [Size][Size]byte
	rank := func() ([]float64, float64) {
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
		sum := 0.0
		for _, rank := range ranks {
			sum += rank
		}
		avg := sum / float64(len(ranks))
		v := 0.0
		for _, rank := range ranks {
			diff := rank - avg
			v += diff * diff
		}
		v /= float64(len(ranks))
		return ranks, v
	}
	indexes := rng.Perm(Size)
	for _, i := range indexes {
		count := 0
		for _, value := range u[i] {
			if value != 0 {
				count++
			}
		}
		perm := rng.Perm(Size)
		count = Size/2 - count
		for j, value := range perm {
			if value != 0 {
				continue
			}
			if count == 0 {
				break
			}
			u[i][j] = 1
			u[j][i] = 1
			count--
		}
	}
	points := make(plotter.XYs, 0, 8)
	for iterations := range *FlagIterations * 1024 {
		ranks, variance := rank()
		points = append(points, plotter.XY{X: float64(iterations), Y: variance})
	search:
		for {
			a, b := rng.Intn(Size), rng.Intn(Size)

			aa, bb := make([]int, 0, 8), make([]int, 0, 8)
			for i, value := range u[a] {
				if value != 0 {
					aa = append(aa, i)
				}
			}
			for i, value := range u[b] {
				if value != 0 {
					bb = append(bb, i)
				}
			}
			rng.Shuffle(len(aa), func(i, j int) {
				aa[i], aa[j] = aa[j], aa[i]
			})
			rng.Shuffle(len(bb), func(i, j int) {
				bb[i], bb[j] = bb[j], bb[i]
			})

			if u[a][b] == 1 && u[b][a] == 1 {
				u[a][b] = 0
				u[b][a] = 0
				r, _ := rank()
				if r[a] < ranks[a] || r[b] < ranks[b] {
					u[a][b] = 1
					u[b][a] = 1
				} else {
					break search
				}
			} else {
				if len(aa) >= 5 {
					u[a][aa[0]] = 0
					u[aa[0]][a] = 0
				}
				if len(bb) >= 5 {
					u[b][bb[0]] = 0
					u[bb[0]][b] = 0
				}
				u[a][b] = 1
				u[b][a] = 1
				r, _ := rank()
				if r[a] < ranks[a] || r[b] < ranks[b] {
					u[a][b] = 0
					u[b][a] = 0
					if len(aa) >= 5 {
						u[a][aa[0]] = 1
						u[aa[0]][a] = 1
					}
					if len(bb) >= 5 {
						u[b][bb[0]] = 1
						u[bb[0]][b] = 1
					}
				} else {
					break search
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
	}

	p := plot.New()

	p.Title.Text = "y vs x"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	scatter, err := plotter.NewScatter(points)
	if err != nil {
		panic(err)
	}
	scatter.GlyphStyle.Radius = vg.Length(1)
	scatter.GlyphStyle.Shape = draw.CircleGlyph{}
	p.Add(scatter)

	err = p.Save(8*vg.Inch, 8*vg.Inch, "plot.png")
	if err != nil {
		panic(err)
	}
}
