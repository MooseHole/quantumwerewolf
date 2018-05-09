package quantumwerewolf

import (
	"database/sql"
	"log"
	"net/http"
	"os"

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
		ResetVars()
		c.HTML(http.StatusOK, "playerSetup.gtpl", nil)
	})
	router.GET("/startGameSetup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "gameSetup.gtpl", gin.H{
			"DefaultRoleName": roleTypes[0].Name,
			"Roles":           GameSetup.Roles,
		})
	})
	router.GET("/setupPlayers", getPlayerHandler)
	router.POST("/setupPlayers", createPlayerHandler)
	router.GET("/setupGame", getRolesHandler)
	router.POST("/setupGame", setRolesHandler)
	router.GET("/getGames", getGamesHandler)
	// router.GET("/game", setGame)
	router.GET("/game", func(c *gin.Context) {
		setGame(c)
		showGame(c)
	})
	router.POST("/processActions", processActions)

	router.GET("/games", func(c *gin.Context) {
		c.HTML(http.StatusOK, "gameList.gtpl", nil)
	})
	router.GET("/drop", dropTables)

	router.Run(":" + port)

	return false
}
