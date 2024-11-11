package infrastructure

import (
	"database/sql"
	"log/slog"
	"os"

	_ "modernc.org/sqlite"
)

func NewDatabase(log *slog.Logger, name string) *sql.DB {
	// check if the databse exist
	_, err := os.Stat(name)
	shouldSeed := false
	if os.IsNotExist(err) {
		log.Warn("Database does not exist", "warn", err.Error())
		// create the database
		shouldSeed = true
	}

	db, err := sql.Open("sqlite", name)
	if err != nil {
		log.Error("Error opening database", "error", err.Error())
		return nil
	}

	err = db.Ping()
	if err != nil {
		log.Error("Error pinging database", "error", err.Error())
		return nil
	}

	if shouldSeed {
		log.Warn("Initialising Database")
		initialiseDB(db, "./")
	}
	return db
}

func SetupTestDatabase() *sql.DB {
	dbName := "../../db/test.sqlite3"
	os.Remove(dbName)
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		panic("Failed to open db database connection")
	}
	initialiseDB(db, "../../")
	return db
}

func executeFile(db *sql.DB, fileName string) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		panic("Failed to read file: " + fileName)
	}
	_, err = db.Exec(string(b))
	if err != nil {
		panic("Failed to execute file: " + fileName)
	}

}
func initialiseDB(db *sql.DB, basePath string) {
	// read seed.sql file and execute it
	setup := basePath + "db/setup.sql"
	executeFile(db, setup)
	seed := basePath + "db/seed.sql"
	executeFile(db, seed)
}
