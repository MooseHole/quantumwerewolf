package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Player has the player name and number
type Player struct {
	Number string `json:"number"`
	Name   string `json:"name"`
}

var players []Player

func getPlayerHandler(c *gin.Context) {
	//Convert the "players" variable to json
	playerListBytes, err := json.Marshal(players)

	// If there is an error, print it to the console, and return a server
	// error response to the user
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error getting player: %v", err))
		return
	}
	// If all goes well, write the JSON list of players to the response
	c.Writer.Write(playerListBytes)
}

func createPlayerHandler(c *gin.Context) {
	// Create a new instance of Player
	player := Player{}

	// We send all our data as HTML form data
	// the `ParseForm` method of the request, parses the
	// form values
	err := c.Request.ParseForm()

	// In case of any error, we respond with an error to the user
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error creating player: %v", err))
		return
	}

	// Get the information about the player from the form info
	player.Number = c.Request.Form.Get("number")
	player.Name = c.Request.Form.Get("name")

	// Append our existing list of players with a new entry
	players = append(players, player)

	//Finally, we redirect the user to the original HTMl page
	c.HTML(http.StatusOK, "players.gtpl", nil)
}
