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

func TestFactorialNegativeInput(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Factorial should panic when given a negative input.")
		}
	}()

	quantumutilities.Factorial(-1)
	quantumutilities.Factorial(-10)
	quantumutilities.Factorial(-21354)
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
	output2 := make([]int, 0, 3)
	output5 := make([]int, 0, 3)
	original = append(original, 0, 1, 2)
	input1 := copySlice(original)
	input2 := copySlice(original)
	input5 := copySlice(original)
	output1 := copySlice(original)
	output2 = append(output2, 0, 2, 1)
	output5 = append(output5, 2, 1, 0)

	tables := []struct {
		s []int
		k uint64
		p []int
	}{
		{input1, 0, output1},
		{input2, 1, output2},
		{input5, 5, output5},
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

func TestGetBytesGetInterface(t *testing.T) {
	var blob []byte
	var err error
	expectedInt := 53
	actualInt := expectedInt
	if expectedInt != actualInt {
		t.Errorf("GetBytes for integer did not set up correctly. expectedInt %d != actualInt %d", expectedInt, actualInt)
	}
	blob, err = quantumutilities.GetBytes(actualInt)
	if err != nil {
		t.Errorf("GetBytes for integer resulted in an error: %v.", err)
	}
	actualInt++
	if expectedInt == actualInt {
		t.Errorf("TestGetBytesGetInterface for integer did not modify correctly. expectedInt %d == actualInt %d", expectedInt, actualInt)
	}
	err = quantumutilities.GetInterface(blob, &actualInt)
	if err != nil {
		t.Errorf("GetInterface for integer resulted in an error: %v.", err)
	}
	if expectedInt != actualInt {
		t.Errorf("GetInterface for integer did not get correct result. expectedInt %d != actualInt %d", expectedInt, actualInt)
	}

	expectedMap := make(map[string]int)
	expectedMap["testValueA"] = 1
	expectedMap["testValueB"] = 2
	expectedMap["testValueC"] = 3
	blob, err = quantumutilities.GetBytes(expectedMap)
	if err != nil {
		t.Errorf("GetBytes for map resulted in an error: %v.", err)
	}
	actualMap := make(map[string]int)
	for k, v := range expectedMap {
		if v == actualMap[k] {
			t.Errorf("TestGetBytesGetInterface for map did not modify correctly. expectedInt[%v] %d == actualInt[%v] %d", k, v, k, actualMap[k])
		}
	}
	err = quantumutilities.GetInterface(blob, &actualMap)
	if err != nil {
		t.Errorf("GetInterface for map resulted in an error: %v.", err)
	}
	for k, v := range expectedMap {
		if v != actualMap[k] {
			t.Errorf("GetInterface for map did not get correct result. expectedInt[%v] %d != actualInt[%v] %d", k, v, k, actualMap[k])
		}
	}
}
