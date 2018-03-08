package main

import (
	"fmt"
	"math"
	"os"
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
	ID     int
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
		boxes[i].Truck = -1
	}

	trucks := make([]Truck, MaxTruck)
	for t := range trucks {
		trucks[t].ID = t
	}
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

type Delta struct {
	Min, Max, Delta float64
	MinT, MaxT      int
}

func delta(trucks []Truck) Delta {
	d := Delta{
		Min:   math.MaxFloat64,
		Max:   -math.MaxFloat64,
		Delta: math.MaxFloat64,
		MinT:  -1,
		MaxT:  -1,
	}
	if trucks == nil {
		return d
	}
	return d.Update(trucks...)
}

func (d Delta) Update(trucks ...Truck) Delta {
	for _, t := range trucks {
		if d.MinT == t.ID {
			d.Min = t.Weight
		}
		if d.MaxT == t.ID {
			d.Max = t.Weight
		}
		if t.Weight < d.Min {
			d.Min = t.Weight
			d.MinT = t.ID
		}
		if t.Weight > d.Max {
			d.Max = t.Weight
			d.MaxT = t.ID
		}
	}
	d.Delta = d.Max - d.Min
	return d
}

func simple(boxes []Box, trucks []Truck) {
	unplacedBox := true

	iter := 0
	for unplacedBox {
		unplacedBox = false
		unplacedCount := 0

		for b := 0; b < len(boxes); b++ {
			box := boxes[b]
			if box.Truck != -1 {
				continue
			}
			minT := -1
			deltaMin := delta(trucks)

			for t, truck := range trucks {
				if truck.Weight == 0 {
					minT = t
					break
				}
				nt, ok := truck.Add(box)
				if ok {
					nd := deltaMin.Update(nt)
					if nd.Delta < deltaMin.Delta {
						deltaMin = nd
						minT = t
					}
				}
			}
			if minT != -1 {
				trucks[minT], _ = trucks[minT].Add(box)
				boxes[b].Truck = minT
			} else {
				unplacedBox = true
				unplacedCount++
			}
		}
		iter++
		debug("iter %d: %d remaining box\n", iter, unplacedCount)
	}
}

func debug(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
}
