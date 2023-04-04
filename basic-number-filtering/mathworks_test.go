package mathworks

import (
	"reflect"
	"testing"
)

func TestEvenNumbers(t *testing.T) {
	// given
	numRange := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	want := []int{2, 4, 6, 8, 10}

	// when
	output := evenNumbers(numRange)

	// then
	if !reflect.DeepEqual(output, want) {
		t.Errorf("wanted: %v but got: %v", want, output)
	}
}

func TestOddNumbers(t *testing.T) {
	// given
	numRange := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	want := []int{1, 3, 5, 7, 9}

	// when
	output := oddNumbers(numRange)

	// then
	if !reflect.DeepEqual(output, want) {
		t.Errorf("wanted: %v but got: %v", want, output)
	}
}

func TestPrimeNumbers(t *testing.T) {
	// given
	numRange := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	want := []int{2, 3, 5, 7}

	// when
	output := primeNumbers(numRange)

	// then
	if !reflect.DeepEqual(output, want) {
		t.Errorf("wanted: %v but got: %v", want, output)
	}
}

func TestOddPrimeNumbers(t *testing.T) {
	// given
	numRange := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	want := []int{3, 5, 7}

	// when
	output := oddPrimeNumbers(numRange)

	// then
	if !reflect.DeepEqual(output, want) {
		t.Errorf("wanted: %v but got: %v", want, output)
	}
}

func TestOddPrimeNumberFilter(t *testing.T) {
	// given
	numRange := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	want := []int{3, 5, 7}

	odd := func(n int) bool { return n%2 != 0 }
	prime := func(n int) bool { return isPrime(n) }
	// when
	output := filter(numRange, odd, prime)

	// then
	if !reflect.DeepEqual(output, want) {
		t.Errorf("wanted: %v but got: %v", want, output)
	}
}

func TestEvenMultiplesOf5NumberFilter(t *testing.T) {
	// given
	numRange := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	want := []int{10, 20}

	even := func(n int) bool { return n%2 == 0 }
	multiplesOf5 := func(n int) bool { return n%5 == 0 }
	// when
	output := filter(numRange, even, multiplesOf5)

	// then
	if !reflect.DeepEqual(output, want) {
		t.Errorf("wanted: %v but got: %v", want, output)
	}
}

func TestOddMultiplesOf3GreaterThan10NumberFilter(t *testing.T) {
	// given
	numRange := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	want := []int{15}

	odd := func(n int) bool { return n%2 != 0 }
	multiplesOf := func(n int) Condition { return func(m int) bool { return m%n == 0 } }
	multiplesOf3 := multiplesOf(3)
	//multiplesOf3 := func(n int) bool { return n%3 == 0 }

	greaterThanN := func(n int) Condition { return func(m int) bool { return m > n } }
	greaterThan10 := greaterThanN(10)
	//greaterThan10 := func(n int) bool { return n > 10 }

	// when
	output := filter(numRange, odd, multiplesOf3, greaterThan10)

	// then
	if !reflect.DeepEqual(output, want) {
		t.Errorf("wanted: %v but got: %v", want, output)
	}
}

func TestMatchAllConditionsNumberFilter(t *testing.T) {
	// given
	odd := func(n int) bool { return n%2 != 0 }
	even := func(n int) bool { return !odd(n) }
	greaterThanN := func(n int) Condition { return func(m int) bool { return m > n } }
	greaterThan5 := greaterThanN(5)
	multiplesOf := func(n int) Condition { return func(m int) bool { return m%n == 0 } }
	multiplesOf3 := multiplesOf(3)
	lessThanN := func(n int) Condition { return func(m int) bool { return m < n } }
	lessThan15 := lessThanN(15)

	tt := []struct {
		nums  []int
		conds []Condition
		want  []int
	}{
		{
			nums:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			conds: []Condition{odd, greaterThan5, multiplesOf3},
			want:  []int{9, 15},
		},
		{
			nums:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			conds: []Condition{even, lessThan15, multiplesOf3},
			want:  []int{6, 12},
		},
	}

	// when
	for _, tc := range tt {
		got := filter(tc.nums, tc.conds...)
		// then
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}

func TestMatchAnyConditionsNumberFilter(t *testing.T) {
	// given
	prime := func(n int) bool { return isPrime(n) }
	greaterThanN := func(n int) Condition { return func(m int) bool { return m > n } }
	greaterThan15 := greaterThanN(15)
	multiplesOf := func(n int) Condition { return func(m int) bool { return m%n == 0 } }
	multiplesOf5 := multiplesOf(5)
	multiplesOf3 := multiplesOf(3)
	lessThanN := func(n int) Condition { return func(m int) bool { return m < n } }
	lessThan6 := lessThanN(6)

	tt := []struct {
		nums  []int
		conds []Condition
		want  []int
	}{
		{
			nums:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			conds: []Condition{prime, greaterThan15, multiplesOf5},
			want:  []int{2, 3, 5, 7, 10, 11, 13, 15, 16, 17, 18, 19, 20},
		},
		{
			nums:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			conds: []Condition{lessThan6, multiplesOf3}, // less than 6, multiple of 3
			want:  []int{1, 2, 3, 4, 5, 6, 9, 12, 15, 18},
		},
	}

	// when
	for _, tc := range tt {
		got := filterAny(tc.nums, tc.conds...)
		// then
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}
