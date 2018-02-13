package main

import "fmt"

/**
 * The while loop represents the game.
 * Each iteration represents a turn of the game
 * where you are given inputs (the heights of the mountains)
 * and where you have to print an output (the index of the mountain to fire on)
 * The inputs you are given are automatically updated according to your last actions.
 **/

func main() {
	for {
		max, imax := -1, -1
		for i := 0; i < 8; i++ {
			var h int
			fmt.Scan(&h)
			if h > max {
				max = h
				imax = i
			}
		}

		fmt.Println(imax) // The index of the mountain to fire on.
	}
}
