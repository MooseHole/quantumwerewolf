package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Bird is for testing with this bird structure
type Bird struct {
	Species     string `json:"species"`
	Description string `json:"description"`
}

var birds []Bird

func getBirdHandler(c *gin.Context) {
	if len(birds) == 0 {
		c.HTML(http.StatusOK, "bird.gtpl", nil)
		return
	}

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

	c.String(http.StatusOK, bird.Species)
	//Finally, we redirect the user to the original HTMl page
	c.HTML(http.StatusOK, "bird.gtpl", nil)
}
