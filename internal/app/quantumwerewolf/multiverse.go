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
					collapseForFixedRole(p.Num, int(fixedRoleID))
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

func collapseForFixedRole(playerNumber int, fixedRoleID int) {
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

	if roleTypes[fixedRoleID].CanAttack {
		collapseForAttack(playerNumber)
	}

	if roleTypes[fixedRoleID].CanPeek {
		collapseForPeek(playerNumber)
	}
}

// collapseForAttack should only be called if role is fixed to attacker
func collapseForAttack(attacker int) {
	universeLength := len(multiverse.originalAssignments)
	evaluationUniverse := make([]int, universeLength)
	evaluationRanks := make([]int, universeLength)
	newUniverses := make([]uint64, 0, len(multiverse.universes))
	universesEliminated := false

	attackTargets := make([]int, 0, len(players))
	actionStrings := strings.Split(players[attacker].Actions, tokenEndAction)
	for _, a := range actionStrings {
		attackIndex := strings.Index(a, tokenAttack)
		if attackIndex < 0 {
			// This is not an attack action
			continue
		}
		attackTarget := a[attackIndex:len(a)]
		attackTargets = append(attackTargets, getPlayerByName(attackTarget).Num)
	}

	for _, v := range multiverse.universes {
		copy(evaluationUniverse, multiverse.originalAssignments)
		for i := range evaluationRanks {
			evaluationRanks[i] = i
		}
		evaluationUniverse = quantumutilities.Kthperm(evaluationUniverse, v)
		evaluationRanks = quantumutilities.Kthperm(evaluationRanks, v)

		attackSucceeds := false
		if roleTypes[evaluationUniverse[attacker]].CanAttack {
			highestRankedAttacker := true
			// Check if potential is highest ranked attacker in this universe
			for i := range evaluationUniverse {
				// If same role ID
				if evaluationUniverse[i] == evaluationUniverse[attacker] {
					// If someone else has higher rank
					if evaluationRanks[i] < evaluationRanks[attacker] {
						// If higher ranked was still alive when attack was made TODO CHECK WHEN
						if strings.Index(players[i].Actions, tokenKilled) < 0 {
							// Can't attack due to low rank
							highestRankedAttacker = false
							break
						}
					}
				}
			}

			if highestRankedAttacker {
				for _, target := range attackTargets {
					// Can attack if on other side
					if roleTypes[evaluationUniverse[target]].Evil != roleTypes[evaluationUniverse[attacker]].Evil {
						attackSucceeds = true
					}
				}
			}
		}

		if attackSucceeds {
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

// collapseForPeek should only be called if role is fixed to peeker
func collapseForPeek(peeker int) {
	universeLength := len(multiverse.originalAssignments)
	evaluationUniverse := make([]int, universeLength)
	newUniverses := make([]uint64, 0, len(multiverse.universes))
	universesEliminated := false

	type PeekAction struct {
		Target int
		IsEvil bool
	}

	peekActions := make([]PeekAction, 0, len(players))
	actionStrings := strings.Split(players[peeker].Actions, tokenEndAction)
	for _, a := range actionStrings {
		peekIndex := strings.Index(a, tokenPeek)
		if peekIndex < 0 {
			// This is not a peek action
			continue
		}
		peekTarget := a[peekIndex : len(a)-1]
		peekResult := a[len(a)-1 : len(a)]

		peekAction := PeekAction{}
		peekAction.Target = getPlayerByName(peekTarget).Num
		peekAction.IsEvil = peekResult == tokenEvil
		peekActions = append(peekActions, peekAction)
	}

	for _, v := range multiverse.universes {
		copy(evaluationUniverse, multiverse.originalAssignments)
		evaluationUniverse = quantumutilities.Kthperm(evaluationUniverse, v)

		peekSucceeds := false
		if roleTypes[evaluationUniverse[peeker]].CanPeek {
			peekSucceeds = true
			for _, a := range peekActions {
				targetIsEvil := roleTypes[evaluationUniverse[a.Target]].Evil

				// Observation doesn't match this universe's reality
				if a.IsEvil != targetIsEvil {
					peekSucceeds = false
					break
				}
			}
		}

		if peekSucceeds {
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
	collapseForFixedRole(playerNumber, roleID)
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
	evaluationUniverse := make([]int, universeLength)
	// Keep trying until a universe is found where this player is a seer
	for true {
		copy(evaluationUniverse, multiverse.originalAssignments)
		evaluationUniverse = quantumutilities.Kthperm(evaluationUniverse, randomUniverse())
		if evaluationUniverse[potentialSeer] == seer.ID {
			return roleTypes[evaluationUniverse[target]].Evil
		}
	}

	return false
}
