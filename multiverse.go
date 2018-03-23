package main

import (
	"log"

	_ "github.com/lib/pq"
)

type Universe struct {
	assignment []string
	active     bool
}

var multiverse []Universe

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
		originalAssignments = append(originalAssignments, villager.Id)
	}
	for i := 0; i < roles.Seers; i++ {
		originalAssignments = append(originalAssignments, seer.Id)
	}
	for i := 0; i < roles.Wolves; i++ {
		originalAssignments = append(originalAssignments, wolf.Id)
	}

	for p := make([]int, len(originalAssignments)); p[0] < len(p); nextPerm(p) {
		var perm = getPerm(originalAssignments, p)
		var universeString = "["
		for i, v := range perm {
			if i > 0 {
				universeString += " "
			}
			universeString += roleTypes[v].Name
		}
		universeString += "]"
		log.Print(universeString)
	}
}
