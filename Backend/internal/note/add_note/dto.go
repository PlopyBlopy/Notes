package addnote

type input struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ThemeId     int    `json:"themeId"`
	TagIds      []int  `json:"tagIds"`
	NoteColorId int    `json:"noteColorId"`
}
