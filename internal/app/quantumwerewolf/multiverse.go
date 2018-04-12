package quantumwerewolf

import (
	"fmt"
	"log"
	"math/rand"
	"quantumwerewolf/pkg/quantumutilities"
	"strconv"
)

func getUniverseString(universe uint64) string {
	var universeString string

	// Add gameSetup display
	universeLength := len(multiverse.originalAssignments)
	universeRoleIDs := make([]int, universeLength)
	copy(universeRoleIDs, multiverse.originalAssignments)
	universeRoleIDs = quantumutilities.Kthperm(universeRoleIDs, universe)

	universeRanks := make([]int, universeLength)
	for i := range universeRanks {
		universeRanks[i] = i
	}
	universeRanks = quantumutilities.Kthperm(universeRanks, universe)

	universeString += "["
	for i, v := range universeRoleIDs {
		if i > 0 {
			universeString += " "
		}
		universeString += roleTypes[v].Name[:1]
		universeString += strconv.Itoa(universeRanks[i])
	}
	universeString += "]"

	return fmt.Sprint(universeString)
}

func createMultiverse() {
	setupRoles()
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

func updateRoleTotals() {
	for i := range players {
		for j := range players[i].Role {
			players[i].Role[j] = 0
		}
	}

	universeLength := len(multiverse.originalAssignments)
	universeRoleIDs := make([]int, universeLength)
	for _, v := range multiverse.universes {
		copy(universeRoleIDs, multiverse.originalAssignments)
		universeRoleIDs = quantumutilities.Kthperm(universeRoleIDs, v)
		for i, uv := range universeRoleIDs {
			players[i].Role[uv]++
		}
	}
}

func randomUniverse() uint64 {
	universeNumber := uint64(rand.Int63n(int64(len(multiverse.universes))))
	getUniverseString(universeNumber)
	return universeNumber
}

func getFixedRole(playerNumber int) int {
	universeLength := len(multiverse.originalAssignments)
	foundUniverse := make([]int, universeLength)
	copy(foundUniverse, multiverse.originalAssignments)
	foundUniverse = quantumutilities.Kthperm(foundUniverse, randomUniverse())

	return foundUniverse[playerNumber]
}

func collapseToFixedRole(playerNumber int) {
	roleID := getFixedRole(playerNumber)

	// TODO: Eliminate cases for peeks and attacks
	universeLength := len(multiverse.originalAssignments)
	universeRoleIDs := make([]int, universeLength)
	for i, v := range multiverse.universes {
		copy(universeRoleIDs, multiverse.originalAssignments)
		universeRoleIDs = quantumutilities.Kthperm(universeRoleIDs, v)
		if universeRoleIDs[playerNumber] != roleID {
			multiverse.universes = append(multiverse.universes[:i], multiverse.universes[i+1:]...)
		}
	}
}

// peek returns true if the playernumber is evil
func peek(potentialSeer int, target int) bool {
	updateRoleTotals()

	if players[potentialSeer].Role[seer.ID] == 0 {
		log.Printf("Attempted to peek with player %d but can not peek", potentialSeer)
		return false
	}

	universeLength := len(multiverse.originalAssignments)
	foundUniverse := make([]int, universeLength)
	// Keep trying until a universe is found where this player is a seer
	for true {
		copy(foundUniverse, multiverse.originalAssignments)
		foundUniverse = quantumutilities.Kthperm(foundUniverse, randomUniverse())
		if foundUniverse[potentialSeer] == seer.ID {
			return roleTypes[foundUniverse[target]].Evil
		}
	}

	return false
}
