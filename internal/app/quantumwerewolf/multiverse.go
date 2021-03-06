package quantumwerewolf

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"quantumwerewolf/pkg/quantumutilities"
	"sort"
	"strconv"
)

var dirtyMultiverse bool

func getUniverseString(universe uint64) string {
	var universeString string

	universeLength := len(Multiverse.originalAssignments)
	universeRoleIDs := quantumutilities.Kthperm(Multiverse.originalAssignments, universe)

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

// CreateMultiverse creates the entire Multiverse based on the current game state
func CreateMultiverse() {
	dirtyMultiverse = true
	setupRoles()

	// Make sure to iterate through roleTypes in the same order each time
	keys := make([]int, 0)
	for k := range roleTypes {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		for j := 0; j < GameSetup.Roles[roleTypes[k].Name]; j++ {
			Multiverse.originalAssignments = append(Multiverse.originalAssignments, roleTypes[k].ID)
		}
	}

	randSource := rand.NewSource(Game.Seed)
	Multiverse.rando = rand.New(randSource)
	Multiverse.Universes = quantumutilities.PermUint64Trunc(Multiverse.rando, GameSetup.Universes, 100000)
	UpdateRoleTotals()
}

func updateFixedRoles() {
	for _, o := range observations {
		if o.getType() == "KillObservation" {
			role, err := o.getRole()
			if err == nil {
				CollapseForFixedRole(o.getSubject(), role)
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

	if len(Multiverse.Universes) == 0 {
		return
	}

	updateFixedRoles()

	for _, p := range Players {
		for j := range p.Role.Totals {
			Players[getPlayerIndex(p)].Role.Totals[j] = 0
		}
	}

	for _, v := range Multiverse.Universes {
		universeRoleIDs := quantumutilities.Kthperm(Multiverse.originalAssignments, v)
		for i, uv := range universeRoleIDs {
			Players[getPlayerIndex(getPlayerByNumber(i))].Role.Totals[uv]++
		}
	}

	for _, p := range Players {
		numberOfPossibleRoles := 0
		fixedRole := 0
		for _, r := range roleTypes {
			if p.Role.Totals[r.ID] > 0 {
				numberOfPossibleRoles++
				fixedRole = r.ID
			}
		}
		if numberOfPossibleRoles == 1 {
			Players[getPlayerIndex(p)].Role.IsFixed = true
			Players[getPlayerIndex(p)].Role.Fixed = fixedRole
		}
	}

	dirtyMultiverse = false
}

func randomUniverse() uint64 {
	universeNumber := Multiverse.Universes[uint64(Multiverse.rando.Int63n(int64(len(Multiverse.Universes))))]
	return universeNumber
}

func randomPeekUniverse(peeker int) (uint64, error) {
	temp := make([]uint64, 0, len(Multiverse.Universes))

	for _, v := range Multiverse.Universes {
		evaluationUniverse := quantumutilities.Kthperm(Multiverse.originalAssignments, v)
		if roleTypes[evaluationUniverse[peeker]].CanPeek {
			temp = append(temp, v)
		}
	}

	if len(temp) < 1 {
		return 0, errors.New("No valid peek universe found")
	}
	return temp[uint64(Multiverse.rando.Int63n(int64(len(temp))))], nil
}

func getFixedRole(playerNumber int) int {
	foundUniverse := quantumutilities.Kthperm(Multiverse.originalAssignments, randomUniverse())

	return foundUniverse[playerNumber]
}

// CollapseAll removes universes that are inconsistent for any reason
func CollapseAll() {
	CollapseAllUpTo(Game.RoundNum)
}

// CollapseAllUpTo removes universes that are inconsistent for any reason up to the input round
func CollapseAllUpTo(maxRound int) {
	dirtyMultiverse = true
	FillObservations()
	UpdateRoleTotals()

	for round := 0; round < maxRound; round++ {
		collapsed := true

		for i := 0; collapsed && i < len(Players); i++ {
			collapsed = false

			anyEliminations := true
			for anyEliminations {
				anyEliminations = false
				for _, p := range Players {
					if p.Role.IsFixed {
						if CollapseForFixedRole(p.Num, p.Role.Fixed) {
							anyEliminations = true
							collapsed = true
						}
					}
				}
			}

			for _, o := range observations {
				anyEliminations = true
				for anyEliminations {
					anyEliminations = false

					newUniverses := make([]uint64, 0, len(Multiverse.Universes))
					universesEliminated := false

					for _, v := range Multiverse.Universes {
						if o.collapse(round, v) {
							universesEliminated = true
						} else {
							newUniverses = append(newUniverses, v)
						}
					}

					universesEliminated = eliminateUniverses(universesEliminated, newUniverses)

					if universesEliminated {
						anyEliminations = true
						collapsed = true
					}
				}
			}
		}
	}

	UpdateRoleTotals()
}

// CollapseForFixedRole removes universes that are inconsistent with the role given to the player
// Param playerNumber: The number of the player to collapse for
// Param fixedRoleId: The role id that the player is fixed to
func CollapseForFixedRole(playerNumber int, fixedRoleID int) bool {
	newUniverses := make([]uint64, 0, len(Multiverse.Universes))
	universesEliminated := false

	for _, v := range Multiverse.Universes {
		universeRoleIDs := quantumutilities.Kthperm(Multiverse.originalAssignments, v)
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
	evaluationRanks := make([]int, len(Multiverse.originalAssignments))
	for i := range evaluationRanks {
		evaluationRanks[i] = i
	}
	evaluationRanks = quantumutilities.Kthperm(evaluationRanks, universe)
	evaluationUniverse := quantumutilities.Kthperm(Multiverse.originalAssignments, universe)

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
						for _, teammateKilled := range observations {
							// If higher ranked was dead when attack was made
							if !teammateKilled.getPending() && teammateKilled.getSubject() == teammateIndex && teammateKilled.getRound() > attacker && teammateKilled.getType() == "KillObservation" {
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
	universeLength := len(Multiverse.originalAssignments)
	evaluationRanks := make([]int, universeLength)
	for i := range evaluationRanks {
		evaluationRanks[i] = i
	}
	evaluationRanks = quantumutilities.Kthperm(evaluationRanks, universe)
	evaluationUniverse := quantumutilities.Kthperm(Multiverse.originalAssignments, universe)

	type Attacker struct {
		player    Player
		deadNight int
		rank      int
	}

	attackers := make([]Attacker, 0, universeLength)
	for i, r := range evaluationUniverse {
		if roleTypes[r].CanAttack && roleTypes[r].Evil == evilSide {
			var attacker Attacker
			attacker.player = Players[i]
			attacker.rank = evaluationRanks[getPlayerIndex(attacker.player)]
			attacker.deadNight = 10000 // Initialize to infinity
			for _, o := range observations {
				if o.getSubject() == attacker.player.Num && o.getType() == "KillObservation" {
					attacker.deadNight = o.getRound()
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
	evaluationUniverse := quantumutilities.Kthperm(Multiverse.originalAssignments, universe)

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

func eliminateUniverses(universesEliminated bool, newUniverses []uint64) bool {
	if universesEliminated && len(newUniverses) > 0 {
		dirtyMultiverse = true
		Multiverse.Universes = make([]uint64, 0, len(newUniverses))
		for _, v := range newUniverses {
			Multiverse.Universes = append(Multiverse.Universes, v)
		}

		return true
	}

	if universesEliminated {
		log.Printf("Attempted to remove all universes")
	}

	return false
}

func collapseToFixedRole(playerNumber int) int {
	CollapseAll()
	roleID := getFixedRole(playerNumber)
	CollapseForFixedRole(playerNumber, roleID)
	CollapseAll()
	return roleID
}

// Peek returns true if the playernumber is evil
func Peek(peeker int, target int) bool {
	UpdateRoleTotals()

	if playerCanPeek(getPlayerByNumber(peeker)) {
		peekUniverse, err := randomPeekUniverse(peeker)
		if err != nil {
			log.Printf("Attempted to peek with player %d but had error: %v", peeker, err)
			return false
		}

		evaluationUniverse := quantumutilities.Kthperm(Multiverse.originalAssignments, peekUniverse)
		return roleTypes[evaluationUniverse[target]].Evil
	}

	log.Printf("Attempted to peek with player %d but can not peek", peeker)
	return false
}
