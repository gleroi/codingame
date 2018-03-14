package main

import (
	"fmt"
	"math"
	"os"
	"sort"
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

func (t Truck) Remove(b Box) Truck {
	t.Volume -= b.V
	t.Weight -= b.W
	return t
}

const MaxTruckVolume = 100
const MaxTruck = 100

func main() {
	var boxCount int
	fmt.Scan(&boxCount)

	boxes := make([]Box, boxCount)
	for i := 0; i < boxCount; i++ {
		fmt.Scan(&boxes[i].W, &boxes[i].V)
		boxes[i].Truck = 0
	}

	trucks := make([]Truck, MaxTruck)
	for t := range trucks {
		trucks[t].ID = t
	}
	firstfit(boxes, trucks)
	optimize(boxes, trucks)

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

func NewDelta() Delta {
	return Delta{
		Min:   math.MaxFloat64,
		Max:   -math.MaxFloat64,
		Delta: math.MaxFloat64,
		MinT:  -1,
		MaxT:  -1,
	}
}

func delta(trucks []Truck) Delta {
	d := NewDelta()
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

func firstfit(boxes []Box, trucks []Truck) {
	//TODO: modify to put boxes in all trucks
	bIds := make([]int, len(boxes))
	for id := range boxes {
		bIds[id] = id
	}
	sort.Slice(bIds, func(i, j int) bool {
		return boxes[i].W >= boxes[j].W
	})

	for _, b := range bIds {
		box := boxes[b]

		delta := NewDelta()
		bestT := -1
		for t, truck := range trucks {
			nt, ok := truck.Add(box)
			if ok {
				nd := delta.Update(nt)
				if nd.Delta < delta.Delta {
					bestT = t
					delta = nd
				}
			}
		}
		if bestT == -1 {
			debug("box %d not placed\n", b)
		} else {
			trucks[bestT], _ = trucks[bestT].Add(box)
			boxes[b].Truck = bestT
		}
	}
}

func split(boxes []Box, trucks []Truck, T1, T2 int) {
	// all boxes are in one truck T1
	// split boxes in two truck T1 and T2'
	// such as: weight are equal
	// until T and T' volume is less than 100

}

func findBoxes(boxes []Box, t int) []int {
	ids := make([]int, 0, len(boxes))
	for b, box := range boxes {
		if box.Truck == t {
			ids = append(ids, b)
		}
	}
	return ids
}

func optimize(boxes []Box, trucks []Truck) {
	for iter := 0; iter < 5000; iter++ {
		delta := delta(trucks)
		debug("delta: %f, min: %d, max: %d\n", delta.Delta, delta.MinT, delta.MaxT)
		debug("delta: min: %f, max: %f\n", delta.Min, delta.Max)
		boxesMin := findBoxes(boxes, delta.MinT)
		boxesMax := findBoxes(boxes, delta.MaxT)
		b1, b2 := -1, -1
		d := math.MaxFloat64
		for _, bMin := range boxesMin {
			boxMin := boxes[bMin]
			for _, bMax := range boxesMax {
				boxMax := boxes[bMax]
				deltaBoxes := math.Abs(boxMax.W - boxMin.W)
				dNext := math.Abs(delta.Delta - deltaBoxes)
				if dNext < d && canSwap(boxes, trucks, bMin, bMax) {
					b1 = bMin
					b2 = bMax
					d = dNext
				}
			}
		}
		if b1 != -1 && b2 != -1 {
			swap(boxes, trucks, b1, b2)
		} else {
			debug("remaining iter %d\n", iter)
			return
		}
	}
}

func canSwap(boxes []Box, trucks []Truck, b1Id int, b2Id int) bool {
	b1, b2 := boxes[b1Id], boxes[b2Id]
	t1, t2 := trucks[b1.Truck], trucks[b2.Truck]

	t1 = t1.Remove(b1)
	t2 = t2.Remove(b2)

	var ok bool
	t1, ok = t1.Add(b2)
	if !ok {
		return false
	}
	t2, ok = t2.Add(b1)
	return ok
}

func swap(boxes []Box, trucks []Truck, b1Id int, b2Id int) {
	b1, b2 := boxes[b1Id], boxes[b2Id]
	t1, t2 := trucks[b1.Truck], trucks[b2.Truck]

	t1 = t1.Remove(b1)
	t2 = t2.Remove(b2)

	var ok bool
	t1, ok = t1.Add(b2)
	if !ok {
		panic(fmt.Errorf("cannot swap %d from %d to %d\n", b2Id, t2.ID, t1.ID))
	}
	t2, ok = t2.Add(b1)
	if !ok {
		panic(fmt.Errorf("cannot swap %d from %d to %d\n", b1Id, t1.ID, t2.ID))
	}
	debug("swap\n  %+v\nand\n  %+v\n", b1, b2)
	trucks[t1.ID] = t1
	boxes[b2Id].Truck = t1.ID
	trucks[t2.ID] = t2
	boxes[b1Id].Truck = t2.ID
}

func fillTrucks(boxes []Box, trucks []Truck) {
	for b := 0; b < len(boxes); b++ {
		box := boxes[b]

		for t, truck := range trucks {
			nt, ok := truck.Add(box)
			if ok {
				trucks[t] = nt
				boxes[b].Truck = nt.ID
				break
			}
		}
	}
}

const DEBUG = false

func debug(format string, v ...interface{}) {
	if DEBUG {
		fmt.Fprintf(os.Stderr, format, v...)
	}
}
