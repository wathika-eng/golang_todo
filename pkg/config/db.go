package config

import (
	"context"
	"database/sql"
	"fmt"
	logging "golang_todo/pkg/logger"
	"os"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
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
		nErr := fmt.Sprintf("Error connecting to the database: %v", err.Error())
		logging.Logger.Error(nErr)
		os.Exit(1)
		return nil
	}
	if err := healthCheck(sqldb); err != nil {
		logging.Logger.Error(err.Error())
		os.Exit(1)
	}
	DB := bun.NewDB(sqldb, pgdialect.New())
	logging.Logger.Info("✅ Database connected successfully")

	return DB
}

// healthCheck checks the database connection and logs relevant stats
func healthCheck(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	startTime := time.Now()
	err := db.PingContext(ctx)
	duration := time.Since(startTime)

	if err != nil {
		return fmt.Errorf("❌ Database health check failed: %v", err)
	}

	stats := db.Stats()
	logging.Logger.Info(fmt.Sprintf("✅ Database is healthy (Ping Time: %v)", duration))
	logging.Logger.Info(fmt.Sprintf("📊 DB Stats - Open Connections: %d, In Use: %d, Idle: %d, Wait Count: %d",
		stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount))

	return nil
}
