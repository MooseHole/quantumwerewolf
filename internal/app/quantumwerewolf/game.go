package quantumwerewolf

import (
	"quantumwerewolf/pkg/quantumutilities"
	"strconv"

	"github.com/gin-gonic/gin"
)

func rebuildGame(c *gin.Context, gameID int) {
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

	err = quantumutilities.GetInterface(rolesByteArray, gameSetup.Roles)
	if quantumutilities.HandleErr(c, err, "Error getting game roles interface") {
		return
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
