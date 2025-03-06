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

func (r *NotesRepository) GetAllNotes() (*types.Note, error) {
	var notes types.Note
	err := r.db.NewSelect().Model(&notes).Scan(ctx, &notes)
	if err != nil {
		return nil, fmt.Errorf("error inserting new note: %v", err.Error())
	}
	fmt.Println(notes)
	return &notes, nil
}
