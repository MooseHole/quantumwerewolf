package quantumwerewolf

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"quantumwerewolf/pkg/quantumutilities"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

func showGame(c *gin.Context) {
	playersByName := make([]Player, GameSetup.Total, GameSetup.Total)
	playersByNum := make([]Player, GameSetup.Total, GameSetup.Total)
	for i, v := range Players {
		playersByName[i] = v
		playersByNum[i] = v
	}
	sort.Slice(playersByName, func(i, j int) bool { return playersByName[i].Name < playersByNum[j].Name })
	sort.Slice(playersByNum, func(i, j int) bool { return playersByNum[i].Num < playersByNum[j].Num })
	var roundString = ""
	if Game.RoundNight {
		roundString += "Night "
	} else {
		roundString += "Day "
	}
	roundString += strconv.Itoa(Game.RoundNum)

	FillObservations()
	actionMessages := make(map[int][]string)

	for _, p := range Players {
		actionMessages[p.Num] = make([]string, 0)
		actionMessages[p.Num] = append(actionMessages[p.Num], p.Name)
		actionMessages[p.Num] = append(actionMessages[p.Num], "["+GameSetup.Name+"] "+roundString)
		actionMessages[p.Num] = append(actionMessages[p.Num], "You are player number "+strconv.Itoa(p.Num))
		for round := 0; round < Game.RoundNum; round++ {
			for _, o := range observations {
				if !o.getPending() && o.getRound() == round && o.getSubject() == p.Num {
					actionMessages[p.Num] = append(actionMessages[p.Num], o.actionMessage())
				}
			}
		}
	}

	type ActionSelections struct {
		Name       string
		RevealName string
		RevealRole string
		Peek       map[string]string
		Attack     map[string]string
		Vote       map[string]string
		Peeked     map[string]string
		PeekResult map[string]string
		Attacked   map[string]string
		Voted      map[string]string
		Killed     map[string]string
		Lynched    map[string]string
		Percents   map[string]int
	}
	actionSubjects := make(map[int]ActionSelections)
	FillObservations()

	winner, evilWins := checkWin()
	winMessage := ""
	if winner {
		if evilWins {
			winMessage = "EVIL WINS!"
		} else {
			winMessage = "GOOD WINS!"
		}
	}

	for _, s := range playersByNum {
		var selection ActionSelections
		var playerIsDead = false

		selection.Percents = make(map[string]int)
		selection.Peeked = make(map[string]string)
		selection.PeekResult = make(map[string]string)
		selection.Attacked = make(map[string]string)
		selection.Voted = make(map[string]string)
		selection.Killed = make(map[string]string)

		selection.Name = s.Name
		selection.RevealName = "--Hidden--"
		selection.RevealRole = "--Undetermined--"
		selection.Percents["Good"] = playerGoodPercent(s)
		selection.Percents["Evil"] = playerEvilPercent(s)
		selection.Percents["Dead"] = PlayerDeadPercent(s)
		for _, v := range roleTypes {
			selection.Percents[v.Name] = playerRolePercent(s, v.ID)
		}

		for _, o := range observations {
			if o.dead(s.Num) {
				playerIsDead = true
				break
			}
		}

		for _, o := range observations {
			if o.getSubject() == s.Num && o.getType() == "KillObservation" {
				role, err := o.getRole()
				if err == nil {
					selection.Killed[strconv.Itoa(o.getRound())] = roleTypes[role].Name
					break
				}
			}
		}

		if s.Role.IsFixed {
			selection.RevealRole = "--Resolved--"
		}

		if winner || playerIsDead {
			selection.RevealName = s.Name
			if s.Role.IsFixed {
				selection.RevealRole = roleTypes[s.Role.Fixed].Name
			}
		}

		for _, o := range observations {
			if o.getSubject() == s.Num {
				target, err := o.getTarget()
				if err == nil {
					switch observationType := o.getType(); observationType {
					case "PeekObservation":
						var resultString = "=good"
						isEvil, err := o.getIsEvil()
						if isEvil {
							resultString = "=evil"
						}
						if err != nil {
							resultString = "=??"
						}
						selection.Peeked[strconv.Itoa(o.getRound())] = getPlayerByNumber(target).Name
						selection.PeekResult[strconv.Itoa(o.getRound())] = resultString
					case "AttackObservation":
						selection.Attacked[strconv.Itoa(o.getRound())] = getPlayerByNumber(target).Name
					case "VoteObservation":
						selection.Voted[strconv.Itoa(o.getRound())] = getPlayerByNumber(target).Name
					}
				}
			}
		}

		// Set up next actions
		selection.Peek = make(map[string]string)
		selection.Attack = make(map[string]string)
		selection.Vote = make(map[string]string)

		selection.Peek["--NONE--"] = ""
		selection.Attack["--NONE--"] = ""
		selection.Vote["--NONE--"] = ""

		if !playerIsDead {
			// Don't allow dead players to do actions
			for _, t := range Players {
				skipTarget := false

				// Don't add actions for dead targets
				for _, o := range observations {
					if o.dead(t.Num) {
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
					for _, o := range observations {
						target, err := o.getTarget()
						if err == nil {
							if !o.getPending() && o.getSubject() == s.Num && target == t.Num {
								hasPeeked = true
								break
							}
						}
					}
					if !hasPeeked {
						selection.Peek[t.Name] = t.Name
					}
				}

				if playerCanAttack(s) {
					hasAttacked := false
					for _, o := range observations {
						target, err := o.getTarget()
						if err == nil {
							if !o.getPending() && o.getSubject() == s.Num && target == t.Num {
								hasAttacked = true
								break
							}
						}
					}
					if !hasAttacked {
						selection.Attack[t.Name] = t.Name
					}
				}

				selection.Vote[t.Name] = t.Name
			}
		}

		actionSubjects[s.Num] = selection
	}

	if Game.RoundNight {
		for k, v := range actionSubjects {
			if len(v.Attack) > 1 {
				actionMessages[k] = append(actionMessages[k], "You may attack one of the following:")
				for _, a := range v.Attack {
					if len(a) > 0 {
						actionMessages[k] = append(actionMessages[k], a)
					}
				}
			}
			if len(v.Peek) > 1 {
				actionMessages[k] = append(actionMessages[k], "You may peek at one of the following:")
				for _, a := range v.Peek {
					if len(a) > 0 {
						actionMessages[k] = append(actionMessages[k], a)
					}
				}
			}
		}
	}

	rounds := make([]string, Game.RoundNum+1)
	for i := range rounds {
		rounds[i] = strconv.Itoa(i)
	}

	universes := make(map[int]string)
	for _, u := range Multiverse.Universes {
		universes[int(u)] = getUniverseString(u)
	}

	c.HTML(http.StatusOK, "game.gtpl", gin.H{
		"GameID":         Game.Number,
		"Name":           GameSetup.Name,
		"TotalPlayers":   GameSetup.Total,
		"Roles":          GameSetup.Roles,
		"RoundNum":       strconv.Itoa(Game.RoundNum),
		"Round":          roundString,
		"Rounds":         rounds,
		"Universes":      universes,
		"IsNight":        Game.RoundNight,
		"PlayersByName":  playersByName,
		"PlayersByNum":   playersByNum,
		"ActionMessages": actionMessages,
		"ActionSubjects": actionSubjects,
		"WinMessage":     winMessage,
		"Graph":          multiverseProgression(c, Game.Number),
	})
}

func rebuildGame(c *gin.Context, gameID int) {
	ResetVars()

	gameQuery := "SELECT id, name, players, roles, universes, round, nightPhase, randomSeed FROM game"
	gameQuery += " WHERE id=" + strconv.Itoa(gameID)
	gameQuery += " LIMIT 1"

	row, err := db.Query(gameQuery)
	if quantumutilities.HandleErr(c, err, "Error selecting game ["+gameQuery+"]") {
		return
	}

	if row.Next() {
		rolesByteArray := make([]byte, 0, 100)
		err = row.Scan(&Game.Number, &GameSetup.Name, &GameSetup.Total, &rolesByteArray, &GameSetup.Universes, &Game.RoundNum, &Game.RoundNight, &Game.Seed)
		row.Close()

		if quantumutilities.HandleErr(c, err, "Error scanning game variables ["+gameQuery+"]") {
			return
		}

		err = quantumutilities.GetInterface(rolesByteArray, &GameSetup.Roles)
		if quantumutilities.HandleErr(c, err, "Error getting game roles interface ["+gameQuery+"]") {
			return
		}
	}
	row.Close()

	playerQuery := "SELECT name, num, actions FROM players"
	playerQuery += " WHERE gameid=" + strconv.Itoa(gameID)
	playerQuery += " LIMIT " + strconv.Itoa(GameSetup.Total)

	row, err = db.Query(playerQuery)
	if quantumutilities.HandleErr(c, err, "Error selecting players ["+playerQuery+"]") {
		return
	}

	Players = nil
	for row.Next() {
		var player Player
		err = row.Scan(&player.Name, &player.Num, &player.Actions)
		if quantumutilities.HandleErr(c, err, "Error scanning player variables ["+playerQuery+"]") {
			return
		}
		player.Role.Totals = make(map[int]int)
		Players = append(Players, player)
	}
	row.Close()

	CreateMultiverse()
}

func multiverseProgression(c *gin.Context, gameID int) string {
	rebuildGame(c, gameID)
	progression := make(map[int][]uint64)
	deadAmount := make(map[int][]int)
	progression[-1] = append(progression[-1], Multiverse.Universes...)
	for i := 0; i < len(Players); i++ {
		deadAmount[-1] = append(deadAmount[-1], PlayerDeadPercent(getPlayerByNumber(i)))
	}
	for round := 0; round <= Game.RoundNum; round++ {
		CollapseAllUpTo(round)
		for i := 0; i < len(Players); i++ {
			deadAmount[round] = append(deadAmount[round], PlayerDeadPercent(getPlayerByNumber(i)))
		}
		progression[round] = append(progression[round], Multiverse.Universes...)
	}

	img := image.NewRGBA(image.Rect(0, 0, Game.RoundNum+1, len(Players)))
	for round := -1; round <= Game.RoundNum; round++ {
		totalRed := make([]int, len(Players))
		totalGreen := make([]int, len(Players))
		totalBlue := make([]int, len(Players))
		for _, universeNum := range progression[round] {
			universe := quantumutilities.Kthperm(Multiverse.originalAssignments, universeNum)
			for i, role := range universe {
				totalRed[i] += int(roleTypes[role].Color.R)
				totalGreen[i] += int(roleTypes[role].Color.G)
				totalBlue[i] += int(roleTypes[role].Color.B)
			}
		}
		for i := 0; i < len(Players); i++ {
			myColor := color.RGBA{uint8((totalRed[i] / len(progression[round])) * (100 - deadAmount[round][i]/2) / 100), uint8((totalGreen[i] / len(progression[round])) * (100 - deadAmount[round][i]/2) / 100), uint8((totalBlue[i] / len(progression[round])) * (100 - deadAmount[round][i]/2) / 100), 255}
			img.Set(round+1, i, myColor)
		}
	}

	out := new(bytes.Buffer)
	err := png.Encode(out, img)
	quantumutilities.HandleErr(c, err, "Error processing multiverse graph")

	base64Img := base64.StdEncoding.EncodeToString(out.Bytes())
	return base64Img
}

func setGame(c *gin.Context) {
	err := c.Request.ParseForm()
	if quantumutilities.HandleErr(c, err, "Error setting GameSetup") {
		return
	}

	gameID, err := strconv.ParseInt(c.Query("gameId")[0:], 10, 32)

	rebuildGame(c, int(gameID))
	CollapseAll()
}

func processActions(c *gin.Context) {
	err := c.Request.ParseForm()
	if quantumutilities.HandleErr(c, err, "Error processing actions") {
		return
	}

	var gameID = c.Request.FormValue("gameId")
	gameIDNum, err := strconv.ParseInt(gameID, 10, 32)

	for _, p := range Players {
		var attackSelection = c.Request.FormValue(strconv.Itoa(p.Num) + "Attack")
		var peekSelection = c.Request.FormValue(strconv.Itoa(p.Num) + "Peek")
		var voteSelection = c.Request.FormValue(strconv.Itoa(p.Num) + "Vote")
		if len(attackSelection) > 0 {
			var observation AttackObservation
			observation.Pending = true
			observation.Round = Game.RoundNum
			observation.Subject = p.Num
			observation.Target = getPlayerByName(attackSelection).Num
			observation.add()
		}
		if len(peekSelection) > 0 {
			var observation PeekObservation
			observation.Pending = true
			observation.Round = Game.RoundNum
			observation.Subject = p.Num
			observation.Target = getPlayerByName(peekSelection).Num
			observation.IsEvil = false // Determined at commit time
			observation.add()
		}
		if len(voteSelection) > 0 {
			var observation VoteObservation
			observation.Pending = true
			observation.Round = Game.RoundNum
			observation.Subject = p.Num
			observation.Target = getPlayerByName(voteSelection).Num
			observation.add()
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
		var voteTargets = make(map[int]int)
		for _, o := range observations {
			if o.getType() == "VoteObservation" && Game.RoundNum == o.getRound() {
				target, err := o.getTarget()
				if err == nil {
					voteTargets[target]++
				}
			}
		}

		remainingPlayers := 0
		for _, p := range Players {
			if !playerIsDead(p) {
				remainingPlayers++
			}
		}

		for t, n := range voteTargets {
			if n > remainingPlayers/2 {
				votedPlayer := getPlayerByNumber(t)

				var observation LynchObservation
				observation.Pending = false
				observation.Round = Game.RoundNum
				observation.Subject = votedPlayer.Num
				observation.add()
				break
			}
		}

		// Take all observations out of pending
		CommitObservations()

		CollapseAll()

		// Immediately kill anybody who is dead
		for _, p := range Players {
			// TODO: Issue 3: PlayerDeadPercent did not give 100% dead when it should have.
			deadPercent := PlayerDeadPercent(p)
			if deadPercent == 100 {
				var observation KillObservation
				observation.Pending = false
				observation.Round = Game.RoundNum
				observation.Subject = p.Num
				observation.Role = collapseToFixedRole(p.Num)
				observation.add()
			}
		}

		var nightBoolString = ""
		if Game.RoundNight {
			Game.RoundNum++
			Game.RoundNight = false
			nightBoolString = "FALSE"
		} else {
			Game.RoundNight = true
			nightBoolString = "TRUE"
		}

		var dbStatement = "UPDATE game SET "
		dbStatement += "round=" + strconv.Itoa(Game.RoundNum)
		dbStatement += ", nightPhase=" + nightBoolString
		dbStatement += " WHERE id=" + gameID
		quantumutilities.DbExec(c, db, dbStatement)
	}

	FillActionsWithObservations()
	for _, p := range Players {
		var dbStatement = "UPDATE players SET "
		dbStatement += "actions = '" + p.Actions + "' WHERE num=" + strconv.Itoa(p.Num) + " AND gameId=" + gameID
		quantumutilities.DbExec(c, db, dbStatement)
	}

	rebuildGame(c, int(gameIDNum))
	CollapseAll()
	showGame(c)
}
