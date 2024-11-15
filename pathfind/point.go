// File:		point.go
// Created by:	Hoven
// Created on:	2024-11-15
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pathfind

import "fmt"

type Point interface {
	fmt.Stringer
	Equals(other Point) bool
	Shift(x, y int) Point
	GetX() int
	GetY() int
}

type SimplePoint struct {
	X, Y int
}

func (p *SimplePoint) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func (p *SimplePoint) Equals(other Point) bool {
	if o, ok := other.(*SimplePoint); ok {
		return p.X == o.X && p.Y == o.Y
	}
	return false
}

func (p *SimplePoint) Shift(x, y int) Point {
	return &SimplePoint{p.X + x, p.Y + y}
}

func (p *SimplePoint) GetX() int {
	return p.X
}

func (p *SimplePoint) GetY() int {
	return p.Y
}
