package migrations

import (
	"context"
	"fmt"
	logging "golang_todo/pkg/logger"
	"golang_todo/pkg/types"
	"os"

	"github.com/uptrace/bun"
)

var ctx = context.Background()

func Migrate(db *bun.DB) {
	// users table
	_, err := db.NewCreateTable().IfNotExists().
		Model((*types.User)(nil)).Exec(ctx)
	if err != nil {
		nErr := fmt.Sprintf("❌ Failed to create users table: %v", err.Error())
		logging.Logger.Error(nErr)
		os.Exit(1)
	}
	logging.Logger.Info("✅ Users table created successfully!")
	//notes table
	_, err = db.NewCreateTable().IfNotExists().
		Model((*types.Note)(nil)).
		ForeignKey(`("user_id") REFERENCES "users" ("id") ON DELETE CASCADE`).
		Exec(ctx)
	if err != nil {
		nErr := fmt.Sprintf("❌ Failed to create notes table: %v", err.Error())
		logging.Logger.Error(nErr)
		os.Exit(1)
	}
	logging.Logger.Info("✅ notes table created successfully!")
}

func Drop(db *bun.DB) {
	err := db.ResetModel(ctx, (*types.User)(nil), (*types.Note)(nil))
	if err != nil {
		nErr := fmt.Sprintf("❌ Failed to drop tables: %v", err.Error())
		logging.Logger.Error(nErr)
		os.Exit(1)
	}
}
