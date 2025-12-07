package service

import (
	"fmt"

	"github.com/PlopyBlopy/notebot/pkg/note"
)

type INoteManager interface {
	AddNote(note.CreateNote) error
	GetFilteredNoteCards(completed bool, search string, limit, cursor, themeId int, tags ...int) ([]note.NoteCard, int, error)
	DeleteNote(id int) error
}

type NoteService struct {
	noteManager INoteManager
}

func NewNoteService(nm INoteManager) (*NoteService, error) {
	return &NoteService{noteManager: nm}, nil
}

func (ns NoteService) AddNote(note note.CreateNote) error {
	err := ns.noteManager.AddNote(note)
	if err != nil {
		return fmt.Errorf("failed add note: %w", err)
	}
	return nil
}

func (ns NoteService) GetFilteredNoteCards(completed bool, search string, limit, cursor, themeId int, tagIds ...int) ([]note.NoteCard, int, error) {
	noteCards, c, err := ns.noteManager.GetFilteredNoteCards(completed, search, limit, cursor, themeId, tagIds...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get filtered note cards: %w", err)
	}
	return noteCards, c, nil
}

func (ns NoteService) DeleteNote(id int) error {
	err := ns.noteManager.DeleteNote(id)
	if err != nil {
		return err
	}

	return nil
}
