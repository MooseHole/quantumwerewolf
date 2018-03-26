package quantumwerewolf

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"quantumwerewolf/pkg/quantumutilities"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

func getRolesHandler(c *gin.Context) {
	//Convert the "gameSetup" variable to json
	roleListBytes, err := json.Marshal(gameSetup)

	// If there is an error, print it to the console, and return a server
	// error response to the user
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error getting gameSetup: %v", err))
		return
	}
	// If all goes well, write the JSON list of gameSetup to the response
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

	gameSetup.Total++
	if gameSetup.Total > 2 {
		gameSetup.Seers = 1
		gameSetup.Wolves = gameSetup.Total / 3
		gameSetup.Villagers = gameSetup.Total - gameSetup.Seers - gameSetup.Wolves
		amountAssigned := 0
		for _, v := range roleTypes {
			// ID 0 is reserved for Villagers
			if v.ID != 0 {
				assign := 0
				// If a whole number
				if v.DefaultAmount >= 1 {
					assign = int(v.DefaultAmount)
				} else {
					assign = int(float32(gameSetup.Total) * v.DefaultAmount)
				}
				gameSetup.Roles[v.Name] = assign
				amountAssigned += assign
			}
		}

		gameSetup.Roles[roleTypes[0].Name] = gameSetup.Total - amountAssigned
	}

	// Get the information about the player from the form info
	player.Name = c.Request.Form.Get("playerName")

	// Append our existing list of players with a new entry
	players = append(players, player)

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
			fmt.Sprintf("Error setting gameSetup: %v", err))
		return
	}

	gameSetup.Name = c.Request.FormValue("gameName")

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

	if int(s+w) <= gameSetup.Total {
		gameSetup.Seers = int(s)
		gameSetup.Wolves = int(w)
		gameSetup.Villagers = gameSetup.Total - gameSetup.Seers - gameSetup.Wolves
	}

	for _, v := range roleTypes {
		value, err := strconv.ParseInt(c.Request.FormValue(v.Name)[0:], 10, 64)
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error converting seers: %v", err))
		}
		gameSetup.Roles[v.Name] = int(value)
	}

	k, err := strconv.ParseInt(c.Request.FormValue("keep")[0:], 10, 64)
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error converting keep: %v", err))
	}
	gameSetup.Keep = int(k)

	createMultiverse()
	startGame(c)
	resetVars()
	c.HTML(http.StatusOK, "games.gtpl", nil)
}

func startGame(c *gin.Context) {
	dbStatement := ""

	dbStatement = "CREATE TABLE IF NOT EXISTS game ("
	dbStatement += "id SERIAL PRIMARY KEY"
	dbStatement += ", name text"
	dbStatement += ", players integer"
	dbStatement += ", seers integer"
	dbStatement += ", wolves integer"
	dbStatement += ", keepPercent integer"
	dbStatement += ", round integer"
	dbStatement += ", nightPhase boolean"
	dbStatement += ", randomSeed integer"
	dbStatement += ")"
	quantumutilities.DbExec(c, db, dbStatement)

	dbStatement = "CREATE TABLE IF NOT EXISTS players ("
	dbStatement += "id BIGSERIAL PRIMARY KEY"
	dbStatement += ", name text"
	dbStatement += ", num integer"
	dbStatement += ", gameid integer"
	dbStatement += ", actions text"
	dbStatement += ")"
	quantumutilities.DbExec(c, db, dbStatement)

	dbStatement = "INSERT INTO game ("
	dbStatement += "name, players, seers, wolves, keepPercent, round, nightPhase, randomSeed"
	dbStatement += ") VALUES ("
	dbStatement += "'" + gameSetup.Name + "'"
	dbStatement += ", " + strconv.Itoa(gameSetup.Total)
	dbStatement += ", " + strconv.Itoa(gameSetup.Seers)
	dbStatement += ", " + strconv.Itoa(gameSetup.Wolves)
	dbStatement += ", " + strconv.Itoa(gameSetup.Keep)
	dbStatement += ", " + strconv.Itoa(game.RoundNum)
	dbStatement += ", TRUE"
	dbStatement += ", " + strconv.Itoa(int(rand.Int31()))
	dbStatement += ") RETURNING id"
	var gameID = quantumutilities.DbExecReturn(c, db, dbStatement)

	// Assign random player numbers
	perm := rand.Perm(len(players))
	log.Printf("len(players) %d", len(players))
	for i, p := range players {
		dbStatement = "INSERT INTO players ("
		dbStatement += "name, num, gameid, actions"
		dbStatement += ") VALUES ("
		dbStatement += "'" + p.Name + "'"
		dbStatement += ", " + strconv.Itoa(perm[i])
		dbStatement += ", " + strconv.Itoa(gameID)
		dbStatement += ", ''"
		dbStatement += ")"
		log.Printf("Going to execute %q", dbStatement)
		quantumutilities.DbExec(c, db, dbStatement)
	}
}

func dropTables(c *gin.Context) {

	if _, err := db.Exec("DROP TABLE IF EXISTS game"); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error dropping game: %q", err))
		return
	}

	if _, err := db.Exec("DROP TABLE IF EXISTS players"); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error dropping player: %q", err))
		return
	}

	c.String(http.StatusOK, "Done dropping")
}
