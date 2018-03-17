package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
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
		stringValue := strings.Join(value, "")
		if stringValue != "" && strings.Contains(key, "player") {
			insertStatement := "INSERT INTO players (name, num) VALUES ('" + stringValue + "', " + strconv.Itoa(playerNum) + ")"
			c.String(http.StatusOK, fmt.Sprintf(insertStatement+"\n"))
			if _, err := db.Exec(insertStatement); err != nil {
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

// Bird is for testing with this bird structure
type Bird struct {
	Species     string `json:"species"`
	Description string `json:"description"`
}

var birds []Bird

func getBirdHandler(c *gin.Context) {
	//Convert the "birds" variable to json
	birdListBytes, err := json.Marshal(birds)

	// If there is an error, print it to the console, and return a server
	// error response to the user
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error getting bird: %v", err))
		return
	}
	// If all goes well, write the JSON list of birds to the response
	c.Writer.Write(birdListBytes)
}

func createBirdHandler(c *gin.Context) {
	// Create a new instance of Bird
	bird := Bird{}

	// We send all our data as HTML form data
	// the `ParseForm` method of the request, parses the
	// form values
	err := c.Request.ParseForm()

	// In case of any error, we respond with an error to the user
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error creating bird: %v", err))
		return
	}

	// Get the information about the bird from the form info
	bird.Species = c.Request.Form.Get("species")
	bird.Description = c.Request.Form.Get("description")

	// Append our existing list of birds with a new entry
	birds = append(birds, bird)

	//Finally, we redirect the user to the original HTMl page
	c.HTML(http.StatusOK, "bird.gtpl", nil)
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

	router.POST("/playersResult", addPlayers)

	router.GET("/repeat", repeatFunc)
	router.GET("/db", dbFunc)

	router.GET("/birdy", func(c *gin.Context) {
		c.HTML(http.StatusOK, "bird.gtpl", nil)
	})
	router.GET("/bird", getBirdHandler)
	router.POST("/bird", createBirdHandler)

	router.Run(":" + port)
}
