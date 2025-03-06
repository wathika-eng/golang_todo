package repository

import (
	"fmt"
	"golang_todo/pkg/types"

	"github.com/uptrace/bun"
)

type NotesRepository struct {
	db *bun.DB
}

func NewNotesRepo(db *bun.DB) *NotesRepository {
	return &NotesRepository{
		db: db,
	}
}

func (r *NotesRepository) InsertNotes(notes types.Note) error {
	resp, err := r.db.NewInsert().Model(&notes).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error inserting new note: %v", err.Error())
	}
	println(resp)
	return nil
}
