// Copyright 2018 The SAT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/pointlander/go-galib"
)

const (
	// ProblemSize is the size of the problems
	ProblemSize = 100
	// RestartAfter is the number of generations after which to Restarts
	RestartAfter = 10
)

// Optimize runs ga on a give problem
func Optimize(size int, problem *Node, statistics *Statistics) {
	scores, restarts := 0, 0
	score := func(g *ga.GAFixedBitstringGenome) float64 {
		scores++
		vars := make([]int32, len(g.Gene))
		for i, bit := range g.Gene {
			if bit {
				vars[i] = One
			}
		}
		return float64(One-problem.Eval(vars)) / One
	}

	m := ga.NewMultiMutator()
	m.Add(new(ga.GAShiftMutator))
	m.Add(new(ga.GASwitchMutator))
	m.Add(new(ga.GAMutatorRandom))
	param := ga.GAParameter{
		Initializer: new(ga.GARandomInitializer),
		Selector:    ga.NewGATournamentSelector(0.7, 5),
		PMutate:     0.2,
		PBreed:      0.2,
		Breeder:     new(ga.GA2PointBreeder),
		Mutator:     m,
	}

	found := false
	for !found {
		gao := ga.NewGA(param)
		genome := ga.NewFixedBitstringGenome(make([]bool, size), score)
		gao.Init(128, genome)

		stuck, count := 0.0, 0
		gao.OptimizeUntil(func(best ga.GAGenome) bool {
			score := best.Score()
			if *verbose {
				fmt.Printf("best = %v\n", score)
			}
			if score == stuck {
				count++
			} else {
				stuck = score
				count = 0
			}
			if count > RestartAfter {
				return true
			}

			found = score < 1e-3
			return found
		})
		if !found {
			restarts++
		}
	}

	statistics.Scores.Add(float64(scores))
	statistics.Restarts.Add(float64(restarts))
}

// OptimizeMeta runs ga on a give problem
func OptimizeMeta(size int, problem *Node, statistics *Statistics) {
	scores, restarts := 0, 0
	score := func(g *ga.GAFixedBitstringGenome) float64 {
		scores++
		vars := make([]int32, len(g.Gene))
		for i, bit := range g.Gene {
			if bit {
				vars[i] = One
			}
		}
		return float64(One-problem.Eval(vars)) / One
	}

	m := ga.NewMultiMutator()
	m.Add(new(ga.GAShiftMutator))
	m.Add(new(ga.GASwitchMutator))
	m.Add(new(ga.GAMutatorRandom))
	param := ga.GAParameter{
		Initializer: new(ga.GARandomInitializer),
		Selector:    ga.NewGATournamentSelector(0.7, 5),
		PMutate:     0.2,
		PBreed:      0.2,
		Breeder:     new(ga.GA2PointBreeder),
		Mutator:     m,
	}

	found := false
	for !found {
		gao := ga.NewGA(param)
		genome := ga.NewFixedBitstringGenome(make([]bool, size), score)
		gao.Init(128, genome)

		stuck, count := 0.0, 0
		gao.OptimizeUntil(func(best ga.GAGenome) bool {
			score := best.Score()
			if *verbose {
				fmt.Printf("best = %v\n", score)
			}
			if score == stuck {
				count++
			} else {
				stuck = score
				count = 0
			}
			if count > RestartAfter {
				count = 0
				g := best.(*ga.GAFixedBitstringGenome)
				a := &Node{
					NodeType: Or,
					Nodes:    make([]*Node, len(g.Gene)),
				}
				for i, bit := range g.Gene {
					if bit {
						b := &Node{
							NodeType: Not,
							Nodes:    make([]*Node, 1),
						}
						b.Nodes[0] = &Node{
							NodeType: Var,
							A:        i,
						}
						a.Nodes[i] = b
					} else {
						b := &Node{
							NodeType: Var,
							A:        i,
						}
						a.Nodes[i] = b
					}
				}
				problem.Nodes = append(problem.Nodes, a)
				gao.Reset()
				restarts++
			}

			found = score < 1e-3
			return found
		})
		if !found {
			restarts++
		}
	}

	statistics.Scores.Add(float64(scores))
	statistics.Restarts.Add(float64(restarts))
}

var verbose = flag.Bool("v", false, "verbose mode")

func main() {
	flag.Parse()

	rand.Seed(time.Now().Unix())

	run := func(name string, newProblem func() (int, *Node),
		optimize func(size int, problem *Node, statistics *Statistics)) {
		start := time.Now()
		statistics := Statistics{}
		for i := 0; i < 10; i++ {
			size, problem := newProblem()
			optimize(size, problem, &statistics)
		}
		fmt.Println(name)
		fmt.Println(statistics.String())
		fmt.Println(time.Now().Sub(start))
		fmt.Println("")
	}

	twoPeak := func() (int, *Node) {
		return TwoPeak(ProblemSize)
	}
	run("TwoPeak normal", twoPeak, Optimize)

	falsePeak := func() (int, *Node) {
		return FalsePeak(ProblemSize)
	}
	run("FalsePeak normal", falsePeak, Optimize)

	hamiltonian := func() (int, *Node) {
		return HamiltonianCircuit(8)
	}
	run("HamiltonianCircuit normal", hamiltonian, Optimize)
	run("HamiltonianCircuit meta", hamiltonian, OptimizeMeta)
}
