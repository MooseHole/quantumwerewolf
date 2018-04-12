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
	Role    map[int]int `json:"role"`
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

func resetVars() {
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
}

func getPlayerByName(playerName string) Player {
	for _, p := range players {
		if p.Name == playerName {
			return p
		}
	}

	log.Printf("Attempted to get unknown player: %v", playerName)
	var unknownPlayer Player
	return unknownPlayer
}
