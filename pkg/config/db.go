package config

import (
	"database/sql"
	"fmt"
	logging "golang_todo/pkg/logger"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var (
	DB *bun.DB
)

func InitDB() *bun.DB {
	var dsn string
	if Envs.CONNECTION_STRING == "" {
		dsn = fmt.Sprintf("%v://%v:%v@%v:%v/%v",
			Envs.DB_TYPE, Envs.DB_USER, Envs.DB_PASSWORD,
			Envs.DB_HOST, Envs.DB_PORT, Envs.DB_NAME)
	} else {
		dsn = Envs.CONNECTION_STRING
	}
	logging.Logger.Info(dsn)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	if sqldb == nil {
		logging.Logger.Error("Failed to create SQL DB connection")
		os.Exit(1)
	}
	if err := sqldb.Ping(); err != nil {
		logging.Logger.Error("Error connecting to the database: %v", err)
		os.Exit(1)
	}
	DB = bun.NewDB(sqldb, pgdialect.New())
	logging.Logger.Info(string(sqldb.Stats().OpenConnections))
	return DB
	// if err := healthCheck(DB.DB); err != nil {
	// 	log.Fatal(err)
	// }
}

// func healthCheck(db *sql.DB) error {
// 	err := db.Ping()
// 	if err != nil {
// 		return errors.New("failed to ping the database")
// 	}
// 	inUse := db.Stats().InUse
// 	idle := db.Stats().Idle
// 	fmt.Printf("%v, %v\n", inUse, idle)
// 	return nil
// }
