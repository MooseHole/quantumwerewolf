package quantumwerewolf

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	db *sql.DB
)

// SetupRoutes sets up the routes
func SetupRoutes() bool {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	var err error

	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("forms/*.gtpl")
	router.Static("static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "gameList.gtpl", nil)
	})
	router.GET("/startPlayersSetup", func(c *gin.Context) {
		resetVars()
		c.HTML(http.StatusOK, "playerSetup.gtpl", nil)
	})
	router.GET("/startGameSetup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "gameSetup.gtpl", gin.H{
			"DefaultRoleName": roleTypes[0].Name,
			"Roles":           gameSetup.Roles,
		})
	})
	router.GET("/setupPlayers", getPlayerHandler)
	router.POST("/setupPlayers", createPlayerHandler)
	router.GET("/setupGame", getRolesHandler)
	router.POST("/setupGame", setRolesHandler)
	router.POST("/startGame", startGame)
	router.GET("/getGames", getGamesHandler)
	// router.GET("/game", setGame)
	router.GET("/game", func(c *gin.Context) {
		setGame(c)
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
			"Name":          gameSetup.Name,
			"TotalPlayers":  gameSetup.Total,
			"Roles":         gameSetup.Roles,
			"Round":         roundString,
			"IsNight":       game.RoundNight,
			"PlayersByName": playersByName,
			"PlayersByNum":  playersByNum,
		})
	})

	router.GET("/games", func(c *gin.Context) {
		c.HTML(http.StatusOK, "gameList.gtpl", nil)
	})
	router.GET("/drop", dropTables)

	router.Run(":" + port)

	return false
}
