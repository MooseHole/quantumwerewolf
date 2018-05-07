package quantumwerewolf

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"quantumwerewolf/pkg/quantumutilities"
	"strconv"
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
	multiverse.universes = quantumutilities.PermUint64Trunc(multiverseRandom, gameSetup.Universes, 100000)
	UpdateRoleTotals()
	collapseAll()
}

func updateFixedRoles() {
	for _, o := range killObservations {
		collapseForFixedRole(o.Subject, o.Role)
	}
}

// UpdateRoleTotals figures everything out based on actions and fixed roles
// It is slow so it uses dirtyMultiverse to reduce iterations
func UpdateRoleTotals() {
	if !dirtyMultiverse {
		return
	}

	if len(multiverse.universes) == 0 {
		return
	}

	updateFixedRoles()

	for _, p := range players {
		for j := range p.Role.Totals {
			players[getPlayerIndex(p)].Role.Totals[j] = 0
		}
	}

	universeLength := len(multiverse.originalAssignments)
	universeRoleIDs := make([]int, universeLength)
	for _, v := range multiverse.universes {
		copy(universeRoleIDs, multiverse.originalAssignments)
		universeRoleIDs = quantumutilities.Kthperm(universeRoleIDs, v)
		for i, uv := range universeRoleIDs {
			players[getPlayerIndex(getPlayerByNumber(i))].Role.Totals[uv]++
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
			players[getPlayerIndex(p)].Role.IsFixed = true
			players[getPlayerIndex(p)].Role.Fixed = fixedRole
		}
	}

	dirtyMultiverse = false
}

func randomUniverse() uint64 {
	universeNumber := multiverse.universes[uint64(rand.Int63n(int64(len(multiverse.universes))))]
	return universeNumber
}

func randomPeekUniverse(peeker int) (uint64, error) {
	temp := make([]uint64, 0, len(multiverse.universes))
	universeLength := len(multiverse.originalAssignments)
	evaluationUniverse := make([]int, universeLength)

	for _, v := range multiverse.universes {
		copy(evaluationUniverse, multiverse.originalAssignments)
		evaluationUniverse = quantumutilities.Kthperm(evaluationUniverse, v)
		if roleTypes[evaluationUniverse[peeker]].CanPeek {
			temp = append(temp, v)
		}
	}

	if len(temp) < 1 {
		return 0, errors.New("No valid peek universe found")
	}
	return temp[uint64(rand.Int63n(int64(len(temp))))], nil
}

func getFixedRole(playerNumber int) int {
	universeLength := len(multiverse.originalAssignments)
	foundUniverse := make([]int, universeLength)
	copy(foundUniverse, multiverse.originalAssignments)
	foundUniverse = quantumutilities.Kthperm(foundUniverse, randomUniverse())

	return foundUniverse[playerNumber]
}

