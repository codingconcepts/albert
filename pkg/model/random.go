package model

import (
	prng "crypto/rand"
	"encoding/hex"
	"math"
	"math/rand"
	"time"
)

const (
	// InboxMinLength is the minimum size (in bytes) of an inbox
	InboxMinLength = 10

	// InboxMaxLength is the maximum size (in bytes) of an inbox
	InboxMaxLength = 20
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

// InboxName generates a cryptographically-secure random string
// which can be used as an inbox for orchestrators and agents
// during communication.
func InboxName(length int) (name string, err error) {
	// if for any reason a ridiculous length has been provided,
	// swap it out for a sensible default.
	if length < InboxMinLength || length > InboxMaxLength {
		length = InboxMaxLength
	}

	data := make([]byte, length)
	if _, err = prng.Read(data); err != nil {
		return
	}

	return hex.EncodeToString(data), nil
}
