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
	collapseAll()

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

	for _, p := range players {
		for j := range p.Role.Totals {
			p.Role.Totals[j] = 0
		}
	}

	universeLength := len(multiverse.originalAssignments)
	universeRoleIDs := make([]int, universeLength)
	for _, v := range multiverse.universes {
		copy(universeRoleIDs, multiverse.originalAssignments)
		universeRoleIDs = quantumutilities.Kthperm(universeRoleIDs, v)
		for i, uv := range universeRoleIDs {
			getPlayerByNumber(i).Role.Totals[uv]++
		}
	}

	for _, p := range players {
		numberOfPossibleRoles := 0
		fixedRole := 0
		for _, r := range roleTypes {
			if p.Role.Totals[r.ID] > 0 {
				numberOfPossibleRoles++
				fixedRole = r.ID
			}
		}
		if numberOfPossibleRoles == 1 {
			p.Role.IsFixed = true
			p.Role.Fixed = fixedRole
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

func collapseAll() {
	UpdateRoleTotals()
	anyEliminations := false
	for _, p := range players {
		if p.Role.IsFixed {
			anyEliminations = collapseForFixedRole(p.Num, p.Role.Fixed)
		}

		if playerCanAttack(p) {
			anyEliminations = anyEliminations || collapseForAttack(p.Num)
		}

		if playerCanPeek(p) {
			anyEliminations = anyEliminations || collapseForPeek(p.Num)
		}
	}

	anyEliminations = anyEliminations || collapseForPriorDeaths()

	if anyEliminations {
		collapseAll()
	}
}

func collapseForFixedRole(playerNumber int, fixedRoleID int) bool {
	FillObservations()
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

	anyEliminations := universesEliminated
	eliminateUniverses(universesEliminated, newUniverses)

	return anyEliminations
}

// AttackTarget returns true if the attacker successfully completed an attack on this target in the given universe
func AttackTarget(universe uint64, attacker int, target int) bool {
	universeLength := len(multiverse.originalAssignments)
	evaluationUniverse := make([]int, universeLength)
	evaluationRanks := make([]int, universeLength)

	copy(evaluationUniverse, multiverse.originalAssignments)
	for i := range evaluationRanks {
		evaluationRanks[i] = i
	}

	evaluationUniverse = quantumutilities.Kthperm(evaluationUniverse, universe)
	evaluationRanks = quantumutilities.Kthperm(evaluationRanks, universe)

	attackSucceeds := false
	if roleTypes[evaluationUniverse[attacker]].CanAttack {
		// Can only attack if on other side from target
		if roleTypes[evaluationUniverse[target]].Evil != roleTypes[evaluationUniverse[attacker]].Evil {
			attackSucceeds = true

			// Check if potential is highest ranked attacker in this universe
			for teammateIndex := range evaluationUniverse {
				// If same role ID
				if evaluationUniverse[teammateIndex] == evaluationUniverse[attacker] {
					// If someone else has higher rank
					if evaluationRanks[teammateIndex] < evaluationRanks[attacker] {
						wasTeammateDead := false
						for _, teammateKilled := range killObservations {
							// If higher ranked was dead when attack was made
							if !teammateKilled.Pending && teammateKilled.Subject == teammateIndex && teammateKilled.Round > attacker {
								wasTeammateDead = true
								break
							}
						}
						if !wasTeammateDead {
							attackSucceeds = false
							break
						}
					}
				}
			}
		}
	}

	return attackSucceeds
}

// AttackFriend returns true if the attacker kills a teammate
func AttackFriend(universe uint64, attacker int, target int) bool {
	universeLength := len(multiverse.originalAssignments)
	evaluationUniverse := make([]int, universeLength)
	evaluationRanks := make([]int, universeLength)

	copy(evaluationUniverse, multiverse.originalAssignments)
	for i := range evaluationRanks {
		evaluationRanks[i] = i
	}

	evaluationUniverse = quantumutilities.Kthperm(evaluationUniverse, universe)
	evaluationRanks = quantumutilities.Kthperm(evaluationRanks, universe)

	attackedFriend := false
	if roleTypes[evaluationUniverse[attacker]].CanAttack {
		// If on same side as target
		if roleTypes[evaluationUniverse[target]].Evil == roleTypes[evaluationUniverse[attacker]].Evil {
			attackedFriend = true

			// Check if potential is highest ranked attacker in this universe
			for teammateIndex := range evaluationUniverse {
				// If same role ID
				if evaluationUniverse[teammateIndex] == evaluationUniverse[attacker] {
					// If someone else has higher rank
					if evaluationRanks[teammateIndex] < evaluationRanks[attacker] {
						wasTeammateDead := false
						for _, teammateKilled := range killObservations {
							// If higher ranked was dead when attack was made
							if !teammateKilled.Pending && teammateKilled.Subject == teammateIndex && teammateKilled.Round > attacker {
								wasTeammateDead = true
								break
							}
						}
						if !wasTeammateDead {
							attackedFriend = false
							break
						}
					}
				}
			}
		}
	}

	return attackedFriend
}

// PeekOk returns false if a player is a seer and gets an untrue result
func PeekOk(universe uint64, peeker int) bool {
	if getPlayerByNumber(peeker).Role.IsFixed && roleTypes[getPlayerByNumber(peeker).Role.Fixed].CanPeek {
		universeLength := len(multiverse.originalAssignments)
		evaluationUniverse := make([]int, universeLength)
		copy(evaluationUniverse, multiverse.originalAssignments)

		FillObservations()
		for _, peek := range peekObservations {
			if !peek.Pending && roleTypes[evaluationUniverse[peeker]].CanPeek {
				if roleTypes[evaluationUniverse[peek.Target]].Evil != peek.IsEvil {
					return false
				}
			}
		}
	}

	return true
}

// AttackOk returns false if a player is the dominant attacker and attacks a teammate
func AttackOk(universe uint64, attacker int) bool {
	universeLength := len(multiverse.originalAssignments)
	evaluationUniverse := make([]int, universeLength)
	copy(evaluationUniverse, multiverse.originalAssignments)

	FillObservations()
	for _, attack := range attackObservations {
		if !attack.Pending && AttackFriend(universe, attacker, attack.Target) {
			return false
		}
	}

	return true
}

func collapseForAttack(attacker int) bool {
	newUniverses := make([]uint64, 0, len(multiverse.universes))
	universesEliminated := false

	for _, v := range multiverse.universes {
		if AttackOk(v, attacker) {
			newUniverses = append(newUniverses, v)
		} else {
			universesEliminated = true
		}
	}

	// TODO add a killed observation if someone is dead in all universes

	eliminateUniverses(universesEliminated, newUniverses)
	return universesEliminated
}

// collapseForPriorDeaths eliminates universes where a lynchee was attacked befoe
func collapseForPriorDeaths() bool {
	// TODO fill this in
	return false
}

func collapseForPeek(peeker int) bool {
	newUniverses := make([]uint64, 0, len(multiverse.universes))
	universesEliminated := false

	for _, v := range multiverse.universes {
		if PeekOk(v, peeker) {
			newUniverses = append(newUniverses, v)
		} else {
			universesEliminated = true
		}
	}

	eliminateUniverses(universesEliminated, newUniverses)
	return universesEliminated
}

func eliminateUniverses(universesEliminated bool, newUniverses []uint64) {
	if universesEliminated && len(newUniverses) > 0 {
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
	collapseAll()
	return roleID
}

// Peek returns true if the playernumber is evil
func Peek(potentialSeer int, target int) bool {
	UpdateRoleTotals()

	if getPlayerByNumber(potentialSeer).Role.Totals[seer.ID] == 0 {
		log.Printf("Attempted to peek with player %d but can not peek", potentialSeer)
	}

	universeLength := len(multiverse.originalAssignments)
	evaluationUniverse := make([]int, universeLength)
	// Keep trying until a universe is found where this player is a seer
	for true {
		copy(evaluationUniverse, multiverse.originalAssignments)
		evaluationUniverse = quantumutilities.Kthperm(evaluationUniverse, randomUniverse())
		if evaluationUniverse[potentialSeer] == seer.ID {
			if roleTypes[evaluationUniverse[target]].Evil {
				return true
			}

			return false
		}
	}

	return false
}
