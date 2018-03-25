package main

import (
	"quantumwerewolf/internal/app/quantumwerewolf"

	_ "github.com/lib/pq"
)

func main() {
	quantumwerewolf.SetupRoutes()
}
