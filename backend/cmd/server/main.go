package main

import (
	"log"

	"github.com/jobearz/bookmark-manager/db"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer database.Close()

	// code here
}
