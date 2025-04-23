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

func InitDB() (*bun.DB, error) {
	var dsn string
	if Envs.ConnectionString == "" {
		dsn = fmt.Sprintf("%v://%v:%v@%v:%v/%v?sslmode=disable",
			Envs.DbType, Envs.DbUser, Envs.DbPassword,
			Envs.DbHost, Envs.DbPort, Envs.DbName)
	} else {
		dsn = Envs.ConnectionString
	}
	logging.Logger.Info(dsn)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	if sqldb == nil {
		logging.Logger.Error("Failed to create SQL DB connection")
		os.Exit(1)
	}
	if err := sqldb.Ping(); err != nil {
		nErr := fmt.Sprintf("Error connecting to the database: %v", err.Error())
		return nil, fmt.Errorf(nErr)
	}
	if err := healthCheck(sqldb); err != nil {
		return nil, err
	}
	DB := bun.NewDB(sqldb, pgdialect.New())
	logging.Logger.Info("✅ Database connected successfully")

	return DB, nil
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

	// Set connection pool settings (optional, can be moved to initialization)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)

	// Get database statistics
	stats := db.Stats()

	// Log health check results
	logging.Logger.Info(fmt.Sprintf("✅ Database is healthy (Ping Time: %v)", duration))
	logging.Logger.Info(fmt.Sprintf("📊 DB Stats - Open Connections: %d, In Use: %d, Idle: %d, Wait Count: %d",
		stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount))

	// Check for slow database performance
	if duration > 100*time.Millisecond {
		logging.Logger.Warn(fmt.Sprintf("⚠️ Database ping time is high: %v", duration))
	}
	if duration > 500*time.Millisecond {
		logging.Logger.Error(fmt.Sprintf("🚨 Database ping time is critically high: %v", duration))
	}
	if stats.WaitCount > 0 {
		logging.Logger.Warn(fmt.Sprintf("⚠️ Connection pool exhausted: %d connections waiting", stats.WaitCount))
	}

	return nil
}
