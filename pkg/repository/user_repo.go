// handle database operations
package repository

import (
	"context"
	"fmt"
	"golang_todo/pkg/types"

	"github.com/uptrace/bun"
)

// holds a reference to the database

type UserRepo struct {
	db *bun.DB
}

// a constructor which ensures all db calls go through one place
func NewUserRepo(db *bun.DB) *UserRepo {
	return &UserRepo{db: db}
}

var ctx = context.Background()

func (r *UserRepo) InsertUser(user *types.User) error {
	resp, err := r.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error inserting user: %v", err)
	}
	fmt.Print(resp)
	return nil
}
