package main

import (
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {
	assertSum := func (t *testing.T, got, want int)  {
		if got != want {
        	t.Errorf("got %d want %d given", got, want)
    	}
	}

	t.Run(
		"Sum up 1 to 3",
		func (t *testing.T)  {
			numbers := []int{1, 2, 3}

			got := Sum(numbers)
			want := 6

			assertSum(t, got, want)
		},
	)
	t.Run(
		"Sum up 1 to 5",
		func (t *testing.T)  {
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
