package quantumwerewolf

import (
	"fmt"
	"log"
	"math/rand"
	"quantumwerewolf/pkg/quantumutilities"
)

func getUniverseString(universe uint64) string {
	var universeString string

	// Add roles display
	universeAssignments := make([]int, len(multiverse.originalAssignments))
	copy(universeAssignments, multiverse.originalAssignments)
	universeAssignments = quantumutilities.Kthperm(universeAssignments, universe)

	universeString += "["
	for i, v := range universeAssignments {
		if i > 0 {
			universeString += " "
		}
		universeString += roleTypes[v].Name[:1]
	}
	universeString += "]"

	return fmt.Sprint(universeString)
}

func createMultiverse() {
	setupRoles()
	for i := 0; i < roles.Villagers; i++ {
		multiverse.originalAssignments = append(multiverse.originalAssignments, villager.ID)
	}
	for i := 0; i < roles.Seers; i++ {
		multiverse.originalAssignments = append(multiverse.originalAssignments, seer.ID)
	}
	for i := 0; i < roles.Wolves; i++ {
		multiverse.originalAssignments = append(multiverse.originalAssignments, wolf.ID)
	}

	randSource := rand.NewSource(game.Seed)
	multiverseRandom := rand.New(randSource)
	possibleUniverses := quantumutilities.Factorial(roles.Total)
	multiverse.universes = quantumutilities.PermUint64Trunc(multiverseRandom, possibleUniverses, 100000)

	for _, v := range multiverse.universes {
		log.Printf(getUniverseString(v))
	}
}
