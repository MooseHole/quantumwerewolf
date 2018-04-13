package quantumwerewolf

import (
	"net/http"
	"quantumwerewolf/pkg/quantumutilities"
	"sort"
	"strconv"
	"strings"

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
	ResetVars()

	gameQuery := "SELECT id, name, players, roles, keepPercent, round, nightPhase, randomSeed FROM game"
	gameQuery += " WHERE id=" + strconv.Itoa(gameID)
	gameQuery += " LIMIT 1"

	row, err := db.Query(gameQuery)
	if quantumutilities.HandleErr(c, err, "Error selecting game ["+gameQuery+"]") {
		return
	}

	if row.Next() {
		rolesByteArray := make([]byte, 0, 100)
		err = row.Scan(&game.Number, &gameSetup.Name, &gameSetup.Total, &rolesByteArray, &gameSetup.Keep, &game.RoundNum, &game.RoundNight, &game.Seed)
		row.Close()

		if quantumutilities.HandleErr(c, err, "Error scanning game variables ["+gameQuery+"]") {
			return
		}

		err = quantumutilities.GetInterface(rolesByteArray, &gameSetup.Roles)
		if quantumutilities.HandleErr(c, err, "Error getting game roles interface ["+gameQuery+"]") {
			return
		}
	}
	row.Close()

	playerQuery := "SELECT name, num, actions FROM players"
	playerQuery += " WHERE gameid=" + strconv.Itoa(gameID)
	playerQuery += " LIMIT " + strconv.Itoa(gameSetup.Total)

	row, err = db.Query(playerQuery)
	if quantumutilities.HandleErr(c, err, "Error selecting players ["+playerQuery+"]") {
		return
	}

	players = nil
	for row.Next() {
		var player Player
		err = row.Scan(&player.Name, &player.Num, &player.Actions)
		if quantumutilities.HandleErr(c, err, "Error scanning player variables ["+playerQuery+"]") {
			return
		}
		player.Role = make(map[int]int)
		players = append(players, player)
	}
	row.Close()

	CreateMultiverse()
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
	gameIDNum, err := strconv.ParseInt(gameID, 10, 32)

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
		rebuildGame(c, int(gameIDNum))
		lynchSubstring := strconv.Itoa(game.RoundNum) + tokenLynch
		lynchSubstringLength := len(lynchSubstring)
		var lynchTargets = make(map[string]int)
		for _, p := range players {
			actionStrings := strings.Split(p.Actions, tokenEndAction)
			for _, a := range actionStrings {
				if len(a) >= lynchSubstringLength && a[0:lynchSubstringLength] == lynchSubstring {
					lynchTarget := a[lynchSubstringLength:]
					lynchTargets[lynchTarget]++
				}
			}
		}

		for t, n := range lynchTargets {
			if n > len(players)/2 {
				lynchedPlayer := getPlayerByName(t)

				fixedRole := collapseToFixedRole(lynchedPlayer.Num)

				var dbStatement = "UPDATE players SET "
				dbStatement += "actions = "
				dbStatement += "'" + lynchedPlayer.Actions + strconv.Itoa(game.RoundNum) + tokenKilled + strconv.Itoa(fixedRole) + tokenEndAction + "'"
				dbStatement += " WHERE num=" + strconv.Itoa(lynchedPlayer.Num) + " AND gameId=" + gameID
				quantumutilities.DbExec(c, db, dbStatement)

				break
			}
		}

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
		dbStatement += " WHERE id=" + gameID
		quantumutilities.DbExec(c, db, dbStatement)
	}

	rebuildGame(c, int(gameIDNum))
	showGame(c)
}
