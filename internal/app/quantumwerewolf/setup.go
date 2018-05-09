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
	playerListBytes, err := json.Marshal(Players)

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
	//Convert the "GameSetup" variable to json
	roleListBytes, err := json.Marshal(GameSetup)

	// If there is an error, print it to the console, and return a server
	// error response to the user
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error getting GameSetup: %v", err))
		return
	}
	// If all goes well, write the JSON list of GameSetup to the response
	c.Writer.Write(roleListBytes)
}

func createPlayerHandler(c *gin.Context) {
	setupRoles()

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

	GameSetup.Total++
	if GameSetup.Total > 2 {
		amountAssigned := 0
		for _, v := range roleTypes {
			// ID 0 is reserved for Villagers
			if v.ID != 0 {
				assign := 0
				// If a whole number
				if v.DefaultAmount >= 1 {
					assign = int(v.DefaultAmount)
				} else {
					assign = int(float32(GameSetup.Total) * v.DefaultAmount)
				}
				GameSetup.Roles[v.Name] = assign
				amountAssigned += assign
			}
		}

		GameSetup.Roles[roleTypes[0].Name] = GameSetup.Total - amountAssigned
	}

	// Get the information about the player from the form info
	player.Name = c.Request.Form.Get("playerName")

	// Append our existing list of players with a new entry
	Players = append(Players, player)

	//Finally, we redirect the user to the original HTMl page
	c.HTML(http.StatusOK, "playerSetup.gtpl", nil)
}

func setRolesHandler(c *gin.Context) {
	setupRoles()

	err := c.Request.ParseForm()
	if quantumutilities.HandleErr(c, err, "Error setting GameSetup") {
		return
	}

	GameSetup.Name = c.Request.FormValue("gameName")

	specialRoles := 0
	for _, v := range roleTypes {
		if v.ID != 0 {
			value, err := strconv.ParseInt(c.Request.FormValue(v.Name)[0:], 10, 64)
			if err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error converting seers: %v", err))
			}
			GameSetup.Roles[v.Name] = int(value)
			specialRoles += int(value)
		}
	}
	GameSetup.Roles[roleTypes[0].Name] = GameSetup.Total - specialRoles

	k, err := strconv.ParseInt(c.Request.FormValue("keep")[0:], 10, 64)
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error converting keep: %v", err))
	}
	GameSetup.Keep = int(k)

	u, err := strconv.ParseInt(c.Request.FormValue("universes")[0:], 10, 64)
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error converting universes: %v", err))
	}
	GameSetup.Universes = uint64(u)

	CreateGame(c)
	ResetVars()
	c.HTML(http.StatusOK, "gameList.gtpl", nil)
}

// CreateGame sets up and stores the initial game parameters in the database
func CreateGame(c *gin.Context) {
	dbStatement := ""

	dbStatement = "CREATE TABLE IF NOT EXISTS game ("
	dbStatement += "id SERIAL PRIMARY KEY"
	dbStatement += ", name text"
	dbStatement += ", players integer"
	dbStatement += ", roles bytea"
	dbStatement += ", universes integer"
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

	roleBlob, err := quantumutilities.GetBytes(GameSetup.Roles)
	quantumutilities.HandleErr(c, err, "Error getting Roles as bytes")
	roleBytesString := fmt.Sprintf("'\\x%x'", roleBlob)
	dbStatement = "INSERT INTO game ("
	dbStatement += "name, players, roles, universes, round, nightPhase, randomSeed"
	dbStatement += ") VALUES ("
	dbStatement += "'" + GameSetup.Name + "'"
	dbStatement += ", " + strconv.Itoa(GameSetup.Total)
	dbStatement += ", " + roleBytesString
	dbStatement += ", " + strconv.FormatUint(GameSetup.Universes, 10)
	dbStatement += ", " + strconv.Itoa(Game.RoundNum)
	dbStatement += ", TRUE"
	dbStatement += ", " + strconv.Itoa(int(rand.Int31()))
	dbStatement += ") RETURNING id"
	var gameID = quantumutilities.DbExecReturn(c, db, dbStatement)

	// Assign random player numbers
	perm := rand.Perm(len(Players))
	log.Printf("len(players) %d", len(Players))
	for i, p := range Players {
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
