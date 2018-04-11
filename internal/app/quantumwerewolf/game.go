package quantumwerewolf

import (
	"net/http"
	"quantumwerewolf/pkg/quantumutilities"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

func showGame(c *gin.Context) {
	playersByName := make([]Player, gameSetup.Total, gameSetup.Total)
	playersByNum := make([]Player, gameSetup.Total, gameSetup.Total)
	for i, v := range players {
		playersByName[i] = v
		playersByNum[i] = v
	}
	sort.Slice(playersByName, func(i, j int) bool { return playersByName[i].Name < playersByNum[j].Name })
	sort.Slice(playersByNum, func(i, j int) bool { return playersByNum[i].Num < playersByNum[j].Num })
	var roundString = ""
	if game.RoundNight {
		roundString += "Night "
	} else {
		roundString += "Day "
	}
	roundString += strconv.Itoa(game.RoundNum)
	c.HTML(http.StatusOK, "game.gtpl", gin.H{
		"GameID":        game.Number,
		"Name":          gameSetup.Name,
		"TotalPlayers":  gameSetup.Total,
		"Roles":         gameSetup.Roles,
		"Round":         roundString,
		"IsNight":       game.RoundNight,
		"PlayersByName": playersByName,
		"PlayersByNum":  playersByNum,
	})
}

func rebuildGame(c *gin.Context, gameID int) {
	resetVars()

	gameQuery := "SELECT id, name, players, roles, keepPercent, round, nightPhase, randomSeed FROM game"
	gameQuery += " WHERE id=" + strconv.Itoa(gameID)
	gameQuery += " LIMIT 1"

	row, err := db.Query(gameQuery)
	if quantumutilities.HandleErr(c, err, "Error selecting game") {
		return
	}

	row.Next()
	rolesByteArray := make([]byte, 0, 100)
	err = row.Scan(&game.Number, &gameSetup.Name, &gameSetup.Total, &rolesByteArray, &gameSetup.Keep, &game.RoundNum, &game.RoundNight, &game.Seed)
	row.Close()
	if quantumutilities.HandleErr(c, err, "Error scanning game variables") {
		return
	}

	err = quantumutilities.GetInterface(rolesByteArray, &gameSetup.Roles)
	if quantumutilities.HandleErr(c, err, "Error getting game roles interface") {
		return
	}

	playerQuery := "SELECT name, num, actions FROM players"
	playerQuery += " WHERE gameid=" + strconv.Itoa(gameID)
	playerQuery += " LIMIT " + strconv.Itoa(gameSetup.Total)

	row, err = db.Query(playerQuery)
	if quantumutilities.HandleErr(c, err, "Error selecting players") {
		return
	}

	for row.Next() {
		var player Player
		err = row.Scan(&player.Name, &player.Num, &player.Actions)
		if quantumutilities.HandleErr(c, err, "Error scanning player variables") {
			return
		}
		players = append(players, player)
	}
	row.Close()

	createMultiverse()
}

func setGame(c *gin.Context) {
	err := c.Request.ParseForm()
	if quantumutilities.HandleErr(c, err, "Error setting gameSetup") {
		return
	}

	gameID, err := strconv.ParseInt(c.Query("gameId")[0:], 10, 32)

	rebuildGame(c, int(gameID))
}

func processActions(c *gin.Context) {
	err := c.Request.ParseForm()
	if quantumutilities.HandleErr(c, err, "Error processing actions") {
		return
	}

	var gameID = c.Request.FormValue("gameId")

	for _, p := range players {
		var attackSelection = c.Request.FormValue(p.Name + "Attack")
		var peekSelection = c.Request.FormValue(p.Name + "Peek")
		var lynchSelection = c.Request.FormValue(p.Name + "Lynch")
		if len(attackSelection) > 0 {
			p.Actions += strconv.Itoa(game.RoundNum) + tokenAttack + attackSelection + tokenEndAction
		}
		if len(peekSelection) > 0 {
			p.Actions += strconv.Itoa(game.RoundNum) + tokenPeek + peekSelection + tokenEndAction
		}
		if len(lynchSelection) > 0 {
			p.Actions += strconv.Itoa(game.RoundNum) + tokenLynch + lynchSelection + tokenEndAction
		}
		var dbStatement = "UPDATE players SET "
		dbStatement += "actions = "
		dbStatement += "'" + p.Actions + "'"
		dbStatement += " WHERE num=" + strconv.Itoa(p.Num) + " AND gameId=" + gameID
		quantumutilities.DbExec(c, db, dbStatement)
	}

	var advance = c.Request.Form["advance"]
	var advanceRound = false
	for _, s := range advance {
		if s == "true" {
			advanceRound = true
		}
	}

	if advanceRound {
		var nightBoolString = ""
		if game.RoundNight {
			game.RoundNum++
			game.RoundNight = false
			nightBoolString = "FALSE"
		} else {
			game.RoundNight = true
			nightBoolString = "TRUE"
		}

		var dbStatement = "UPDATE game SET "
		dbStatement += "round=" + strconv.Itoa(game.RoundNum)
		dbStatement += ", nightPhase=" + nightBoolString
		dbStatement += "WHERE id=" + gameID
		quantumutilities.DbExec(c, db, dbStatement)
	}

	setGame(c)
	showGame(c)
}
