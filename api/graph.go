package main

import (
	"errors"
)

type Node[T any] struct {
	Id    string `json:"id"`
	Value T      `json:"value"`
}

type Edge[T any] struct {
	From *Node[T] `json:"from"`
	To   *Node[T] `json:"to"`
	Code int      `json:"code"`
}

type Graph[T any] struct {
	Nodes        []*Node[T] `json:"nodes"`
	Edges        []*Edge[T] `json:"edges"`
	isNodeUnique func(g *Graph[T], value T) bool
}

func NewGraph[T any](
	isNodeUnique func(g *Graph[T], value T) bool,
) *Graph[T] {
	return &Graph[T]{
		Nodes:        []*Node[T]{},
		Edges:        []*Edge[T]{},
		isNodeUnique: isNodeUnique,
	}
}

func (g *Graph[T]) AddNode(value T) (*Node[T], error) {
	if !g.isNodeUnique(g, value) {
		return nil, errors.New("node is not unique")
	}
	newNode := &Node[T]{Value: value}
	g.Nodes = append(g.Nodes, newNode)
	return newNode, nil
}

func (g *Graph[T]) AddEdge(from, to *Node[T]) (*Edge[T], error) {
	newEdge := &Edge[T]{
		From: from,
		To:   to,
	}
	g.Edges = append(g.Edges, newEdge)
	return newEdge, nil
}

func (g *Graph[T]) FindNodeByValue(
	findFunc func(n *Node[T]) bool,
) (*Node[T], error) {
	for _, node := range g.Nodes {
		if found := findFunc(node); found == true {
			return node, nil
		}
	}
	return nil, errors.New("value not found")
}
