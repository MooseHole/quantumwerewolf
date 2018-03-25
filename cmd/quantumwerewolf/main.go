package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/MooseHole/quantumwerewolf/internal/app/quantumwerewolf"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

func main() {
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
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/startPlayersSetup", func(c *gin.Context) {
		resetVars()
		c.HTML(http.StatusOK, "players.gtpl", nil)
	})
	router.GET("/startGameSetup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "gameSetup.gtpl", nil)
	})

	router.GET("/setupPlayers", quantumwerewolf.getPlayerHandler)
	router.POST("/setupPlayers", quantumwerewolf.createPlayerHandler)
	router.GET("/setupGame", getRolesHandler)
	router.POST("/setupGame", setRolesHandler)
	router.POST("/startGame", startGame)
	router.GET("/getGames", games.getGamesHandler)
	router.GET("/games", func(c *gin.Context) {
		c.HTML(http.StatusOK, "games.gtpl", nil)
	})
	router.GET("/drop", dropTables)

	router.Run(":" + port)
}