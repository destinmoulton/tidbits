package db

import (
	"embed"
	"fmt"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/sqlite"
)

func runMigrations(path string, migrationFS embed.FS) error {
	u, _ := url.Parse("sqlite:" + path)
	dbmate := dbmate.New(u)
	dbmate.FS = migrationFS
	fmt.Println(migrationFS)

	migrations, err := dbmate.FindMigrations()
	if err != nil {
		return err
	}
	for _, m := range migrations {
		fmt.Println(m.Version, m.FilePath)
	}

	err = dbmate.CreateAndMigrate()
	if err != nil {
		fmt.Println("failing on CreateAndMigrate:", err)
		return err
	}
	return nil
}
