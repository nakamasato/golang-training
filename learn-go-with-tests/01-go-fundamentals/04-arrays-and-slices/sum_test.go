package main

import (
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {
	assertSum := func(t *testing.T, got, want int) {
		if got != want {
			t.Errorf("got %d want %d given", got, want)
		}
	}

	t.Run(
		"Sum up 1 to 3",
		func(t *testing.T) {
			numbers := []int{1, 2, 3}

			got := Sum(numbers)
			want := 6

			assertSum(t, got, want)
		},
	)
	t.Run(
		"Sum up 1 to 5",
		func(t *testing.T) {
			numbers := []int{1, 2, 3, 4, 5}

			got := Sum(numbers)
			want := 15

			assertSum(t, got, want)
		},
	)
	t.Run("Sum two arrays", func(t *testing.T) {
		got := SumAll([]int{1, 2}, []int{0, 9})
		want := []int{3, 9}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("Sum three arrays", func(t *testing.T) {
		got := SumAll([]int{1, 2}, []int{0, 9}, []int{0, 10, 3})
		want := []int{3, 9, 13}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestSumAllTails(t *testing.T) {
	checkSums := func(t testing.TB, got, want []int) {
		t.Helper()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	}
	t.Run("make the sums of tails of", func(t *testing.T) {
		got := SumAllTails([]int{1, 2}, []int{0, 9})
		want := []int{2, 9}
		checkSums(t, got, want)
	})
	t.Run("safely sum empty slices", func(t *testing.T) {
		got := SumAllTails([]int{}, []int{3, 4, 5})
		want := []int{0, 9}
		checkSums(t, got, want)
	})
}
