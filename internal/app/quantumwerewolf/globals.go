package quantumwerewolf

import (
	"log"
	"math/rand"
	"time"
)

// Tokens for the actions
// 3@Alice| means the player attacked Alice in round 3
// 0%Bob|2&Carol|3&Carol|4#| means the player peeked at Bob in round 0, voted to lynch Carol in rounds 2 and 3, and died in round 4
const tokenAttack = "@"
const tokenPeek = "%"
const tokenVote = "&"
const tokenKilled = "#"
const tokenLynched = "+"
const tokenEndAction = "|"
const tokenGood = "^"
const tokenEvil = "~"
const tokenPending = "*"

// GameSettings holds a single game's information
type GameSettings struct {
	Name       string `json:"gameName"`
	Number     int    `json:"gameNumber"`
	RoundNum   int
	RoundNight bool
	Seed       int64
}

// Player holds a single player's name
type Player struct {
	Name    string      `json:"playerName"`
	Num     int         `json:"playerNumber"`
	Actions string      `json:"actions"`
	Role    RoleTracker `json:"role"`
}

// RoleTracker holds state info about a player's roles
type RoleTracker struct {
	Totals  map[int]int `json:"roles"`
	IsFixed bool        `json:"isFixed"`
	Fixed   int         `json:"fixed"`
}

// GameSetupSettings holds the game settings
type GameSetupSettings struct {
	Name      string         `json:"gameName"`
	Roles     map[string]int `json:"roles"`
	Total     int            `json:"totalPlayers"`
	Keep      int            `json:"keepPercent"`
	Universes uint64         `json:"universes"`
}

// MultiverseDefinition holds the state of all Universes
// Universes contains the active universe permutation numbers
// originalAssignments contains the permutation of the 0th universe
type MultiverseDefinition struct {
	Universes           []uint64
	originalAssignments []int
	rando               *rand.Rand
}

// Multiverse is the global definition of the multiverse
var Multiverse MultiverseDefinition

// Players is the global slice of Player instances
var Players []Player

// GameSetup is the global GameSetup settings
var GameSetup GameSetupSettings

// Game is the global GameSettings instance
var Game GameSettings

// ResetVars reinitializes all global variables
func ResetVars() {
	rand.Seed(time.Now().UTC().UnixNano())
	Players = nil
	GameSetup.Name = ""
	GameSetup.Roles = make(map[string]int)
	GameSetup.Total = 0
	GameSetup.Keep = 100
	Game.Name = ""
	Game.Number = -1
	Game.RoundNight = true
	Game.RoundNum = 0
	Game.Seed = rand.Int63()
	Multiverse.Universes = nil
	Multiverse.originalAssignments = nil
	Multiverse.rando = nil
	ResetObservations()
}

func getPlayerByName(playerName string) Player {
	for _, p := range Players {
		if p.Name == playerName {
			return p
		}
	}

	log.Printf("Attempted to get unknown player by name: %v", playerName)
	var unknownPlayer Player
	return unknownPlayer
}

func getPlayerByNumber(playerNumber int) Player {
	for _, p := range Players {
		if p.Num == playerNumber {
			return p
		}
	}

	log.Printf("Attempted to get unknown player by number: %d", playerNumber)
	var unknownPlayer Player
	return unknownPlayer
}

func getPlayerIndex(player Player) int {
	for i, p := range Players {
		if p.Num == player.Num {
			return i
		}
	}

	return -1
}

func playerCanPeek(player Player) bool {
	UpdateRoleTotals()
	for k, v := range player.Role.Totals {
		if v > 0 && roleTypes[k].CanPeek {
			return true
		}
	}

	return false
}

func playerCanAttack(player Player) bool {
	UpdateRoleTotals()
	for k, v := range player.Role.Totals {
		if v > 0 && roleTypes[k].CanAttack {
			return true
		}
	}

	return false
}

