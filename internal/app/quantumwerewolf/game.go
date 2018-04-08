package quantumwerewolf

import (
	"quantumwerewolf/pkg/quantumutilities"
	"strconv"

	"github.com/gin-gonic/gin"
)

func rebuildGame(c *gin.Context, gameID int) {
	resetVars()

	gameQuery := "SELECT name, players, roles, keepPercent, round, nightPhase, randomSeed FROM game"
	gameQuery += " WHERE id=" + strconv.Itoa(gameID)
	gameQuery += " LIMIT 1"

	row, err := db.Query(gameQuery)
	if quantumutilities.HandleErr(c, err, "Error selecting game") {
		return
	}

	row.Next()
	rolesByteArray := make([]byte, 0, 100)
	err = row.Scan(&gameSetup.Name, &gameSetup.Total, &rolesByteArray, &gameSetup.Keep, &game.RoundNum, &game.RoundNight, &game.Seed)
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

	var advance = c.Request.FormValue("advance")

	for _, p := range players {
		var attackSelection = p.Name + "Attack"
		var peekSelection = p.Name + "Peek"
		var lynchSelection = p.Name + "Lynch"
		if len(attackSelection) > 0 {
			p.Actions += "A:" + c.Request.FormValue(attackSelection) + "|"
		}
		if len(peekSelection) > 0 {
			p.Actions += "P:" + c.Request.FormValue(peekSelection) + "|"
		}
		if len(lynchSelection) > 0 {
			p.Actions += "L:" + c.Request.FormValue(lynchSelection) + "|"
		}
		var dbStatement = "UPDATE players SET "
		dbStatement += "actions = "
		dbStatement += "'" + p.Actions + "'"
		dbStatement += " WHERE num=" + strconv.Itoa(p.Num) + " AND gameId=" + strconv.Itoa(game.Number)
		quantumutilities.DbExec(c, db, dbStatement)
	}

	if advance == "advance" {
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
		dbStatement += "WHERE id=" + strconv.Itoa(game.Number)
	}
}
