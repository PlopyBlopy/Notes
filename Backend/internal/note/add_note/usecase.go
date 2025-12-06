package addnote

import "github.com/PlopyBlopy/notebot/pkg/note"

type INoteService interface {
	AddNote(note note.CreateNote) error
}

func NewUsecase(noteService INoteService) func(i input) error {
	return func(i input) error {
		createNote := note.CreateNote{
			Title:       i.Title,
			Description: i.Description,
			ThemeId:     i.ThemeId,
			TagIds:      i.TagIds,
			NoteColorId: i.NoteColorId,
		}

		err := noteService.AddNote(createNote)
		if err != nil {
			return err
		}

		return nil
	}
}
