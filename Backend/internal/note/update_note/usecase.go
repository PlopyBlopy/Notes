package updatenote

import (
	"fmt"

	"github.com/PlopyBlopy/notebot/pkg/note"
)

type INoteService interface {
	UpdateNote(note note.UpdateNote) error
}

func NewUsecase(noteService INoteService) func(i input) error {
	return func(i input) error {
		updateNote := note.UpdateNote{
			Id:          i.Id,
			Title:       i.Title,
			Description: i.Description,
			ThemeId:     i.ThemeId,
			TagIds:      i.TagIds,
			NoteColorId: i.NoteColorId,
		}

		err := noteService.UpdateNote(updateNote)
		if err != nil {
			return fmt.Errorf("failed update note: %w", err)
		}

		return nil
	}
}
