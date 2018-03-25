package quantumwerewolf

import (
	"encoding/json"
	"quantumwerewolf/pkg/quantumutilities"

	"github.com/gin-gonic/gin"
)

func getGamesHandler(c *gin.Context) {
	rows, err := db.Query("SELECT id, name FROM game")
	if quantumutilities.HandleErr(c, err, "Error selecting games") {
		return
	}

	var games []Game
	for rows.Next() {
		var thisGame Game
		err = rows.Scan(&thisGame.Number, &thisGame.Name)
		if quantumutilities.HandleErr(c, err, "Error scanning games") {
			return
		}

		games = append(games, thisGame)
	}

	gameListBytes, err := json.Marshal(games)

	if quantumutilities.HandleErr(c, err, "Error getting games") {
		return
	}

	c.Writer.Write(gameListBytes)
}
