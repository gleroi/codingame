package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

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

func (t Truck) Add(b Box) (Truck, bool) {
	if t.Volume+b.V <= 100 {
		t.Volume += b.V
		t.Weight += b.W
		return t, true
	}
	return t, false
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
	simple(boxes, trucks)

	// for each box, print its truck id
	for _, box := range boxes {
		fmt.Printf("%d ", box.Truck)
	}
	fmt.Println()
}

func mean(boxes []Box) (float64, float64) {
	mW := 0.0
	for _, box := range boxes {
		mW += box.W
	}
	return MaxTruckVolume, mW / MaxTruck
}

func delta(trucks []Truck) float64 {
	min, max := math.MaxFloat64, -1.0

	for _, t := range trucks {
		if t.Weight < min {
			min = t.Weight
		}
		if t.Weight > max {
			max = t.Weight
		}
	}

	return max - min
}

func simple(boxes []Box, trucks []Truck) {
	for b, box := range boxes {
		done := false
		for t, truck := range trucks {
			trucks[t], done = truck.Add(box)
			if done {
				boxes[b].Truck = t
				break
			}
		}
		if !done {
			debug("box %d: %v with not trucks\n", b, box)
		}
	}
}

func annealing(boxes []Box, trucks []Truck) {
	delta := delta(trucks)
	randMax := len(boxes)
	rng := rand.New(rand.NewSource(42))

	start := time.Now()
	for delta > 0 && time.Now().Sub(start).Seconds() < 49.988 {
		b1 := rng.Intn(randMax)
		b2 := rng.Intn(randMax)
		if b1 == b2 {
			continue
		}
		if boxes[b1].Truck == boxes[b2].Truck {
			continue
		}
		// validate swap

		// validate better delta

	}
}

func debug(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
}
