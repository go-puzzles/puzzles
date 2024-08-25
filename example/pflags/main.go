package main

import (
	"fmt"
	
	"github.com/go-puzzles/puzzles/pflags"
)

type Config struct {
	Name string
}

var (
	stringTest = pflags.String("stringTest", "string", "string pflags test")
	boolTest   = pflags.Bool("boolTest", false, "bool pflags test")
	floatTest  = pflags.Float64("floatTest", 0, "float test")
	intTest    = pflags.Int("intTest", 0, "int test")
	sliceTest  = pflags.StringSlice("sliceTest", []string{"this", "is", "test"}, "slice test")
	structTest = pflags.Struct("structTest", (*Config)(nil), "struct test")
	
	stringRTest = pflags.StringRequired("stringRequired", "string required test")
	stringPTest = pflags.StringP("stringP", "p", "stringP", "string shordhand test")
)

func main() {
	pflags.Parse()
	
	fmt.Println(stringTest())
	fmt.Println(boolTest())
	fmt.Println(floatTest())
	fmt.Println(intTest())
	fmt.Println(sliceTest())
	
	conf := new(Config)
	if err := structTest(conf); err != nil {
		panic(err)
	}
	fmt.Println(conf)
	
	fmt.Println(stringRTest())
}
