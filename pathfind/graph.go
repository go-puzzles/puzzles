// File:		graph.go
// Created by:	Hoven
// Created on:	2024-11-15
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pathfind

type Graph interface {
	Neighbors(p Point) []Point
	IsInGraph(p Point) bool
	IsBlocked(p Point) bool
	SetBlock(p Point)
	IsVisited(p Point) bool
	SetVisited(p Point)
	GetEdgeWeight(from, to Point) int
}

var (
	directions = []SimplePoint{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
	}
)

type SimpleGraph struct {
	graph   [][]int
	weights map[SimplePoint]map[SimplePoint]int
}

func NewSimpleGraph(width, height int) *SimpleGraph {
	graph := make([][]int, height)
	for i := range graph {
		graph[i] = make([]int, width)
	}
	weights := make(map[SimplePoint]map[SimplePoint]int)
	return &SimpleGraph{graph: graph, weights: weights}
}

func (g *SimpleGraph) Neighbors(p Point) []Point {
	var neighbors []Point
	for _, dir := range directions {
		newPoint := p.Shift(dir.X, dir.Y)
		if !g.IsInGraph(newPoint) {
			continue
		}
		neighbors = append(neighbors, newPoint)
	}
	return neighbors
}

func (g *SimpleGraph) IsInGraph(p Point) bool {
	return p.GetX() >= 0 && p.GetX() < len(g.graph) && p.GetY() >= 0 && p.GetY() < len(g.graph[0])
}

func (g *SimpleGraph) IsBlocked(p Point) bool {
	return g.graph[p.GetX()][p.GetY()] == -1
}

func (g *SimpleGraph) SetBlock(p Point) {
	g.graph[p.GetX()][p.GetY()] = -1
}

func (g *SimpleGraph) SetVisited(p Point) {
	g.graph[p.GetX()][p.GetY()] = 1
}

func (g *SimpleGraph) IsVisited(p Point) bool {
	return g.graph[p.GetX()][p.GetY()] == 1
}

func (g *SimpleGraph) GetEdgeWeight(from, to Point) int {
	return 1
}
