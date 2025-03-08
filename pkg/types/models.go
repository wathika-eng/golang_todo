package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID         uint       `json:"user_id" bun:",pk,type:uuid,default:uuid_generate_v7()"`
	Email      string     `json:"email" binding:"required,email,contains=@gmail.com" bun:"email,notnull,unique"`
	Role       string     `json:"role" bun:"default:'user'"`
	Password   string     `json:"password" binding:"required,gt=8" bun:"password,notnull"`
	Notes      []Note     `json:"notes" bun:"rel:has-many,join:id=user_id"`
	Created_At time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	Last_Login time.Time  `bun:"last_login,notnull,default:current_timestamp"`
	Updated_At *time.Time `bun:"updated_at,nullzero"`
}

type Note struct {
	bun.BaseModel `bun:"table:notes"`
	ID            uuid.UUID `json:"notes_id" bun:",pk,type:uuid,default:uuid_generate_v4()"`
	// unique:user_notes
	Body       string     `json:"notes_body" binding:"required,gt=5" bun:"body,notnull"`
	Completed  bool       `json:"completed"  bun:"completed,default:false"`
	Created_At time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	Updated_At *time.Time `bun:"updated_at,nullzero,default:current_timestamp"`
	DeletedAt  *time.Time `bun:"deleted_at,soft_delete,nullzero"`
	UserID     uint       `json:"user_id" bun:"user_id,notnull"`
}
