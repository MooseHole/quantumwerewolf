package quantumwerewolf

import (
	"database/sql"
	tt "html/template"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type TplCtx struct {
	UserEmail tt.JS
}

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
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.LoadHTMLGlob("forms/*.gtpl")
	router.Static("static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/startPlayersSetup", func(c *gin.Context) {
		resetVars()
		c.HTML(http.StatusOK, "players.gtpl", nil)
	})
	//	router.GET("/startGameSetup", func(c *gin.Context) {
	//		c.HTML(http.StatusOK, "gameSetup.gtpl", nil)
	//	})
	router.GET("/startGameSetup", func(c *gin.Context) {
		c.Header("ContentType", "text/plain")
		c.HTML(http.StatusOK, "gameSetupTemp.gtpl", gin.H{
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
	router.GET("/games", func(c *gin.Context) {
		c.HTML(http.StatusOK, "games.gtpl", nil)
	})
	router.GET("/drop", dropTables)

	router.Run(":" + port)

	return false
}
