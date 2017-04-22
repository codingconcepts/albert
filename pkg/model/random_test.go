package model

import (
	"fmt"
	"math"
	"testing"
	"testing/quick"

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

func TestInboxName(t *testing.T) {
	testCases := []struct {
		length         int
		expectedLength int
	}{
		{length: 10, expectedLength: InboxMaxLength},
		{length: 11, expectedLength: 11 * 2},
		{length: 19, expectedLength: 19 * 2},
		{length: -1, expectedLength: InboxMaxLength * 2},
		{length: 0, expectedLength: InboxMaxLength * 2},
		{length: 1, expectedLength: InboxMaxLength * 2},
		{length: 9, expectedLength: InboxMaxLength * 2},
		{length: InboxMaxLength, expectedLength: InboxMaxLength * 2},
		{length: 21, expectedLength: InboxMaxLength * 2},
		{length: 101, expectedLength: InboxMaxLength * 2},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("TestInboxName_%d", testCase.length), func(t *testing.T) {
			inboxName, err := InboxName(testCase.length)
			t.Log(inboxName)

			test.ErrorNil(t, err)
			test.Equals(t, testCase.expectedLength, len(inboxName))
		})
	}
}

func TestBlackBoxCheckBetween(t *testing.T) {
	f := func(min int, max int) bool {
		result := Between(min, max)

		// Between ensures min is less than or equal to max,
		// so perform the same switch here
		if min > max {
			min, max = max, min
		}
		return result >= min && result <= max
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestBlackBoxTakeRandom(t *testing.T) {
	f := func(input []string, perc float64) bool {
		// TakeRandom ensures perc is between 0 and 1,
		// so perform the same check here
		if perc < 0 || perc > 1 {
			return true
		}

		result := TakeRandom(input, perc)

		actPerc := math.Trunc(float64((100 * len(result)) / len(input)))
		expPerc := math.Trunc(perc)

		return actPerc == expPerc
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestBlackBoxInboxName(t *testing.T) {
	f := func(length int) bool {
		result, err := InboxName(length)
		return err == nil && len(result) >= InboxMinLength && len(result) <= InboxMaxLength*2
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func BenchmarkBetween(b *testing.B) {
	testCases := []struct {
		min int
		max int
	}{
		{min: 1, max: 10},
		{min: 1, max: 100},
		{min: 1, max: 10000},
		{min: 1, max: 1000000},
	}

	for _, testCase := range testCases {
		b.Run(fmt.Sprintf("BenchmarkBetween_%d_%d", testCase.min, testCase.max), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Between(testCase.min, testCase.max)
			}
		})
	}
}

func BenchmarkTakeRandom(b *testing.B) {
	testCases := []struct {
		input []string
		perc  float64
	}{
		{input: []string{"a", "b", "c", "d"}, perc: 0.25},
		{input: []string{"a", "b", "c", "d"}, perc: 0.50},
		{input: []string{"a", "b", "c", "d"}, perc: 0.75},
		{input: []string{"a", "b", "c", "d"}, perc: 1},
	}

	for _, testCase := range testCases {
		b.Run(fmt.Sprintf("BenchmarkTakeRandom_%s_%.2f", testCase.input, testCase.perc), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				TakeRandom(testCase.input, testCase.perc)
			}
		})
	}
}
