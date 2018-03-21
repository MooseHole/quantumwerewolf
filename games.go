package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Game holds a single game's information
type Game struct {
	Name   string `json:"gameName"`
	Number int    `json:"gameNumber"`
}

var games []Game

func getGamesHandler(c *gin.Context) {
	log.Printf("Going to query: SELECT id, name FROM game")
	rows, err := db.Query("SELECT id, name FROM game")
	log.Printf("Query complete")
	if err != nil {
		log.Printf("Error selecting games: %v", err)
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error selecting games: %v", err))
		return
	}

	for rows.Next() {
		log.Printf("rows next")
		var game Game
		err = rows.Scan(&game.Number, &game.Name)
		if err != nil {
			log.Printf("Error scanning games: %v", err)
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error scanning games: %v", err))
			return
		}

		log.Printf("game.Number %d, game.Name %v", game.Number, game.Name)

		games = append(games, game)
	}

	gameListBytes, err := json.Marshal(games)

	if err != nil {
		log.Printf("Error getting games: %v", err)
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error getting games: %v", err))
		return
	}

	c.Writer.Write(gameListBytes)
}
