package main

import (
	"fmt"
	"log"
	"math"

	_ "github.com/lib/pq"
)

// Universe is a single unit in multiverse
type Universe struct {
	permutation uint64
	active      bool
}

var multiverse []Universe
var originalAssignments []int

func (u Universe) String() string {
	var universeString string

	// Add active display
	if u.active {
		universeString = "A"
	} else {
		universeString = "I"
	}

	// Add roles display
	universeAssignments := make([]int, len(originalAssignments))
	copy(universeAssignments, originalAssignments)
	universeAssignments = kthperm(universeAssignments, u.permutation)

	universeString += "["
	for i, v := range universeAssignments {
		if i > 0 {
			universeString += " "
		}
		universeString += roleTypes[v].Name
	}
	universeString += "]"
	return fmt.Sprint(universeString)
}

func factorial(n int) uint64 {
	factVal := uint64(1)
	if n < 0 {
		fmt.Print("Factorial of negative number doesn't exist.")
	} else {
		for i := 1; i <= n; i++ {
			factVal *= uint64(i) // mismatched types int64 and int
		}

	}
	return factVal
}

func kthperm(S []int, k uint64) []int {
	var P []int
	for len(S) > 0 {
		f := factorial(len(S) - 1)
		i := int(math.Floor(float64(k) / float64(f)))
		x := S[i]
		k = k % f
		P = append(P, x)
		S = append(S[:i], S[i+1:]...)
	}

	return P
}

func createMultiverse() {
	setupRoles()
	for i := 0; i < roles.Villagers; i++ {
		originalAssignments = append(originalAssignments, villager.ID)
	}
	for i := 0; i < roles.Seers; i++ {
		originalAssignments = append(originalAssignments, seer.ID)
	}
	for i := 0; i < roles.Wolves; i++ {
		originalAssignments = append(originalAssignments, wolf.ID)
	}

	for i := uint64(0); i < factorial(roles.Total); i++ {
		var universe Universe
		universe.permutation = i
		universe.active = true
		log.Print(universe)
	}
}
