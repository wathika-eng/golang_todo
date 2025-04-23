package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        uuid.UUID  `json:"user_id" bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Email     string     `json:"email" binding:"required,email" bun:"email,notnull,unique"`
	Role      string     `json:"role" bun:"default:'user'"`
	Password  string     `json:"password" binding:"required,gt=8" bun:"password,notnull"`
	Notes     []Note     `json:"notes" bun:"rel:has-many,join:id=user_id"`
	CreatedAt time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	LastLogin time.Time  `bun:"last_login,notnull,default:current_timestamp"`
	UpdatedAt *time.Time `bun:"updated_at,nullzero"`
}

type Note struct {
	bun.BaseModel `bun:"table:notes"`
	ID            uuid.UUID `json:"notes_id" bun:",pk,type:uuid,default:uuid_generate_v4()"`
	// unique:user_notes
	Body      string     `json:"notes_body" binding:"required,gt=5" bun:"body,notnull"`
	Completed bool       `json:"completed"  bun:"completed,default:false"`
	CreatedAt time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time  `bun:"updated_at,nullzero"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete,nullzero"`
	UserID    uuid.UUID  `json:"user_id" bun:"type:uuid,notnull"`
}
