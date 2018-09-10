package random

import (
	"math/rand"
)

// Roll generates a random number between 1 and the number of sides for each die
func Roll(dice, sides int) []int {
	var results []int
	for i := 0; i < dice; i++ {
		roll := rand.Intn(sides) + 1
		results = append(results, roll)
	}
	return results
}