func collapseAll() {
	dirtyMultiverse = true
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

	universesEliminated = eliminateUniverses(universesEliminated, newUniverses)

	return universesEliminated
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

// DominantAttacker returns the player that was the dominant attacker on the night in question in the given universe
func DominantAttacker(universe uint64, night int, evilSide bool) Player {
	universeLength := len(multiverse.originalAssignments)
	evaluationUniverse := make([]int, universeLength)
	evaluationRanks := make([]int, universeLength)

	copy(evaluationUniverse, multiverse.originalAssignments)
	for i := range evaluationRanks {
		evaluationRanks[i] = i
	}

	evaluationUniverse = quantumutilities.Kthperm(evaluationUniverse, universe)
	evaluationRanks = quantumutilities.Kthperm(evaluationRanks, universe)

	type Attacker struct {
		player    Player
		deadNight int
		rank      int
	}

	attackers := make([]Attacker, 0, universeLength)
	for i, r := range evaluationUniverse {
		if roleTypes[r].CanAttack && roleTypes[r].Evil == evilSide {
			var attacker Attacker
			attacker.player = players[i]
			attacker.rank = evaluationRanks[getPlayerIndex(attacker.player)]
			attacker.deadNight = 10000 // Initialize to infinity
			for _, o := range killObservations {
				if o.Subject == attacker.player.Num {
					attacker.deadNight = o.Round
				}
			}
			attackers = append(attackers, attacker)
		}
	}

	maxRank := 0
	for _, a := range attackers {
		if a.deadNight > night {
			maxRank = int(math.Max(float64(maxRank), float64(a.rank)))
		}
	}

	for _, a := range attackers {
		if a.rank == maxRank {
			return a.player
		}
	}

	log.Printf("Attempted to get unknown dominant attacker  universe %d  night %d  evilSide %v", universe, night, evilSide)
	var unknownPlayer Player
	return unknownPlayer
}

// AttackFriend returns true if the attacker kills a teammate
func AttackFriend(universe uint64, attacker int, target int, night int) bool {
	universeLength := len(multiverse.originalAssignments)
	evaluationUniverse := make([]int, universeLength)
	copy(evaluationUniverse, multiverse.originalAssignments)
	evaluationUniverse = quantumutilities.Kthperm(evaluationUniverse, universe)

	attackedFriend := false
	if roleTypes[evaluationUniverse[attacker]].CanAttack {
		// If on same side as target this is not ok
		if roleTypes[evaluationUniverse[target]].Evil == roleTypes[evaluationUniverse[attacker]].Evil {
			attackedFriend = true

			// If the attacker is mot the dominant though
			dominant := DominantAttacker(universe, night, roleTypes[evaluationUniverse[attacker]].Evil)
			if dominant.Num != attacker {
				attackedFriend = false
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
	FillObservations()
	for _, attack := range attackObservations {
		if !attack.Pending && AttackFriend(universe, attacker, attack.Target, attack.Round) {
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

	universesEliminated = eliminateUniverses(universesEliminated, newUniverses)
	return universesEliminated
}

// collapseForPriorDeaths eliminates universes where a lynchee was attacked before
func collapseForPriorDeaths() bool {
	newUniverses := make([]uint64, 0, len(multiverse.universes))
	universesEliminated := false

	for _, v := range multiverse.universes {
		eliminate := false
		for _, lynch := range lynchObservations {
			for _, kill := range killObservations {
				if lynch.Subject == kill.Subject && lynch.Round != kill.Round {
					eliminate = true
					break
				}
			}
			if eliminate {
				break
			}
		}

		if !eliminate {
			newUniverses = append(newUniverses, v)
		} else {
			universesEliminated = true
		}
	}

	universesEliminated = eliminateUniverses(universesEliminated, newUniverses)
	return universesEliminated
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

	universesEliminated = eliminateUniverses(universesEliminated, newUniverses)
	return universesEliminated
}

func eliminateUniverses(universesEliminated bool, newUniverses []uint64) bool {
	if universesEliminated && len(newUniverses) > 0 {
		dirtyMultiverse = true
		multiverse.universes = make([]uint64, 0, len(newUniverses))
		for _, v := range newUniverses {
			multiverse.universes = append(multiverse.universes, v)
		}

		return true
	}

	return false
}

func collapseToFixedRole(playerNumber int) int {
	roleID := getFixedRole(playerNumber)
	collapseForFixedRole(playerNumber, roleID)
	collapseAll()
	return roleID
}

// Peek returns true if the playernumber is evil
func Peek(peeker int, target int) bool {
	UpdateRoleTotals()

	if playerCanPeek(getPlayerByNumber(peeker)) {
		log.Printf("Attempted to peek with player %d but can not peek", peeker)
	}

	peekUniverse, err := randomPeekUniverse(peeker)
	if err != nil {
		log.Printf("Attempted to peek with player %d but had error: %v", peeker, err)
		return false
	}

	evaluationUniverse := make([]int, len(multiverse.originalAssignments))
	copy(evaluationUniverse, multiverse.originalAssignments)
	evaluationUniverse = quantumutilities.Kthperm(evaluationUniverse, peekUniverse)
	return roleTypes[evaluationUniverse[target]].Evil
}
