package main

import (
	"log"

	"github.com/linus5304/social/internal/db"
	"github.com/linus5304/social/internal/env"
	"github.com/linus5304/social/internal/store"
)

func main() {

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5433/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("database connectino pool established")

	store := store.NewStorage(db)
	app := &application{
		config: cfg,
		store:  store,
	}

	log.Fatal(app.run(app.mount()))
}