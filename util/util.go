package util

import "math/rand"

func AbsInt(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func Shuffle(numbers []int) []int {
	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})
	return numbers
}
