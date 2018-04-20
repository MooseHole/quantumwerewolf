package quantumwerewolf

import (
	"fmt"
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

	FillObservations()
	actionMessages := ""
	for _, o := range peekObservations {
		if !o.Pending && o.Round == game.RoundNum-1 {
			resultString := "good"
			if o.IsEvil {
				resultString = "evil"
			}
			actionMessages += fmt.Sprintf("%s peeked at %s and found them %s.<br>", getPlayerByNumber(o.Subject).Name, getPlayerByNumber(o.Target).Name, resultString)
		}
	}
	for _, o := range attackObservations {
		if !o.Pending && o.Round == game.RoundNum-1 {
			actionMessages += fmt.Sprintf("%s attacked %s.<br>", getPlayerByNumber(o.Subject).Name, getPlayerByNumber(o.Target).Name)
		}
	}
	for _, o := range lynchObservations {
		if !o.Pending && o.Round == game.RoundNum {
			actionMessages += fmt.Sprintf("%s voted to lynch %s.<br>", getPlayerByNumber(o.Subject).Name, getPlayerByNumber(o.Target).Name)
		}
	}
	for _, o := range killObservations {
		if !o.Pending && o.Round == game.RoundNum || (o.Round == game.RoundNum-1 && !game.RoundNight) {
			actionMessages += fmt.Sprintf("%s died and was a %s.<br>", getPlayerByNumber(o.Subject).Name, roleTypes[o.Role].Name)
		}
	}

	type ActionSelections struct {
		Peek        map[string]string
		Attack      map[string]string
		Lynch       map[string]string
		Peeked      map[string]string
		PeekResult  map[string]string
		Attacked    map[string]string
		Lynched     map[string]string
		Killed      map[string]string
		GoodPercent int
		EvilPercent int
		DeadPercent int
	}
	actionSubjects := make(map[string]ActionSelections)
	FillObservations()

	for _, s := range players {
		var selection ActionSelections
		var playerIsDead = false

		selection.GoodPercent = playerGoodPercent(s)
		selection.EvilPercent = playerEvilPercent(s)
		selection.DeadPercent = playerDeadPercent(s)
		selection.Peeked = make(map[string]string)
		selection.PeekResult = make(map[string]string)
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
				var resultString = "=good"
				if o.IsEvil {
					resultString = "=evil"
				}
				selection.Peeked[strconv.Itoa(o.Round)] = getPlayerByNumber(o.Target).Name
				selection.PeekResult[strconv.Itoa(o.Round)] = resultString
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
					if !o.Pending && o.Subject == t.Num {
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

				if playerCanPeek(s) {
					hasPeeked := false
					for _, o := range peekObservations {
						if !o.Pending && o.Subject == s.Num && o.Target == t.Num {
							hasPeeked = true
							break
						}
					}
					if !hasPeeked {
						selection.Peek[t.Name] = t.Name
					}
				}

				if playerCanAttack(s) {
					hasAttacked := false
					for _, o := range attackObservations {
						if !o.Pending && o.Subject == s.Num && o.Target == t.Num {
							hasAttacked = true
							break
						}
					}
					if !hasAttacked {
						selection.Attack[t.Name] = t.Name
					}
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

	universes := make(map[int]string)
	for _, u := range multiverse.universes {
		universes[int(u)] = getUniverseString(u)
	}

	c.HTML(http.StatusOK, "game.gtpl", gin.H{
		"GameID":         game.Number,
		"Name":           gameSetup.Name,
		"TotalPlayers":   gameSetup.Total,
		"Roles":          gameSetup.Roles,
		"RoundNum":       strconv.Itoa(game.RoundNum),
		"Round":          roundString,
		"Rounds":         rounds,
		"Universes":      universes,
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
			var observation AttackObservation
			observation.Pending = true
			observation.Round = game.RoundNum
			observation.Subject = p.Num
			observation.Target = getPlayerByName(attackSelection).Num
			addAttackObservation(observation)
		}
		if len(peekSelection) > 0 {
			var observation PeekObservation
			observation.Pending = true
			observation.Round = game.RoundNum
			observation.Subject = p.Num
			observation.Target = getPlayerByName(peekSelection).Num
			observation.IsEvil = false // Determined at commit time
			addPeekObservation(observation)
		}
		if len(lynchSelection) > 0 {
			var observation LynchObservation
			observation.Pending = true
			observation.Round = game.RoundNum
			observation.Subject = p.Num
			observation.Target = getPlayerByName(lynchSelection).Num
			addLynchObservation(observation)
		}
	}

	var advance = c.Request.Form["advance"]
	var advanceRound = false
	for _, s := range advance {
		if s == "true" {
			advanceRound = true
		}
	}

	if advanceRound {
		var lynchTargets = make(map[int]int)
		for _, o := range lynchObservations {
			if game.RoundNum == o.Round {
				lynchTargets[o.Target]++
			}
		}

		for t, n := range lynchTargets {
			if n > len(players)/2 {
				lynchedPlayer := getPlayerByNumber(t)

				var observation KillObservation
				observation.Pending = true
				observation.Round = game.RoundNum
				observation.Subject = lynchedPlayer.Num
				observation.Role = collapseToFixedRole(lynchedPlayer.Num)
				addKillObservation(observation)
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

		CommitObservations()

		var dbStatement = "UPDATE game SET "
		dbStatement += "round=" + strconv.Itoa(game.RoundNum)
		dbStatement += ", nightPhase=" + nightBoolString
		dbStatement += " WHERE id=" + gameID
		quantumutilities.DbExec(c, db, dbStatement)
	}

	FillActionsWithObservations()
	for _, p := range players {
		var dbStatement = "UPDATE players SET "
		dbStatement += "actions = '" + p.Actions + "' WHERE num=" + strconv.Itoa(p.Num) + " AND gameId=" + gameID
		quantumutilities.DbExec(c, db, dbStatement)
	}

	rebuildGame(c, int(gameIDNum))
	showGame(c)
}
