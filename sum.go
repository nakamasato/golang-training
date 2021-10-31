package main


func Sum(numbers []int) int {
	sum := 0
	for _, number := range numbers {
        sum += number
    }
	return sum
}

func SumAll(numbersToSum ...[]int) (sums []int) {
	for _, numbers := range numbersToSum {
		sums = append(sums, Sum(numbers))
	}
    return
}

func SumAllTails(numbersToSum ...[]int) (sums [] int) {
	for _, numbers := range numbersToSum {
		if len(numbers) == 0 {
			sums = append(sums, 0)
		} else {
			sums = append(sums, Sum(numbers[1:]))
		}
	}
	return
}
