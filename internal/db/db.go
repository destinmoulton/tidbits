package db

import (
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"path/filepath"
	"tidbits/internal/db/sqlc"
	"tidbits/internal/logger"

	"tidbits/internal/utils"
)

const dbfilename = "tidbits.sqlite3"

type TidbitsDB struct {
	db      *sql.DB
	log     *logger.Logger
	queries *sqlc.Queries
}

type table struct {
	name string
}

func NewTidbitsDB(log *logger.Logger, migrations embed.FS) (*TidbitsDB, error) {
	confpath := utils.GetConfigDir()
	dbpath := filepath.Join(confpath, dbfilename)

	// Run migrations (db/migrations)
	err := runMigrations(dbpath, migrations)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbpath)

	if err != nil {
		return nil, err
	}

	queries := sqlc.New(db)

	return &TidbitsDB{
		db:      db,
		log:     log,
		queries: queries,
	}, nil
}

func (t *TidbitsDB) Close() {
	t.db.Close()
}

func (t *TidbitsDB) Init() {
	t.getListTables()
}

func (t *TidbitsDB) getListTables() {
	// Query to list tables
	query := `SELECT name FROM sqlite_master WHERE type='table'`

	// Execute the query
	rows, err := t.db.Query(query)

	if err != nil {
		t.log.Fatal("failed to query the database for list of tables: ", err)
	}
	defer rows.Close()

	// Iterate through the result set and print table names
	fmt.Println("Tables in the database:")
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(tableName)
	}
}
