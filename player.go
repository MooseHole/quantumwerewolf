package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Player holds a single player's name
type Player struct {
	Number int    `json:"number"`
	Name   string `json:"playerName"`
}

// Roles holds the role settings
type Roles struct {
	Total     int `json:"totalPlayers"`
	Villagers int `json:"totalVillagers"`
	Seers     int `json:"totalSeers"`
	Wolves    int `json:"totalWolves"`
}

var players []Player
var roles Roles

func getPlayerHandler(c *gin.Context) {
	// If first time
	if len(players) == 0 {
		roles.Total = 0
		roles.Villagers = 0
		roles.Seers = 0
		roles.Wolves = 0
		c.HTML(http.StatusOK, "players.gtpl", nil)
		return
	}

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

func getRolesHandler(c *gin.Context) {
	//Convert the "roles" variable to json
	roleListBytes, err := json.Marshal(roles)

	// If there is an error, print it to the console, and return a server
	// error response to the user
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error getting roles: %v", err))
		return
	}
	// If all goes well, write the JSON list of roles to the response
	c.Writer.Write(roleListBytes)
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

	roles.Total++
	if roles.Total > 2 {
		roles.Seers = 1
		roles.Wolves = roles.Total / 3
		roles.Villagers = roles.Total - roles.Seers - roles.Wolves
	}

	// Get the information about the player from the form info
	player.Number = -1
	player.Name = c.Request.Form.Get("playerName")

	// Append our existing list of players with a new entry
	players = append(players, player)

	perm := rand.Perm(roles.Total)
	i := 0
	for _, p := range players {
		p.Number = perm[i]
		//		c.String(http.StatusOK, strconv.Itoa(i)+" "+strconv.Itoa(perm[i])+" / "+strconv.Itoa(roles.Total))
	}

	//	c.String(http.StatusOK, players[0].Number)
	//Finally, we redirect the user to the original HTMl page
	c.HTML(http.StatusOK, "players.gtpl", nil)
}

func setRolesHandler(c *gin.Context) {
	// We send all our data as HTML form data
	// the `ParseForm` method of the request, parses the
	// form values
	err := c.Request.ParseForm()

	// In case of any error, we respond with an error to the user
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error setting roles: %v", err))
		return
	}

	s, err := strconv.ParseInt(c.Request.FormValue("seers")[0:], 10, 64)
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error converting seers: %v", err))
	}

	w, err := strconv.ParseInt(c.Request.FormValue("wolves")[0:], 10, 64)
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error converting wolves: %v", err))
	}

	if int(s+w) <= roles.Total {
		roles.Seers = int(s)
		roles.Wolves = int(w)
		roles.Villagers = roles.Total - roles.Seers - roles.Wolves
	}

	//	c.String(http.StatusOK, players[0].Number)
	//Finally, we redirect the user to the original HTMl page
	c.HTML(http.StatusOK, "players.gtpl", nil)
}
