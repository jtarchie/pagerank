/*
Package pagerank implements the *weighted* PageRank algorithm.
*/
package pagerank

import (
	"math"
)

type node struct {
	weight   float64
	outbound float64
}

// Graph holds node and edge data.
type Graph[T comparable] struct {
	edges map[T](map[T]float64)
	nodes map[T]*node
}

// NewGraph initializes and returns a new graph.
func NewGraph[T comparable]() *Graph[T] {
	return &Graph[T]{
		edges: make(map[T](map[T]float64)),
		nodes: make(map[T]*node),
	}
}

// Link creates a weighted edge between a source-target node pair.
// If the edge already exists, the weight is incremented.
func (self *Graph[T]) Link(source, target T, weight float64) {
	if _, ok := self.nodes[source]; !ok {
		self.nodes[source] = &node{
			weight:   0,
			outbound: 0,
		}
	}

	self.nodes[source].outbound += weight

	if _, ok := self.nodes[target]; !ok {
		self.nodes[target] = &node{
			weight:   0,
			outbound: 0,
		}
	}

	if _, ok := self.edges[source]; !ok {
		self.edges[source] = map[T]float64{}
	}

	self.edges[source][target] += weight
}

// Rank computes the PageRank of every node in the directed graph.
// α (alpha) is the damping factor, usually set to 0.85.
// ε (epsilon) is the convergence criteria, usually set to a tiny value.
//
// This method will run as many iterations as needed, until the graph converges.
func (self *Graph[T]) Rank(α, ε float64, callback func(id T, rank float64)) {
	Δ := float64(1.0)
	inverse := 1 / float64(len(self.nodes))

	// Normalize all the edge weights so that their sum amounts to 1.
	for source := range self.edges {
		if self.nodes[source].outbound > 0 {
			for target := range self.edges[source] {
				self.edges[source][target] /= self.nodes[source].outbound
			}
		}
	}

	for key := range self.nodes {
		self.nodes[key].weight = inverse
	}

	for Δ > ε {
		leak := float64(0)
		nodes := map[T]float64{}

		for key, value := range self.nodes {
			nodes[key] = value.weight

			if value.outbound == 0 {
				leak += value.weight
			}

			self.nodes[key].weight = 0
		}

		leak *= α

		for source := range self.nodes {
			for target, weight := range self.edges[source] {
				self.nodes[target].weight += α * nodes[source] * weight
			}

			self.nodes[source].weight += (1-α)*inverse + leak*inverse
		}

		Δ = 0

		for key, value := range self.nodes {
			Δ += math.Abs(value.weight - nodes[key])
		}
	}

	for key, value := range self.nodes {
		callback(key, value.weight)
	}
}

// Reset clears all the current graph data.
func (self *Graph[T]) Reset() {
	self.edges = make(map[T](map[T]float64))
	self.nodes = make(map[T]*node)
}
