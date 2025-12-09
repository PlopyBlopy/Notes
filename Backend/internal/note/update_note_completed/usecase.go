package updateNoteCompleted

type INoteService interface {
	UpdateNoteCompleted(id int, completed bool) error
}

func NewUsecase(noteService INoteService) func(i input) error {
	return func(i input) error {
		err := noteService.UpdateNoteCompleted(i.Id, i.Completed)
		if err != nil {
			return err
		}

		return nil
	}
}
