package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Player holds a single player's name
type Game struct {
	Name   string `json:"gameName"`
	Number int    `json:"gameNumber"`
}

var games []Game

func getGamesHandler(c *gin.Context) {
	rows, err := db.Query("SELECT id, name FROM game")
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error selecting games: %v", err))
		return
	}

	for rows.Next() {
		var game Game
		err = rows.Scan(&game.Number, &game.Name)
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error scanning games: %v", err))
			return
		}

		games = append(games, game)
	}

	gameListBytes, err := json.Marshal(games)

	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error getting games: %v", err))
		return
	}

	c.Writer.Write(gameListBytes)
}
