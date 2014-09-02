package connection

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/russross/meddler"

	"github.com/drone/drone/server/database/migrationutil"
	"github.com/drone/drone/server/helper"
)

type Connection struct {
	DB *sql.DB
}

func NewConnection(driver, datasource string) *Connection {
	c := Connection{}

	db, err := SetupConnection(driver, datasource)
	if err != nil {
		log.Fatalln(err)
	}
	c.DB = db
	return &c
}

func SetupConnection(driver, datasource string) (*sql.DB, error) {
	wd, _ := os.Getwd()

	if driver == "" {
		driver = "sqlite3"
		log.Println("WARNING: database driver is missing, use:", driver)
	}

	if datasource == "" {
		datasource = fmt.Sprintf("%s/drone.sqlite", wd)
		log.Println("WARNING: database datasource is missing, use:", datasource)
	}

	switch driver {
	case "sqlite3":
		meddler.Default = meddler.SQLite
		migrationutil.Driver = migrationutil.SQLite
	case "mysql":
		meddler.Default = meddler.MySQL
		migrationutil.Driver = migrationutil.MySQL
	case "postgres":
		meddler.Default = meddler.PostgreSQL
		migrationutil.Driver = migrationutil.PostgreSQL
	default:
		error_message := fmt.Sprintf("unsupported driver: %s", driver)
		driver_error := errors.New(error_message)
		return nil, driver_error
	}

	helper.Driver = driver

	return sql.Open(driver, datasource)
}

func (c *Connection) MigrateAll() error {
	migrations := migrationutil.New(c.DB)
	return migrations.All().Migrate()
}

func (c *Connection) Close() {
	c.DB.Close()
}
