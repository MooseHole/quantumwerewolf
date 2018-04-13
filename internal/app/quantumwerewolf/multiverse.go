package quantumwerewolf

import (
	"fmt"
	"log"
	"math/rand"
	"quantumwerewolf/pkg/quantumutilities"
	"strconv"
	"strings"
)

var dirtyMultiverse bool

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

// CreateMultiverse creates the entire multiverse based on the current game state
func CreateMultiverse() {
	dirtyMultiverse = true
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
	UpdateRoleTotals()

	for _, v := range multiverse.universes {
		log.Printf(getUniverseString(v))
	}
}

func updateFixedRoles() {
	for _, p := range players {
		actionStrings := strings.Split(p.Actions, tokenEndAction)
		for _, a := range actionStrings {
			killedIndex := strings.Index(a, tokenKilled)
			if killedIndex >= 0 {
				fixedRoleIDString := a[killedIndex+1:]
				fixedRoleID, err := strconv.ParseInt(fixedRoleIDString, 10, 64)
				if err == nil {
					collapse(p.Num, int(fixedRoleID))
				} else {
					log.Printf("updateFixedRoles had error when parsing role id: %v", err)
				}
			}
		}
	}
}

// UpdateRoleTotals figures everything out based on actions and fixed roles
// It is slow so it uses dirtyMultiverse to reduce iterations
func UpdateRoleTotals() {
	if !dirtyMultiverse {
		return
	}

	updateFixedRoles()

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

	dirtyMultiverse = false
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

func collapse(playerNumber int, fixedRoleID int) {
	// TODO: Eliminate cases for peeks and attacks
	universeLength := len(multiverse.originalAssignments)
	universeRoleIDs := make([]int, universeLength)
	newUniverses := make([]uint64, 0, len(multiverse.universes))
	universesEliminated := false

	for _, v := range multiverse.universes {
		copy(universeRoleIDs, multiverse.originalAssignments)
		universeRoleIDs = quantumutilities.Kthperm(universeRoleIDs, v)
		if universeRoleIDs[playerNumber] == fixedRoleID {
			newUniverses = append(newUniverses, v)
		} else {
			universesEliminated = true
		}
	}

	if universesEliminated {
		dirtyMultiverse = true
		multiverse.universes = make([]uint64, 0, len(newUniverses))
		for _, v := range newUniverses {
			multiverse.universes = append(multiverse.universes, v)
		}
	}
}

func collapseToFixedRole(playerNumber int) int {
	roleID := getFixedRole(playerNumber)
	collapse(playerNumber, roleID)
	return roleID
}

// Peek returns true if the playernumber is evil
func Peek(potentialSeer int, target int) bool {
	UpdateRoleTotals()

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
			dirtyMultiverse = true
			return roleTypes[foundUniverse[target]].Evil
		}
	}

	return false
}
