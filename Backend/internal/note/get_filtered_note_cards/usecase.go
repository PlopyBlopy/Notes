package getfilterednotecards

import "github.com/PlopyBlopy/notebot/pkg/note"

type INoteService interface {
	GetFilteredNoteCards(completed bool, search string, limit, cursor, themeId int, tagIds ...int) ([]note.NoteCard, int, error)
}

func NewUsecase(noteService INoteService) func(input) (output, error) {
	return func(i input) (output, error) {
		cards, c, err := noteService.GetFilteredNoteCards(i.Completed, i.Search, i.Limit, i.Cursor, i.ThemeId, i.TagIds...)
		if err != nil {
			return output{}, err
		}

		if len(cards) == 0 {
			return output{}, nil
		}

		output := output{
			Cards:  make([]card, 0, len(cards)),
			Cursor: c,
		}
		for _, c := range cards {
			output.Cards = append(output.Cards, card{
				Note:        Note(c.Note),
				Completed:   c.Completed,
				ThemeId:     c.ThemeId,
				TagIds:      c.TagsId,
				NoteColorId: c.NoteColorId,
				CreatedAt:   c.CreatedAt,
			})
		}

		return output, nil
	}
}
