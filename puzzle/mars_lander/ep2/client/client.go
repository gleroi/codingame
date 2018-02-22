package client

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

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func distance(a, b Vec) float64 {
	return math.Sqrt((b.X-a.X)*(b.X-a.X) + (b.Y-a.Y)*(b.Y-a.Y))
}

func add(a, b Vec) Vec {
	return Vec2f(a.X+b.X, a.Y+b.Y)
}

func sub(a, b Vec) Vec {
	return Vec2f(a.X-b.X, a.Y-b.Y)
}

func norm(a Vec) float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y)
}

func normalize(a Vec) Vec {
	norm := norm(a)
	return Vec2f(a.X/norm, a.Y/norm)
}

func neg(a Vec) Vec {
	return Vec2f(-a.X, -a.Y)
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

func round(f float64) int {
	if f < -0.5 {
		return int(f - 0.5)
	}
	if f > 0.5 {
		return int(f + 0.5)
	}
	return 0
}

func (v Vec) Round() Vec {
	return Vec2f(float64(round(v.X)), float64(round(v.Y)))
}

const GravityVal = 3.711

var Gravity = Vec2f(0, -GravityVal)

type Surface struct {
	Points []Vec
}

func ReadSurface(r io.Reader, debug bool) *Surface {
	var surface Surface
	var surfaceN int
	fmt.Fscan(r, &surfaceN)

	surface.Points = make([]Vec, surfaceN)
	for i := 0; i < surfaceN; i++ {
		fmt.Fscan(r, &surface.Points[i].X, &surface.Points[i].Y)
	}
	if debug {
		fmt.Fprintln(os.Stderr, surfaceN)
		for _, pt := range surface.Points {
			fmt.Fprintln(os.Stderr, pt.X, pt.Y)
		}
	}
	return &surface
}

type Zone struct {
	Start, End Vec
}

func (z *Zone) Middle() Vec {
	return Vec2f(z.Start.X+(z.End.X-z.Start.X)/2, z.Start.Y)
}

func (s *Surface) LandingZone() Zone {
	var a, b Vec
	//There is a unique area of flat ground on the surface of Mars,
	// which is at least 1000 meters wide.

	a = s.Points[0]
	for _, pt := range s.Points {
		if pt.Y == a.Y {
			b = pt
		} else {
			// if landing zone is big enough returns it
			// else restart at pt
			if distance(a, b) >= 1000 {
				return Zone{a, b}
			}
			a = pt
			b = pt
		}
	}
	if distance(a, b) >= 1000 {
		return Zone{a, b}
	}
	panic("landing zone not found!")
}

type Lander struct {
	Pos      Vec
	Speed    Vec
	Fuel     int
	Rotation float64 // +/- 15
	Power    float64 // +/- 1
}

func Radian(a float64) float64 {
	r := a * math.Pi / 180
	return math.Pi/2 + r
}

func degree(r float64) float64 {
	return r * 180 / math.Pi
}

func (l Lander) Next(timeStep float64, power float64, rotation float64) Lander {
	n := l

	rad := Radian(rotation)
	accelaration := add(Gravity, Vec2f(math.Cos(rad), math.Sin(rad)).Scale(power))
	n.Speed = add(l.Speed, accelaration)
	n.Pos = add(l.Pos, add(l.Speed, accelaration.Scale(0.5*timeStep)))

	n.Power = power
	n.Rotation = rotation
	n.Fuel = l.Fuel - round(power)
	return n
}

func (l Lander) Round() Lander {
	l.Speed = l.Speed.Round()
	l.Pos = l.Pos.Round()
	return l
}

func ReadLander(r io.Reader, debug bool) Lander {
	var lander Lander
	fmt.Fscan(r, &lander.Pos.X, &lander.Pos.Y, &lander.Speed.X, &lander.Speed.Y, &lander.Fuel, &lander.Rotation, &lander.Power)
	if debug {
		fmt.Fprintln(os.Stderr, lander.Pos.X, lander.Pos.Y, lander.Speed.X, lander.Speed.Y, lander.Fuel, lander.Rotation, lander.Power)
	}
	return lander
}

func Client(r io.Reader, w io.Writer) {
	/*
			For a landing to be successful, the ship must:

		    land on flat ground
		    land in a vertical position (tilt angle = 0°)
		    vertical speed must be limited ( ≤ 40m/s in absolute value)
		    horizontal speed must be limited ( ≤ 20m/s in absolute value)
	*/
	surface := ReadSurface(r, true)
	lander := ReadLander(r, true)

	for {
		// hSpeed: the horizontal speed (in m/s), can be negative.
		// vSpeed: the vertical speed (in m/s), can be negative.
		// fuel: the quantity of remaining fuel in liters.
		// rotate: the rotation angle in degrees (-90 to 90).
		// power: the thrust power (0 to 4).

		const minVSpeed = 35.0
		const minHSpeed = 15.0
		const maxPower = 4
		const minPower = 0
		debug("running...\n")

		zone := surface.LandingZone()
		remainingSeconds := float64(lander.Fuel)
		if lander.Power > 0 {
			remainingSeconds /= lander.Power
		}

		next := lander.Next(1, lander.Power, lander.Rotation)
		targetPos := add(lander.Pos, sub(zone.Middle(), lander.Pos).Scale(1/remainingSeconds))
		debug("target pos: %v\n", targetPos)
		targetDirection := sub(targetPos, next.Pos)
		debug("target dir: %v\n", targetDirection)

		targetSpeed := targetDirection

		targetAccel := sub(targetSpeed, lander.Speed)
		targetThrust := sub(targetAccel, Gravity)

		debug("target speed: %v\n", targetSpeed)
		debug("target accel: %v\n", targetAccel)
		debug("target thrus: %v\n", targetThrust)

		power := targetThrust.Y
		debug("target power: %v\n", power)
		power = min(power, maxPower)
		power = max(power, minPower)

		targetAngle := normalize(targetAccel)
		rad := math.Acos(targetAngle.X)

		// minAngle := -25
		// maxAngle := 25
		angle := degree(rad)
		debug("target angle: %v %.3f\n", targetAngle, angle-90)
		angle = angle - 90

		next = lander.Next(1, power, 0)
		if next.Pos.Y <= zone.Start.Y {
			angle = 0
		}

		// 2 integers: rotate power. rotate is the desired rotation angle (should be 0 for level 1), power is the desired thrust power (0 to 4).
		debug("write lander data\n")
		fmt.Fprintln(w, round(angle), round(power))
		debug("read lander data\n")
		next = ReadLander(r, true)
		lander = lander.Next(1, power, angle)
	}
}

func debug(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "client: "+format, v...)
}
