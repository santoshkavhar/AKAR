package main

import (
	"github.com/santoshkavhar/akar/v1/akar"
)

func main() {
	// Initial compile and run
	//	if err := compileAndRun(); err != nil {
	//		panic(err)
	//	}
	// TODO: Compile with AKAR as 1 when exiting the program.
	akar.AKAR = 1

	//fmt.Println("a")
	akar.MonitorChanges()

}
