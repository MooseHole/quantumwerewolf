package quantumwerewolf

import (
	"fmt"
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

	FillObservations()
	actionMessages := ""
	for _, o := range peekObservations {
		if o.Round == game.RoundNum-1 {
			resultString := "good"
			if o.IsEvil {
				resultString = "evil"
			}
			actionMessages += fmt.Sprintf("%s peeked at %s and found them %s.<br>", getPlayerByNumber(o.Subject).Name, getPlayerByNumber(o.Target).Name, resultString)
		}
	}
	for _, o := range attackObservations {
		if o.Round == game.RoundNum-1 {
			actionMessages += fmt.Sprintf("%s attacked %s.<br>", getPlayerByNumber(o.Subject).Name, getPlayerByNumber(o.Target).Name)
		}
	}
	for _, o := range lynchObservations {
		if o.Round == game.RoundNum {
			actionMessages += fmt.Sprintf("%s voted to lynch %s.<br>", getPlayerByNumber(o.Subject).Name, getPlayerByNumber(o.Target).Name)
		}
	}
	for _, o := range killObservations {
		if o.Round == game.RoundNum || (o.Round == game.RoundNum-1 && !game.RoundNight) {
			actionMessages += fmt.Sprintf("%s died and was a %s.<br>", getPlayerByNumber(o.Subject).Name, roleTypes[o.Role].Name)
		}
	}

	type ActionSelections struct {
		Peek     map[string]string
		Attack   map[string]string
		Lynch    map[string]string
		Peeked   map[string]string
		Attacked map[string]string
		Lynched  map[string]string
		Killed   map[string]string
	}
	actionSubjects := make(map[string]ActionSelections)
	FillObservations()

	for _, s := range players {
		var selection ActionSelections
		var playerIsDead = false

		selection.Peeked = make(map[string]string)
		selection.Attacked = make(map[string]string)
		selection.Lynched = make(map[string]string)
		selection.Killed = make(map[string]string)

		for _, o := range killObservations {
			if o.Subject == s.Num {
				selection.Killed[strconv.Itoa(o.Round)] = roleTypes[o.Role].Name
				playerIsDead = true
			}
		}

		for _, o := range peekObservations {
			if o.Subject == s.Num {
				var resultString = "good"
				if o.IsEvil {
					resultString = "evil"
				}
				selection.Peeked[strconv.Itoa(o.Round)] = getPlayerByNumber(o.Target).Name + "=" + resultString
			}
		}

		for _, o := range attackObservations {
			if o.Subject == s.Num {
				selection.Attacked[strconv.Itoa(o.Round)] = getPlayerByNumber(o.Target).Name
			}
		}

		for _, o := range lynchObservations {
			if o.Subject == s.Num {
				selection.Lynched[strconv.Itoa(o.Round)] = getPlayerByNumber(o.Target).Name
			}
		}

		// Set up next actions
		selection.Peek = make(map[string]string)
		selection.Attack = make(map[string]string)
		selection.Lynch = make(map[string]string)

		selection.Peek["--NONE--"] = ""
		selection.Attack["--NONE--"] = ""
		selection.Lynch["--NONE--"] = ""

		// Don't allow dead players to do actions
		if !playerIsDead {
			for _, t := range players {
				skipTarget := false

				// Don't add actions for dead targets
				for _, o := range killObservations {
					if o.Subject == t.Num {
						skipTarget = true
						break
					}
				}

				// Don't do actions on yourself
				if s.Num == t.Num {
					skipTarget = true
				}

				if skipTarget {
					continue
				}

				hasPeeked := false
				for _, o := range peekObservations {
					if o.Subject == s.Num && o.Target == t.Num {
						hasPeeked = true
						break
					}
				}
				if !hasPeeked {
					selection.Peek[t.Name] = t.Name
				}

				hasAttacked := false
				for _, o := range attackObservations {
					if o.Subject == s.Num && o.Target == t.Num {
						hasAttacked = true
						break
					}
				}
				if !hasAttacked {
					selection.Attack[t.Name] = t.Name
				}

				selection.Lynch[t.Name] = t.Name
			}
		}

		actionSubjects[s.Name] = selection
	}

	rounds := make([]string, game.RoundNum+1)
	for i := range rounds {
		rounds[i] = strconv.Itoa(i)
	}

	c.HTML(http.StatusOK, "game.gtpl", gin.H{
		"GameID":         game.Number,
		"Name":           gameSetup.Name,
		"TotalPlayers":   gameSetup.Total,
		"Roles":          gameSetup.Roles,
		"RoundNum":       strconv.Itoa(game.RoundNum),
		"Round":          roundString,
		"Rounds":         rounds,
		"IsNight":        game.RoundNight,
		"PlayersByName":  playersByName,
		"PlayersByNum":   playersByNum,
		"ActionMessages": actionMessages,
		"ActionSubjects": actionSubjects,
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
		player.Role.Totals = make(map[int]int)
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
			p.Actions += strconv.Itoa(game.RoundNum) + tokenAttack + strconv.Itoa(getPlayerByName(attackSelection).Num) + tokenEndAction
		}
		if len(peekSelection) > 0 {
			p.Actions += strconv.Itoa(game.RoundNum) + tokenPeek + strconv.Itoa(getPlayerByName(peekSelection).Num) + Peek(p.Num, getPlayerByName(peekSelection).Num) + tokenEndAction
		}
		if len(lynchSelection) > 0 {
			p.Actions += strconv.Itoa(game.RoundNum) + tokenLynch + strconv.Itoa(getPlayerByName(lynchSelection).Num) + tokenEndAction
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
		var lynchTargets = make(map[int]int)
		for _, p := range players {
			actionStrings := strings.Split(p.Actions, tokenEndAction)
			for _, a := range actionStrings {
				if len(a) >= lynchSubstringLength && a[0:lynchSubstringLength] == lynchSubstring {
					lynchTarget, _ := strconv.ParseInt(a[lynchSubstringLength:], 10, 32)
					lynchTargets[int(lynchTarget)]++
				}
			}
		}

		for t, n := range lynchTargets {
			if n > len(players)/2 {
				lynchedPlayer := getPlayerByNumber(t)

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
