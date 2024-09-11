package main

import (
	"flag"
	"fmt"
	"log"
)

func seedAccount(db storage, email, password string) *account {
	account, err := NewAccount(email, password)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.CreateAccount(account); err != nil {
		log.Fatal(err)
	}

	fmt.Println("New account => ", account.Email)

	return account
}

func seedAccounts(db storage) {
	seedAccount(db, "johndoe", "securepassword123")
}

func main() {
	seed := flag.Bool("seed", false, "seed the database")
	flag.Parse()

	db, err := NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("Seeding database")
		seedAccounts(db)
	}

	server := NewAPIServer(":3000", db)
	server.Run()
}
