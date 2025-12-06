package getfilterednotecards

import "time"

type input struct {
	Completed bool   `json:"completed"`
	Search    string `json:"search"`
	Limit     int    `json:"limit"`
	Cursor    int    `json:"cursor"`
	ThemeId   int    `json:"themeId"`
	TagIds    []int  `json:"tagIds"`
}

type output struct {
	Cards  []card `json:"cards"`
	Cursor int    `json:"cursor"`
}

type card struct {
	Note        Note      `json:"note"`
	Completed   bool      `json:"completed"`
	ThemeId     int       `json:"themeId"`
	TagIds      []int     `json:"tagIds"`
	NoteColorId int       `json:"noteColorId"`
	CreatedAt   time.Time `json:"createdAt"`
}
type Note struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