func playerRolePercent(player Player, roleID int) int {
	if len(Multiverse.Universes) == 0 {
		return 0
	}

	UpdateRoleTotals()

	amountRole := player.Role.Totals[roleID]
	totalUniverses := len(Multiverse.Universes)

	returnPercent := (amountRole * 100) / totalUniverses
	if returnPercent >= 100 && amountRole < totalUniverses {
		returnPercent = 99
	}
	if returnPercent <= 0 && amountRole > 0 {
		returnPercent = 1
	}
	return returnPercent
}

func playerEvilPercent(player Player) int {
	if len(Multiverse.Universes) == 0 {
		return 0
	}

	UpdateRoleTotals()

	amountEvil := 0
	for k, v := range player.Role.Totals {
		if roleTypes[k].Evil {
			amountEvil += v
		}
	}

	totalUniverses := len(Multiverse.Universes)
	returnPercent := (amountEvil * 100) / totalUniverses
	if returnPercent >= 100 && amountEvil < totalUniverses {
		returnPercent = 99
	}
	if returnPercent <= 0 && amountEvil > 0 {
		returnPercent = 1
	}
	return returnPercent
}

func playerGoodPercent(player Player) int {
	return 100 - playerEvilPercent(player)
}

func playerIsDead(player Player) bool {
	for _, o := range observations {
		if o.getSubject() == player.Num && o.getType() == "KillObservation" {
			return true
		}
	}

	return false
}

// First bool is true if a team won
// Second bool is true if winning team is evil
func checkWin() (bool, bool) {
	if len(Multiverse.Universes) == 0 {
		return false, false
	}

	goodPlayers := 0
	evilPlayers := 0
	goodMustKillPlayers := 0
	evilMustKillPlayers := 0

	UpdateRoleTotals()

	for _, p := range Players {
		if !playerIsDead(p) {
			amountEvil := 0
			amountGood := 0
			amountEvilMustKill := 0
			amountGoodMustKill := 0
			for k, v := range p.Role.Totals {
				mustKill := roleTypes[k].EnemyMustKill
				if roleTypes[k].Evil {
					amountEvil += v
					if mustKill {
						amountGoodMustKill += v
					}
				} else {
					amountGood += v
					if mustKill {
						amountEvilMustKill += v
					}
				}
			}

			if amountEvil > 0 && amountGood > 0 {
				return false, false
			}

			if amountGood > 0 {
				goodPlayers++
				if amountEvilMustKill > 0 {
					evilMustKillPlayers++
				}
			}

			if amountEvil > 0 {
				evilPlayers++
				if amountGoodMustKill > 0 {
					goodMustKillPlayers++
				}
			}
		}
	}

	if evilMustKillPlayers == 0 && goodPlayers <= evilPlayers {
		return true, true
	}
	if goodMustKillPlayers == 0 && goodPlayers >= evilPlayers {
		return true, false
	}
	return false, false
}

// PlayerDeadPercent returns the percentage of deadness for the input player
func PlayerDeadPercent(player Player) int {
	// Dead if marked dead
	if playerIsDead(player) {
		return 100
	}

	// Dead if lynched
	for _, o := range observations {
		if !o.getPending() && o.getSubject() == player.Num && o.getType() == "LynchObservation" {
			return 100
		}
	}

	// Figure out percentage from attacks
	attacksOnMe := make([]observation, 0, len(observations))
	totalUniverses := 0
	totalDeaths := 0

	for _, o := range observations {
		target, _ := o.getTarget()
		if !o.getPending() && o.getType() == "AttackObservation" && target == player.Num {
			attacksOnMe = append(attacksOnMe, o)
		}
	}

	for _, v := range Multiverse.Universes {
		totalUniverses++
		for _, o := range attacksOnMe {
			if AttackTarget(v, o.getSubject(), player.Num) {
				totalDeaths++
				break
			}
		}
	}

	returnPercent := (totalDeaths * 100) / totalUniverses
	if returnPercent >= 100 && totalDeaths < totalUniverses {
		returnPercent = 99
	}
	if returnPercent <= 0 && totalDeaths > 0 {
		returnPercent = 1
	}
	return returnPercent
}
