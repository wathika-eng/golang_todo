package types

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID         uint       `bun:"id,pk,autoincrement"`
	Email      string     `bun:"email,notnull,unique"`
	Password   string     `bun:"password,notnull"`
	Notes      []Note     `bun:"rel:has-many,join:id=user_id"`
	Created_At time.Time  `bun:"createdat,notnull,default:current_timestamp"`
	Updated_At *time.Time `bun:"updatedat,nullzero"`
}

type Note struct {
	bun.BaseModel `bun:"table:notes"`
	ID            uint       `bun:"id,pk,autoincrement"`
	Body          string     `bun:"body,notnull"`
	Completed     bool       `bun:"completed,default:false"`
	Created_At    time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	Updated_At    *time.Time `bun:"updated_at,nullzero,default:current_timestamp"`
	DeletedAt     *time.Time `bun:"deleted_at,soft_delete,nullzero"`
	UserID        uint       `bun:"user_id,notnull"`
}
