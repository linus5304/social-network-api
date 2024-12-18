package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/linus5304/social/internal/db"
	"github.com/linus5304/social/internal/env"
	"github.com/linus5304/social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5433/social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	store := store.NewStorage(conn)
	db.Seed(store)
}
