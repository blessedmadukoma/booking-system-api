package main

import (
	"database/sql"
	"log"

	"booking-api/api"
	db "booking-api/db/sqlc"
	"booking-api/util"

	_ "github.com/lib/pq"
)

func main() {
	// config, err := util.LoadConfig(".")
	config, err := util.LoadEnvConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
	// connect to database
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.StartServer(config.PORT)
	if err != nil {
		log.Fatal("cannot start server!", err)
	}
}
