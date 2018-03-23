package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Universe is a single unit in multiverse
type Universe struct {
	assignment []int
	active     bool
}

var multiverse []Universe

func (u Universe) String() string {
	var universeString string

	// Add active display
	if u.active {
		universeString = "A"
	} else {
		universeString = "I"
	}

	// Add roles display
	universeString += "["
	for i, v := range u.assignment {
		if i > 0 {
			universeString += " "
		}
		universeString += roleTypes[v].Name
	}
	universeString += "]"
	return fmt.Sprint(universeString)
}

func nextPerm(p []int) {
	for i := len(p) - 1; i >= 0; i-- {
		if i == 0 || p[i] < len(p)-i-1 {
			p[i]++
			return
		}
		p[i] = 0
	}
}

func getPerm(orig, p []int) []int {
	result := append([]int{}, orig...)
	for i, v := range p {
		result[i], result[i+v] = result[i+v], result[i]
	}
	return result
}

func createMultiverse() {
	setupRoles()
	var originalAssignments []int
	for i := 0; i < roles.Villagers; i++ {
		originalAssignments = append(originalAssignments, villager.ID)
	}
	for i := 0; i < roles.Seers; i++ {
		originalAssignments = append(originalAssignments, seer.ID)
	}
	for i := 0; i < roles.Wolves; i++ {
		originalAssignments = append(originalAssignments, wolf.ID)
	}

	for p := make([]int, len(originalAssignments)); p[0] < len(p); nextPerm(p) {
		var universe Universe
		universe.assignment = getPerm(originalAssignments, p)
		universe.active = true
		multiverse = append(multiverse, universe)
		log.Print(universe)
	}
}
