// File:		base.go
// Created by:	Hoven
// Created on:	2025-01-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package basepuzzle

import "github.com/go-puzzles/puzzles/cores"

type BasePuzzle struct {
	PuzzleName string
}

func (p *BasePuzzle) Name() string {
	return p.PuzzleName
}

func (p *BasePuzzle) Before(_ *cores.Options) error {
	return nil
}
