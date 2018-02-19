package main

import (
	"fmt"
	"io"
	"math"
	"os"
)

//import "os"

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func distance(a, b Vec) float64 {
	return math.Sqrt(float64((b.X-a.X)*(b.X-a.X) + (b.Y-a.Y)*(b.Y-a.Y)))
}

type Vec struct {
	X, Y float64
}

func Vec2f(x, y float64) Vec {
	return Vec{X: x, Y: y}
}

func (v Vec) Scale(s float64) Vec {
	return Vec2f(v.X*s, v.Y*s)
}

func (v Vec) Round() Vec {
	return Vec2f(math.Round(v.X), math.Round(v.Y))
}

func add(a, b Vec) Vec {
	return Vec2f(a.X+b.X, a.Y+b.Y)
}

const GravityVal = 3.711

var Gravity = Vec2f(0, -GravityVal)

type surface struct {
	points []Vec
}

func readSurface(r io.Reader, debug bool) *surface {
	var surface surface
	var surfaceN int
	fmt.Fscan(r, &surfaceN)

	surface.points = make([]Vec, surfaceN)
	for i := 0; i < surfaceN; i++ {
		fmt.Fscan(r, &surface.points[i].X, &surface.points[i].Y)
	}
	if debug {
		fmt.Fprintln(os.Stderr, surfaceN)
		for _, pt := range surface.points {
			fmt.Fprintln(os.Stderr, pt.X, pt.Y)
		}
	}
	return &surface
}

func (s *surface) LandingZone() (Vec, Vec) {
	var a, b Vec
	//There is a unique area of flat ground on the surface of Mars,
	// which is at least 1000 meters wide.

	a = s.points[0]
	for _, pt := range s.points {
		if pt.Y == a.Y {
			b = pt
		} else {
			// if landing zone is big enough returns it
			// else restart at pt
			if distance(a, b) >= 1000 {
				return a, b
			}
			a = pt
			b = pt
		}
	}
	if distance(a, b) >= 1000 {
		return a, b
	}
	panic("landing zone not found!")
}

type Force int

type Angle int

type lander struct {
	Pos      Vec
	Speed    Vec
	Fuel     int
	Rotation Angle // +/- 15
	Power    Force // +/- 1
}

func radian(a Angle) float64 {
	r := float64(a) * 180 / math.Pi
	return math.Pi/2 + r
}

func (l lander) Next(timeStep float64, power Force, rotation Angle) lander {
	n := l

	rad := radian(rotation)
	accelaration := add(Gravity, Vec2f(math.Cos(rad), math.Sin(rad)).Scale(float64(power)))
	n.Speed = add(l.Speed, accelaration)
	n.Pos = add(l.Pos, add(l.Speed, accelaration.Scale(0.5*timeStep)))

	n.Power = power
	n.Rotation = rotation
	n.Fuel = l.Fuel - int(power)
	return n
}

func (l lander) Round() lander {
	l.Speed = l.Speed.Round()
	l.Pos = l.Pos.Round()
	return l
}

func readLander(r io.Reader, debug bool) lander {
	var lander lander
	fmt.Fscan(r, &lander.Pos.X, &lander.Pos.Y, &lander.Speed.X, &lander.Speed.Y, &lander.Fuel, &lander.Rotation, &lander.Power)
	if debug {
		fmt.Fprintln(os.Stderr, lander.Pos.X, lander.Pos.Y, lander.Speed.X, lander.Speed.Y, lander.Fuel, lander.Rotation, lander.Power)
	}
	return lander
}

func main() {
	/*
			For a landing to be successful, the ship must:

		    land on flat ground
		    land in a vertical position (tilt angle = 0°)
		    vertical speed must be limited ( ≤ 40m/s in absolute value)
		    horizontal speed must be limited ( ≤ 20m/s in absolute value)
	*/
	_ = readSurface(os.Stdin, true)

	for {
		// hSpeed: the horizontal speed (in m/s), can be negative.
		// vSpeed: the vertical speed (in m/s), can be negative.
		// fuel: the quantity of remaining fuel in liters.
		// rotate: the rotation angle in degrees (-90 to 90).
		// power: the thrust power (0 to 4).
		lander := readLander(os.Stdin, true)

		const minVSpeed = -35
		const maxPower = 4
		const minPower = 0
		power := int(lander.Power)
		// fmt.Fprintln(os.Stderr, "Debug messages...")
		if lander.Speed.Y < minVSpeed {
			power = min(power+1, maxPower)

		} else {
			power = max(power-1, minPower)
		}

		// 2 integers: rotate power. rotate is the desired rotation angle (should be 0 for level 1), power is the desired thrust power (0 to 4).
		fmt.Println("0", power)
	}
}
