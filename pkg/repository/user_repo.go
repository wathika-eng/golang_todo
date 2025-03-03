// handle database operations
package repository

import (
	"context"
	"database/sql"
	"errors"
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

func (r *UserRepo) InsertUser(user types.User) error {
	resp, err := r.db.NewInsert().Model(&user).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error inserting user: %v", err)
	}
	fmt.Print(resp)
	return nil
}

func (r *UserRepo) GetUserByEmail(email string) (*types.User, error) {
	var user types.User
	err := r.db.NewSelect().Model(&user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error fetching user by email: %v", err)
	}
	print(user.Password)
	return &user, nil
}
