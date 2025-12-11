package updatenote

type input struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ThemeId     int    `json:"themeId"`
	TagIds      []int  `json:"tagIds"`
	NoteColorId int    `json:"noteColorId"`
}
