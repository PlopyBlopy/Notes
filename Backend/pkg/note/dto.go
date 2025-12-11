package note

import (
	"time"
)

type Note struct {
	Id          int
	Title       string
	Description string
}
type Tag struct {
	Id      int
	Title   string
	ColorId int
}

type Theme struct {
	Id    int
	Title string
}

// срезы для цветов тега и карточек заметок
type Color struct {
	Id       int
	Name     string
	Variable string
}

type NoteCard struct {
	Note        Note
	Completed   bool
	ThemeId     int
	TagsId      []int
	NoteColorId int
	CreatedAt   time.Time
}

type CreateNote struct {
	Title       string
	Description string
	ThemeId     int
	TagIds      []int
	NoteColorId int
}

type UpdateNote struct {
	Id          int
	Title       string
	Description string
	ThemeId     int
	TagIds      []int
	NoteColorId int
}
