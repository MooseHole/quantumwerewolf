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
const tokenLynch = "&"
const tokenKilled = "#"
const tokenEndAction = "|"
const tokenGood = "^"
const tokenEvil = "~"
const tokenPending = "*"

// Game holds a single game's information
type Game struct {
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

// GameSetup holds the game settings
type GameSetup struct {
	Name  string         `json:"gameName"`
	Roles map[string]int `json:"roles"`
	Total int            `json:"totalPlayers"`
	Keep  int            `json:"keepPercent"`
}

// Multiverse holds the state of all universes
// universes contains the active universe permutation numbers
// originalAssignments contains the permutation of the 0th universe
type Multiverse struct {
	universes           []uint64
	originalAssignments []int
}

var multiverse Multiverse
var players []Player
var gameSetup GameSetup
var game Game

// ResetVars reinitializes all global variables
func ResetVars() {
	rand.Seed(time.Now().UTC().UnixNano())
	players = nil
	gameSetup.Name = ""
	gameSetup.Roles = make(map[string]int)
	gameSetup.Total = 0
	gameSetup.Keep = 100
	game.Name = ""
	game.Number = -1
	game.RoundNight = true
	game.RoundNum = 0
	game.Seed = rand.Int63()
	multiverse.universes = nil
	multiverse.originalAssignments = nil
	ResetObservations()
}

func getPlayerByName(playerName string) Player {
	for _, p := range players {
		if p.Name == playerName {
			return p
		}
	}

	log.Printf("Attempted to get unknown player by name: %v", playerName)
	var unknownPlayer Player
	return unknownPlayer
}

func getPlayerByNumber(playerNumber int) Player {
	for _, p := range players {
		if p.Num == playerNumber {
			return p
		}
	}

	log.Printf("Attempted to get unknown player by number: %d", playerNumber)
	var unknownPlayer Player
	return unknownPlayer
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

func playerEvilPercent(player Player) int {
	UpdateRoleTotals()

	amountEvil := 0
	for k, v := range player.Role.Totals {
		if roleTypes[k].Evil {
			amountEvil += v
		}
	}

	totalUniverses := len(multiverse.universes)
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

func playerDeadPercent(player Player) int {
	attacksOnMe := make([]AttackObservation, 0, len(attackObservations))
	totalUniverses := 0
	totalDeaths := 0

	for _, o := range attackObservations {
		if !o.Pending && o.Target == player.Num {
			attacksOnMe = append(attacksOnMe, o)
		}
	}

	for _, v := range multiverse.universes {
		totalUniverses++
		for _, o := range attacksOnMe {
			if AttackTarget(v, o.Subject, player.Num) {
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
