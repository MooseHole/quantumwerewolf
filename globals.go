package main

import (
	"math/rand"
	"time"

	_ "github.com/lib/pq"
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

// Roles holds the role settings
type Roles struct {
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
var roles Roles
var game Game

func resetVars() {
	rand.Seed(time.Now().UTC().UnixNano())
	players = nil
	roles.Name = ""
	roles.Total = 0
	roles.Villagers = 0
	roles.Seers = 0
	roles.Wolves = 0
	roles.Keep = 100
	game.Name = ""
	game.Number = -1
	game.RoundNight = true
	game.RoundNum = 0
	game.Seed = rand.Int63()
	multiverse.universes = nil
	multiverse.originalAssignments = nil
}
