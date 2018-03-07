package main

import "fmt"

//import "os"

/*

Constraints

  100 ≤ boxCount ≤ 3300
  0 < weight < 100
  0 < volume ≤ 26
  Response time: 50 seconds

*/

type Box struct {
	W, V  float64
	Truck int
}

// Truck carries boxes, with a maximal volume of 100.
type Truck struct {
	Volume float64
	Weight float64
}

const MaxTruckVolume = 100
const MaxTruck = 100

func main() {
	var boxCount int
	fmt.Scan(&boxCount)

	boxes := make([]Box, boxCount)
	for i := 0; i < boxCount; i++ {
		fmt.Scan(&boxes[i].W, &boxes[i].V)
	}

	trucks := make([]Truck, MaxTruck)

	// for each box, print its truck id
	fmt.Println("0 0 0 0 0 ...") // Write action to stdout
}
