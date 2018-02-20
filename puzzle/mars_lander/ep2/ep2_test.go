package main

import (
	"strings"
	"testing"
)

var examples = []struct {
	input        string
	landingStart Vec
	landingEnd   Vec
}{
	{
		input: `7
		0 100
		1000 500
		1500 1500
		3000 1000
		4000 150
		5500 150
		6999 800
		2500 2700 0 0 550 0 0`,
		landingStart: Vec2f(4000, 150),
		landingEnd:   Vec2f(5500, 150),
	},
	{
		input: `10
		0 100
		1000 500
		1500 100
		3000 100
		3500 500
		3700 200
		5000 1500
		5800 300
		6000 1000
		6999 2000
		6500 2800 -100 0 600 90 0`,
		landingStart: Vec2f(1500, 100),
		landingEnd:   Vec2f(3000, 100),
	},
}

func TestIdentifyLandingZone(t *testing.T) {
	for testID, test := range examples {
		r := strings.NewReader(test.input)
		surface := readSurface(r, false)

		zone := surface.LandingZone()
		if test.landingStart != zone.Start || test.landingEnd != zone.End {
			t.Errorf("case %d: expected a landing zone %v -> %v, got %v -> %v (%fm)", testID,
				test.landingStart, test.landingEnd, zone.Start, zone.End, distance(zone.Start, zone.End))
		}
	}
}

var landerTests = []struct {
	steps []string
}{
	{
		steps: []string{
			"2500 2700 0 0 550 0 0",
			"2500 2698 0 -4 550 0 0",
			"2500 2693 0 -7 550 0 0",
			"2500 2683 0 -11 550 0 0",
			"2500 2670 0 -15 550 0 0",
			"2500 2654 0 -19 550 0 0",
			"2500 2633 0 -22 550 0 0",
			"2500 2609 0 -26 550 0 0",
			"2500 2581 0 -30 550 0 0",
			"2500 2550 0 -33 550 0 0",
			"2500 2514 0 -37 550 0 0",
		},
	},
	{
		steps: []string{
			"2500 2700 0 0 550 0 0",
			"2500 2698 0 -4 550 0 0",
			"2500 2693 0 -7 550 0 0",
			"2500 2683 0 -11 550 0 0",
			"2500 2670 0 -15 550 0 0",
			"2500 2654 0 -19 550 0 0",
			"2500 2633 0 -22 550 0 0",
			"2500 2609 0 -26 550 0 0",
			"2500 2581 0 -30 550 0 0",
			"2500 2550 0 -33 550 0 0",
			"2500 2514 0 -37 550 0 0",
			"2500 2476 0 -40 549 0 1",
			"2500 2435 0 -42 547 0 2",
			"2500 2393 0 -42 544 0 3",
			"2500 2351 0 -42 540 0 4",
			"2500 2310 0 -42 536 0 4",
			"2500 2268 0 -41 532 0 4",
			"2500 2227 0 -41 528 0 4",
			"2500 2186 0 -41 524 0 4",
			"2500 2145 0 -41 520 0 4",
			"2500 2105 0 -40 516 0 4",
			"2500 2065 0 -40 512 0 4",
			"2500 2025 0 -40 508 0 4",
			"2500 1985 0 -39 504 0 4",
			"2500 1946 0 -39 500 0 4",
			"2500 1907 0 -39 496 0 4",
			"2500 1869 0 -38 492 0 4",
			"2500 1830 0 -38 488 0 4",
			"2500 1792 0 -38 484 0 4",
			"2500 1755 0 -38 480 0 4",
			"2500 1717 0 -37 476 0 4",
			"2500 1680 0 -37 472 0 4",
			"2500 1643 0 -37 468 0 4",
			"2500 1606 0 -36 464 0 4",
			"2500 1570 0 -36 460 0 4",
			"2500 1534 0 -36 456 0 4",
			"2500 1498 0 -36 452 0 4",
			"2500 1463 0 -35 448 0 4",
			"2500 1427 0 -36 445 0 3",
			"2500 1391 0 -36 441 0 4",
			"2500 1356 0 -35 437 0 4",
			"2500 1320 0 -36 434 0 3",
			"2500 1284 0 -36 430 0 4",
			"2500 1248 0 -36 426 0 4",
			"2500 1213 0 -35 422 0 4",
			"2500 1177 0 -36 419 0 3",
		},
	},
}

func TestLanderNextPosition(t *testing.T) {
	for testID, test := range landerTests {
		current := readLander(strings.NewReader(test.steps[0]), false)
		for i := 1; i < len(test.steps); i++ {
			expected := readLander(strings.NewReader(test.steps[i]), false)

			next := current.Next(1, expected.Power, expected.Rotation)

			/*
				The simulation rounds its result when transmitting it, but seems to use
				the not rounded result to compute next step.
				Observed for: Speed
			*/

			if next.Round() != expected {
				t.Errorf("case %d: step %d -> %d:\nexpected:\n  %+v\ngot:\n  %+v", testID,
					i, i+1, expected, next)
			}
			current = next
		}
	}
}
