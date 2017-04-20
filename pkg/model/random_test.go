package model

import (
	"fmt"
	"testing"

	"github.com/codingconcepts/albert/test"
)

func TestTakeRandom(t *testing.T) {
	applications := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	testCases := []struct {
		perc              float64
		expectedSubLength int
	}{
		{perc: 0.1, expectedSubLength: 1},
		{perc: 0.2, expectedSubLength: 2},
		{perc: 0.3, expectedSubLength: 3},
		{perc: 0.4, expectedSubLength: 4},
		{perc: 0.5, expectedSubLength: 5},
		{perc: 0.6, expectedSubLength: 6},
		{perc: 0.7, expectedSubLength: 7},
		{perc: 0.8, expectedSubLength: 8},
		{perc: 0.9, expectedSubLength: 9},
		{perc: 1, expectedSubLength: 10},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("TestApplicationsTakeRandom_%.1f", testCase.perc), func(t *testing.T) {
			sub := TakeRandom(applications, testCase.perc)

			test.Equals(t, testCase.expectedSubLength, len(sub))

			// assert uniqueness
			itemMap := make(map[string]int)
			for _, item := range sub {
				itemMap[item]++
			}

			for _, value := range itemMap {
				test.Equals(t, 1, value)
			}
		})
	}
}

func TestTakeRandomReturnOneWhenOnlyOneProvided(t *testing.T) {
	applications := []string{"a"}

	testCases := []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("TestApplicationsTakeRandom_%.1f", testCase), func(t *testing.T) {
			sub := TakeRandom(applications, testCase)

			test.Equals(t, 1, len(sub))
		})
	}
}

func TestTakeRandomTakesNoneWhenPercentageIsZero(t *testing.T) {
	applications := []string{"a", "b", "c"}
	sub := TakeRandom(applications, 0)

	test.Equals(t, 0, len(sub))
}

func TestTakeRandomTakesNoneWhenPercentageIsTiny(t *testing.T) {
	applications := []string{"a", "b", "c"}
	sub := TakeRandom(applications, 0.01)

	test.Equals(t, 0, len(sub))
}
