package mathworks

import (
	"math"
)

func evenNumbers(nums []int) []int {
	var output []int
	for _, n := range nums {
		if n%2 == 0 {
			output = append(output, n)
		}
	}
	return output
}

func oddNumbers(nums []int) []int {
	var output []int
	for _, n := range nums {
		if n%2 != 0 {
			output = append(output, n)
		}
	}
	return output
}

func primeNumbers(nums []int) []int {
	var output []int
	for _, n := range nums {
		if isPrime(n) {
			output = append(output, n)
		}
	}
	return output
}

func oddPrimeNumbers(nums []int) []int {
	var output []int
	for _, n := range nums {
		if isPrime(n) && n%2 != 0 {
			output = append(output, n)
		}
	}
	return output
}

func isPrime(n int) bool {
	for i := 2; i <= int(math.Floor(math.Sqrt(float64(n)))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return n > 1
}

type Condition func(n int) bool

func filter(nums []int, conditions ...Condition) []int {
	var output []int

	for _, n := range nums {
		matchesConds := true
		for _, c := range conditions {
			if !c(n) {
				matchesConds = false
			}
		}

		if matchesConds {
			output = append(output, n)
		}
	}
	return output
}

func filterAny(nums []int, conditions ...Condition) []int {
	var output []int

	for _, n := range nums {
		matchesConds := false
		for _, c := range conditions {
			if c(n) {
				matchesConds = true
			}
		}

		if matchesConds {
			output = append(output, n)
		}
	}
	return output
}
