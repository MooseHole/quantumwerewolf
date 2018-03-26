package main

import (
	"internal/app/quantumwerewolf"

	_ "github.com/lib/pq"
)

func main() {
	quantumwerewolf.SetupRoutes()
}
