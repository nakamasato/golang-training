package main

func Sum(numbers []int) int {
	var total int
	for _, n := range numbers {
		total += n
	}
	return total
}

func SumAll(AllNumbers... []int) []int {
	var SummedNumbers []int
	for _, numbers := range AllNumbers {
		SummedNumbers = append(SummedNumbers, Sum(numbers))
	}
	return SummedNumbers
}
