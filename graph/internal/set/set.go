// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package set

import "github.com/jak9708/gonummat/graph"

type Int interface{ ~int | ~int64 }

type Ints[T Int] map[T]struct{}

// Add inserts an element into the set.
func (s Ints[T]) Add(e T) {
	s[e] = struct{}{}
}

// Has reports the existence of the element in the set.
func (s Ints[T]) Has(e T) bool {
	_, ok := s[e]
	return ok
}

// Remove deletes the specified element from the set.
func (s Ints[T]) Remove(e T) {
	delete(s, e)
}

// Count reports the number of elements stored in the set.
func (s Ints[T]) Count() int {
	return len(s)
}

// IntsEqual reports set equality between the parameters. Sets are equal if
// and only if they have the same elements.
func IntsEqual[T Int](a, b Ints[T]) bool {
	if intsSame(a, b) {
		return true
	}

	if len(a) != len(b) {
		return false
	}

	for e := range a {
		if _, ok := b[e]; !ok {
			return false
		}
	}

	return true
}

// Nodes is a set of nodes keyed in their integer identifiers.
type Nodes map[int64]graph.Node

// NewNodes returns a new Nodes.
func NewNodes() Nodes {
	return make(Nodes)
}

// NewNodesSize returns a new Nodes with the given size hint, n.
func NewNodesSize(n int) Nodes {
	return make(Nodes, n)
}

// The simple accessor methods for Nodes are provided to allow ease of
// implementation change should the need arise.

// Add inserts an element into the set.
func (s Nodes) Add(n graph.Node) {
	s[n.ID()] = n
}

// Remove deletes the specified element from the set.
func (s Nodes) Remove(e graph.Node) {
	delete(s, e.ID())
}

// Count returns the number of element in the set.
func (s Nodes) Count() int {
	return len(s)
}

// Has reports the existence of the elements in the set.
func (s Nodes) Has(n graph.Node) bool {
	_, ok := s[n.ID()]
	return ok
}

// CloneNodes returns a clone of src.
func CloneNodes(src Nodes) Nodes {
	dst := make(Nodes, len(src))
	for e, n := range src {
		dst[e] = n
	}
	return dst
}

// Equal reports set equality between the parameters. Sets are equal if
// and only if they have the same elements.
func Equal(a, b Nodes) bool {
	if same(a, b) {
		return true
	}

	if len(a) != len(b) {
		return false
	}

	for e := range a {
		if _, ok := b[e]; !ok {
			return false
		}
	}

	return true
}

// UnionOfNodes returns the union of a and b.
//
// The union of two sets, a and b, is the set containing all the
// elements of each, for instance:
//
//	{a,b,c} UNION {d,e,f} = {a,b,c,d,e,f}
//
// Since sets may not have repetition, unions of two sets that overlap
// do not contain repeat elements, that is:
//
//	{a,b,c} UNION {b,c,d} = {a,b,c,d}
func UnionOfNodes(a, b Nodes) Nodes {
	if same(a, b) {
		return CloneNodes(a)
	}

	dst := make(Nodes)
	for e, n := range a {
		dst[e] = n
	}
	for e, n := range b {
		dst[e] = n
	}

	return dst
}

// IntersectionOfNodes returns the intersection of a and b.
//
// The intersection of two sets, a and b, is the set containing all
// the elements shared between the two sets, for instance:
//
//	{a,b,c} INTERSECT {b,c,d} = {b,c}
//
// The intersection between a set and itself is itself, and thus
// effectively a copy operation:
//
//	{a,b,c} INTERSECT {a,b,c} = {a,b,c}
//
// The intersection between two sets that share no elements is the empty
// set:
//
//	{a,b,c} INTERSECT {d,e,f} = {}
func IntersectionOfNodes(a, b Nodes) Nodes {
	if same(a, b) {
		return CloneNodes(a)
	}
	dst := make(Nodes)
	if len(a) > len(b) {
		a, b = b, a
	}
	for e, n := range a {
		if _, ok := b[e]; ok {
			dst[e] = n
		}
	}
	return dst
}
