package main

import "github.com/santoshkavhar/akar/v1/akar"

func main() {
	// Initial compile and run
//	if err := compileAndRun(); err != nil {
//		panic(err)
//	}
	akar.AKAR = 0

	akar.MonitorChanges()

}
