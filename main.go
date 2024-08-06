package main

import (
	"embed"
	"tidbits/internal/cliflags"
	"tidbits/internal/db"
	gui "tidbits/internal/gui"
	"tidbits/internal/logger"
)

// Use go embed to put them
// the db migrationsFS in the binary

//go:embed db/migrations/*.sql
var migrationsFS embed.FS

func main() {
	flags := cliflags.ParseFlags()
	log := logger.NewLogger(flags)
	defer log.Close()

	tbdb, err := db.NewTidbitsDB(log, migrationsFS)
	if err != nil {
		log.Fatal("failed to connect to db: ", err)
	}
	defer tbdb.Close()

	tbdb.Init()

	appGUI := gui.NewGUI(log, tbdb)
	appGUI.Run()
}
