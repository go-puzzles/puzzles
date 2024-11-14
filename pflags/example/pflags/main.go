package main

import (
	"fmt"
	
	"github.com/go-puzzles/puzzles/pflags"
)

var (
	intFlag = pflags.Int("int", 0, "an int flag")
)

func main() {
	pflags.Parse()
	
	fmt.Println(intFlag())
}
