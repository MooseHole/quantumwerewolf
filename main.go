package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	repeat int
	db     *sql.DB
)

func repeatFunc(c *gin.Context) {
	var buffer bytes.Buffer
	for i := 0; i < repeat; i++ {
		buffer.WriteString("Hello from Go!")
	}
	c.String(http.StatusOK, buffer.String())
}

func dbFunc(c *gin.Context) {
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS ticks (tick timestamp)"); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error creating database table: %q", err))
		return
	}

	if _, err := db.Exec("INSERT INTO ticks VALUES (now())"); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error incrementing tick: %q", err))
		return
	}

	rows, err := db.Query("SELECT tick FROM ticks")
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error reading ticks: %q", err))
		return
	}

	defer rows.Close()
	for rows.Next() {
		var tick time.Time
		if err := rows.Scan(&tick); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error scanning ticks: %q", err))
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("Read from the DB: %s\n", tick.String()))
	}
}

func playersResult(c *gin.Context) {
	c.Request.ParseForm()
	for key, value := range c.Request.PostForm {
		fmt.Println(key, value)
	}
}

func addPlayers(c *gin.Context) {
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS players (name varchar(40), num integer)"); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error creating database table: %q", err))
		return
	}

	c.Request.ParseForm()
	playerNum := 0
	for key, value := range c.Request.PostForm {
		if strings.Contains(key, "player") {
			if _, err := db.Exec("INSERT INTO players (name, num) VALUES (" + strings.Join(value, "") + ", " + strconv.Itoa(playerNum) + ")"); err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error adding player: %q", err))
				return
			}
			playerNum++
		}
		fmt.Println(key, value)
	}

	/*
		rows, err := db.Query("SELECT * FROM players")
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error reading ticks: %q", err))
			return
		}

		defer rows.Close()
		for rows.Next() {
			var tick time.Time
			if err := rows.Scan(&tick); err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error scanning players: %q", err))
				return
			}
			c.String(http.StatusOK, fmt.Sprintf("Read from the DB: %s\n", tick.String()))
		}
	*/
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	var err error
	tStr := os.Getenv("REPEAT")
	repeat, err = strconv.Atoi(tStr)
	if err != nil {
		log.Printf("Error converting $REPEAT to an int: %q - Using default", err)
		repeat = 5
	}

	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	//	http.HandleFunc("/", players)
	//	http.HandleFunc("/players", players)

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.LoadHTMLGlob("forms/*.gtpl")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/players", func(c *gin.Context) {
		c.HTML(http.StatusOK, "players.gtpl", nil)
	})

	router.GET("/repeat", repeatFunc)
	router.GET("/db", dbFunc)

	router.POST("/playersResult", addPlayers)

	router.Run(":" + port)
}
