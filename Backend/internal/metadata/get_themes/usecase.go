package getthemes

import "github.com/PlopyBlopy/notebot/pkg/note"

type IMetadataService interface {
	GetThemes() ([]note.Theme, error)
}

func NewUsecase(metadataService IMetadataService) func() (output, error) {
	return func() (output, error) {

		themes, err := metadataService.GetThemes()
		if err != nil {
			return output{}, err
		}

		if len(themes) == 0 {
			return output{}, nil
		}

		output := output{
			Themes: make([]theme, 0, len(themes)),
		}
		for _, t := range themes {
			output.Themes = append(output.Themes, theme{
				Id:    t.Id,
				Title: t.Title,
			})
		}

		return output, nil
	}
}
