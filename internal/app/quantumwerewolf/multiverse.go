package quantumwerewolf

import (
	"fmt"
	"log"
	"math/rand"
	"quantumwerewolf/pkg/quantumutilities"
)

func getUniverseString(universe uint64) string {
	var universeString string

	// Add gameSetup display
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
	for _, v := range roleTypes {
		for j := 0; j < gameSetup.Roles[v.Name]; j++ {
			multiverse.originalAssignments = append(multiverse.originalAssignments, v.ID)
		}
	}

	randSource := rand.NewSource(game.Seed)
	multiverseRandom := rand.New(randSource)
	possibleUniverses := quantumutilities.Factorial(gameSetup.Total)
	multiverse.universes = quantumutilities.PermUint64Trunc(multiverseRandom, possibleUniverses, 100000)

	for _, v := range multiverse.universes {
		log.Printf(getUniverseString(v))
	}
}
