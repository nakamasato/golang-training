package main

func Sum(numbers []int) int {
	var total int
	for _, n := range numbers {
		total += n
	}
	return total
}

func SumAll(AllNumbers ...[]int) []int {
	var SummedNumbers []int
	for _, numbers := range AllNumbers {
		SummedNumbers = append(SummedNumbers, Sum(numbers))
	}
	return SummedNumbers
}

func SumAllTails(AllNumbers ...[]int) []int {
	var SummedNumbers []int
	for _, numbers := range AllNumbers {
		if len(numbers) == 0 {
			SummedNumbers = append(SummedNumbers, 0)
		} else {
			SummedNumbers = append(SummedNumbers, Sum(numbers[1:]))
		}
	}
	return SummedNumbers
}
