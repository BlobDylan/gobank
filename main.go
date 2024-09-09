package main

import (
	"log"
)

func main() {
	db, err := NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":3000", db)
	server.Run()
}
