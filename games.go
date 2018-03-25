package main

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
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

var game Game

func getGamesHandler(c *gin.Context) {
	rows, err := db.Query("SELECT id, name FROM game")
	if handleErr(c, err, "Error selecting games") {
		return
	}

	var games []Game
	for rows.Next() {
		var thisGame Game
		err = rows.Scan(&thisGame.Number, &thisGame.Name)
		if handleErr(c, err, "Error scanning games") {
			return
		}

		games = append(games, thisGame)
	}

	gameListBytes, err := json.Marshal(games)

	if handleErr(c, err, "Error getting games") {
		return
	}

	c.Writer.Write(gameListBytes)
}
