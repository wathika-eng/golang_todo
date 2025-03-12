package repository

import (
	"fmt"
	"golang_todo/pkg/types"

	"github.com/google/uuid"
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
	_, err := r.db.NewInsert().Model(&notes).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error inserting new note: %v", err.Error())
	}
	return nil
}

func (r *NotesRepository) GetAllNotes(userID uuid.UUID) ([]types.Note, error) {
	var notes []types.Note
	err := r.db.NewSelect().Model(&notes).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching notes: %v", err)
	}
	return notes, nil
}

func (r *NotesRepository) GetNoteByID(noteID uuid.UUID) (*types.Note, error) {
	var note types.Note
	err := r.db.NewSelect().Model(&note).Where("id = ?", noteID).Limit(1).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't get notes with id: %v", noteID)
	}
	return &note, nil
}

func (r *NotesRepository) UpdateWithID(noteID uuid.UUID, updatedFields map[string]interface{}) (*types.Note, error) {
	note, err := r.GetNoteByID(noteID)
	if err != nil {
		return nil, fmt.Errorf("couldn't get note with id %v: %w", noteID, err)
	}

	if len(updatedFields) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query := r.db.NewUpdate().Model(note).WherePK()
	for field, value := range updatedFields {
		query = query.Set(fmt.Sprintf("%s = ?", field), value)
	}

	_, err = query.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error updating note: %w", err)
	}
	return r.GetNoteByID(noteID)
}

func (r *NotesRepository) DeleteWithID(noteID uuid.UUID) (bool, error) {
	var note types.Note

	_, err := r.GetNoteByID(noteID)
	if err != nil {
		return false, fmt.Errorf("couldn't get note with id: %v", noteID)
	}

	_, err = r.db.NewUpdate().
		Model(&note).
		Set("deleted_at = NOW()").
		Where("id = ?", noteID).
		Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("error soft-deleting note: %w", err)
	}

	return true, nil
}

func (r *NotesRepository) SoftDelete(userID uuid.UUID) ([]types.Note, error) {
	var notes []types.Note

	err := r.db.NewSelect().
		Model(&notes).
		Where("user_id = ?", userID).
		Where("deleted_at IS NOT NULL").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching recently deleted notes: %w", err)
	}
	fmt.Println(notes)
	return notes, nil
}
