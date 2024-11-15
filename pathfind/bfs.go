// File:		main.go
// Created by:	Hoven
// Created on:	2024-11-15
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pathfind

import (
	"errors"
	
	"github.com/go-puzzles/puzzles/pqueue"
)

func BFSSearch(graph Graph, start, goal Point) ([]Point, error) {
	if !graph.IsInGraph(start) || !graph.IsInGraph(goal) {
		return nil, errors.New("start point or goal point is not in graph")
	}
	
	if start.Equals(goal) {
		return []Point{start}, nil
	}
	
	parent := make(map[Point]Point)
	totalWeight := make(map[Point]int)
	
	queue := pqueue.NewMemoryQueue[Point]()
	queue.Enqueue(start)
	graph.SetVisited(start)
	totalWeight[start] = 0
	
	for {
		isEmpty, _ := queue.IsEmpty()
		if isEmpty {
			break
		}
		
		point, _ := queue.Dequeue()
		
		if point.Equals(goal) {
			path := []Point{}
			current := point
			for current != nil {
				path = append([]Point{current}, path...)
				current = parent[current]
			}
			return path, nil
		}
		
		for _, n := range graph.Neighbors(point) {
			if graph.IsBlocked(n) || graph.IsVisited(n) {
				continue
			}
			
			newWeight := totalWeight[point] + graph.GetEdgeWeight(point, n)
			if _, exists := totalWeight[n]; !exists || newWeight > totalWeight[n] {
				parent[n] = point
				totalWeight[n] = newWeight
				graph.SetVisited(n)
				queue.Enqueue(n)
			}
		}
	}
	
	return nil, nil
}
