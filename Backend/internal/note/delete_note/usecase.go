package deletenote

type INoteService interface {
	DeleteNote(id int) error
}

func NewUsecase(noteService INoteService) func(i input) error {
	return func(i input) error {
		err := noteService.DeleteNote(i.Id)
		if err != nil {
			return err
		}

		return nil
	}
}
