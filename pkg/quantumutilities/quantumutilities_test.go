package quantumutilities_test

import (
	"quantumwerewolf/pkg/quantumutilities"
	"testing"
)

func TestFactorial(t *testing.T) {
	// Test up to 20 because 21! overflows uint64
	tables := []struct {
		n int
		f uint64
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 6},
		{4, 24},
		{5, 120},
		{6, 720},
		{7, 5040},
		{8, 40320},
		{9, 362880},
		{10, 3628800},
		{11, 39916800},
		{12, 479001600},
		{13, 6227020800},
		{14, 87178291200},
		{15, 1307674368000},
		{16, 20922789888000},
		{17, 355687428096000},
		{18, 6402373705728000},
		{19, 121645100408832000},
		{20, 2432902008176640000},
	}

	for _, table := range tables {
		result := quantumutilities.Factorial(table.n)
		if result != table.f {
			t.Errorf("%d! was incorrect, got: %d, want: %d.", table.n, result, table.f)
		}
	}
}

func copySlice(input []int) []int {
	output := make([]int, len(input), len(input))
	for i, v := range input {
		output[i] = v
	}
	return output
}

func TestKthperm(t *testing.T) {
	original := make([]int, 0, 3)
	original = append(original, 0, 1, 2)
	input1 := copySlice(original)
	input2 := copySlice(original)
	output2 := make([]int, 0, 3)
	output2 = append(output2, 0, 2, 1)

	tables := []struct {
		s []int
		k uint64
		p []int
	}{
		{input1, 0, original},
		{input2, 1, output2},
	}

	for _, table := range tables {
		result := quantumutilities.Kthperm(table.s, table.k)
		for i, v := range result {
			if v != table.p[i] {
				t.Errorf("k=%d was incorrect, got: %v, want: %v.", table.k, result, table.p)
				break
			}
		}
	}
}
