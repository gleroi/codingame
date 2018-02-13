package main

import (
	"fmt"
	"math"
)

//import "os"

type Pt struct {
	X, Y int
}

func (p Pt) Add(d Pt) Pt {
	return Pt{X: p.X + d.X, Y: p.Y + d.Y}
}

type Direction string

const (
	N  Direction = "N"
	NE Direction = "NE"
	E  Direction = "E"
	SE Direction = "SE"
	S  Direction = "S"
	SW Direction = "SW"
	W  Direction = "W"
	NW Direction = "NW"
)

var directions = map[Direction]Pt{
	N:  {0, -1},
	NE: {1, -1},
	E:  {1, 0},
	SE: {1, 1},
	S:  {0, 1},
	SW: {-1, 1},
	W:  {-1, 0},
	NW: {-1, -1},
}

func distance(a, b Pt) float64 {
	return math.Sqrt(math.Pow(float64(b.X-a.X), 2) + math.Pow(float64(b.Y-a.Y), 2))
}

func possible(p Pt) bool {
	return p.X >= 0 && p.X < 40 && p.Y >= 0 && p.Y < 18
}

func main() {
	// lightX: the X position of the light of power
	// lightY: the Y position of the light of power
	// initialTX: Thor's starting X position
	// initialTY: Thor's starting Y position
	var thor Pt
	var light Pt
	fmt.Scan(&light.X, &light.Y, &thor.X, &thor.Y)

	for {
		// remainingTurns: The remaining amount of turns Thor can move. Do not remove this line.
		var remainingTurns int
		fmt.Scan(&remainingTurns)

		// fmt.Fprintln(os.Stderr, "Debug messages...")

		var direction Direction
		d := distance(thor, light)
		for dir, delta := range directions {
			next := thor.Add(delta)
			if !possible(next) {
				continue
			}
			nd := distance(next, light)
			if nd < d {
				d = nd
				direction = dir
			}
		}

		// A single line providing the move to be made: N NE E SE S SW W or NW
		fmt.Println(direction)
		thor = thor.Add(directions[direction])
	}
}
