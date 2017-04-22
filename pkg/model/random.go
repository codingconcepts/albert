package model

import (
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// TakeRandom takes a random percentage of strings from a slice
// of strings, guaranteeing that each item will appear only once.
func TakeRandom(input []string, perc float64) (output []string) {
	// if no percentage has been specified or percentage is invalid,
	// don't kill anything
	if perc == 0 || perc < 0 || perc > 1 {
		return
	}

	// if there is only one application, or perc not defined, return all
	if len(input) == 1 {
		return input
	}

	amount := int(math.Trunc((float64(len(input)) / 1) * perc))
	if amount < 1 {
		return
	}

	output = make([]string, amount)

	for i := 0; i < amount; i++ {
		index := Between(i, len(input))

		input[index], input[i] = input[i], input[index]
		output[i] = input[i]
	}

	return
}

// Between returns a number between min and max inclusively.
func Between(min int, max int) int {
	// swap min and max if max is less than min
	if max < min {
		min, max = max, min
	}

	return rand.Int()%(max-min) + min
}
