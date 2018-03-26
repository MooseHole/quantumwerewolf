package quantumwerewolf

import (
	"math/rand"
	"time"
)

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
	Name string `json:"playerName"`
}

// GameSetup holds the game settings
type GameSetup struct {
	Name      string `json:"gameName"`
	Total     int    `json:"totalPlayers"`
	Villagers int    `json:"totalVillagers"`
	Seers     int    `json:"totalSeers"`
	Wolves    int    `json:"totalWolves"`
	Keep      int    `json:"keepPercent"`
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
	gameSetup.Total = 0
	gameSetup.Villagers = 0
	gameSetup.Seers = 0
	gameSetup.Wolves = 0
	gameSetup.Keep = 100
	game.Name = ""
	game.Number = -1
	game.RoundNight = true
	game.RoundNum = 0
	game.Seed = rand.Int63()
	multiverse.universes = nil
	multiverse.originalAssignments = nil
}
