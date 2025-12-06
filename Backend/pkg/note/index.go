package note

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Index struct {
	NoteTitles       map[string]int // title: noteId
	Themes           map[int][]int  // themeId: noteId
	Tags             map[int][]int  // tagId: noteId
	OffSize          []OffSize
	NoteIndexes      []NoteIndex
	CompletedNotes   []Note
	UncompletedNotes []Note
}

type OffSize struct {
	Off  int64
	Id   int
	Size int
}

type IndexManager struct {
	i               Index
	metadataManager IMetadataManager
}

type IIndexManager interface {
	AddNote(note Note) error
	AddNoteIndex(noteIndex NoteIndex) error

	GetNotes(completed bool, cursor, limit int, noteIds ...int) ([]Note, int, error)
	GetNoteIndexesFilteredNoteIds(noteIds ...int) ([]NoteIndex, error)

	GetFilteredTitleNoteIds(search string) ([]int, error)
	GetFilteredTagNoteIds(tagIds ...int) ([]int, error)
	GetFilteredThemeNoteIds(themeId int) ([]int, error)
}

func NewIndexManager(mm IMetadataManager) (*IndexManager, error) {
	im := &IndexManager{
		i: Index{
			NoteTitles:       map[string]int{},
			Themes:           map[int][]int{},
			Tags:             map[int][]int{},
			NoteIndexes:      []NoteIndex{},
			OffSize:          []OffSize{},
			CompletedNotes:   []Note{},
			UncompletedNotes: []Note{},
		},
		metadataManager: mm,
	}

	return im, nil
}

// need context
func (im *IndexManager) Scan() error {
	stages := [][]func() error{
		{im.scanNoteIndex},
		{im.scanNote, im.scanOffSize, im.scanNoteTheme, im.scanNoteTag},
		{im.scanNoteTitle},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var scanError error
	var once sync.Once
	var wg sync.WaitGroup

	for _, scanStages := range stages {
		if ctx.Err() != nil {
			break
		}

		for _, scan := range scanStages {
			if ctx.Err() != nil {
				break
			}

			wg.Add(1)
			go func(f func() error) {
				defer wg.Done()

				if err := f(); err != nil {
					once.Do(func() {
						scanError = err
						cancel()
					})
				}
			}(scan)
		}

		wg.Wait()
	}

	if scanError != nil {
		return scanError
	}

	return nil
}

func (im *IndexManager) scanNoteIndex() error {
	p := filepath.Join(im.metadataManager.BasePath(), im.metadataManager.IndexPath(), im.metadataManager.NoteIndexFileName())
	b, _ := os.ReadFile(p)
	ni := []NoteIndex{}
	err := json.Unmarshal(b, &ni)
	if err != nil {

	}

	im.i.NoteIndexes = ni

	return nil
}

func (im *IndexManager) scanNote() error {
	p := filepath.Join(im.metadataManager.BasePath(), im.metadataManager.NotePath(), im.metadataManager.NoteFileName())
	b, _ := os.ReadFile(p)
	n := []Note{}
	json.Unmarshal(b, &n)

	for i := 0; i < len(n); i++ {

		if im.i.NoteIndexes[i].Id != n[i].Id {
			return errors.New("failed scanNote: NoteIndexes.Id not equal note.Id")
		}

		if im.i.NoteIndexes[i].Completed {
			im.i.CompletedNotes = append(im.i.CompletedNotes, n[i])
		} else {
			im.i.UncompletedNotes = append(im.i.UncompletedNotes, n[i])
		}
	}
	return nil
}

func (im *IndexManager) scanOffSize() error {
	offSize := make([]OffSize, 0, len(im.i.NoteIndexes))

	for _, ni := range im.i.NoteIndexes {
		offSize = append(offSize, OffSize{
			Id:   ni.Id,
			Off:  ni.Off,
			Size: ni.Size,
		})
	}

	im.i.OffSize = offSize

	return nil
}

func (im *IndexManager) scanNoteTitle() error {
	noteTitles := make(map[string]int, len(im.i.CompletedNotes)+len(im.i.UncompletedNotes))
	completedNotes := im.i.CompletedNotes

	for _, n := range completedNotes {
		key := strings.ToLower(n.Title)
		noteTitles[key] = n.Id
	}

	uncompletedNotes := im.i.UncompletedNotes

	for _, n := range uncompletedNotes {
		key := strings.ToLower(n.Title)
		noteTitles[key] = n.Id
	}

	im.i.NoteTitles = noteTitles

	return nil
}
func (im *IndexManager) scanNoteTheme() error {
	themeIds, err := im.metadataManager.GetThemeIds()
	if err != nil {

	}

	themes := make(map[int][]int, len(themeIds))
	noteIndexes := im.i.NoteIndexes

	for i := 0; i < len(themeIds); i++ {
		for _, noteIndex := range noteIndexes {
			if noteIndex.ThemeId == themeIds[i] {
				themes[themeIds[i]] = append(themes[themeIds[i]], noteIndex.Id)
			} else if themeIds[i] == 0 {
				themes[themeIds[i]] = append(themes[themeIds[i]], noteIndex.Id)
			}
		}
	}

	im.i.Themes = themes

	return nil
}
func (im *IndexManager) scanNoteTag() error {
	TagsIds, err := im.metadataManager.GetTagIds()
	if err != nil {

	}

	tags := make(map[int][]int, len(TagsIds))
	noteIndexes := im.i.NoteIndexes

	for _, tagId := range TagsIds {
		for _, noteIndex := range noteIndexes {
			for _, noteTagId := range noteIndex.TagIds {
				if noteTagId == tagId {
					tags[tagId] = append(tags[tagId], noteIndex.Id)
					break
				}
			}
		}
	}

	im.i.Tags = tags

	return nil
}

func (im *IndexManager) AddNote(note Note) error {
	im.i.UncompletedNotes = append(im.i.UncompletedNotes, note)

	title := strings.ToLower(note.Title)
	im.i.NoteTitles[title] = note.Id

	return nil
}

func (im *IndexManager) AddNoteIndex(noteIndex NoteIndex) error {
	im.i.NoteIndexes = append(im.i.NoteIndexes, noteIndex)
	im.i.OffSize = append(im.i.OffSize, OffSize{
		Id:   noteIndex.Id,
		Off:  noteIndex.Off,
		Size: noteIndex.Size,
	})

	if noteIndex.ThemeId != 0 {
		im.i.Themes[0] = append(im.i.Themes[0], noteIndex.Id)
	}
	im.i.Themes[noteIndex.ThemeId] = append(im.i.Themes[noteIndex.ThemeId], noteIndex.Id)

	for _, tagId := range noteIndex.TagIds {
		im.i.Tags[tagId] = append(im.i.Tags[tagId], noteIndex.Id)
	}

	return nil
}

func (im *IndexManager) GetNotes(completed bool, cursor, limit int, noteIds ...int) ([]Note, int, error) {
	cap := limit
	if len(noteIds) < limit {
		cap = len(noteIds)
	}

	notes := make([]Note, 0, cap)
	var indexNotes []Note

	if completed {
		indexNotes = im.i.CompletedNotes
	} else {
		indexNotes = im.i.UncompletedNotes
	}

	if cursor >= len(indexNotes) {
		return nil, cursor, fmt.Errorf("cursor exit of range")
	}

	noteIdIndex := 0

	for i := cursor; i < len(indexNotes); i++ {
		if len(notes) == cap {
			break
		}

		if noteIdIndex < len(noteIds) && indexNotes[i].Id == noteIds[noteIdIndex] {
			notes = append(notes, indexNotes[i])
			cursor = i + 1
			noteIdIndex++
		}
	}

	return notes, cursor, nil
}

func (im *IndexManager) GetNoteIndexesFilteredNoteIds(noteIds ...int) ([]NoteIndex, error) {
	notes := make([]NoteIndex, 0, len(noteIds))
	noteIndexes := im.i.NoteIndexes

	for _, n := range noteIndexes {
		for _, id := range noteIds {
			if n.Id == id {
				notes = append(notes, n)
			}
		}
	}

	return notes, nil
}

func (im *IndexManager) GetFilteredTitleNoteIds(search string) ([]int, error) {
	ids := []int{}
	titles := im.i.NoteTitles
	title := strings.ToLower(search)

	for k, v := range titles {
		if strings.Contains(k, title) {
			ids = append(ids, v)
		}
	}

	return ids, nil
}

func (im *IndexManager) GetFilteredTagNoteIds(tagIds ...int) ([]int, error) {
	tags := im.i.Tags

	// срезы что содержатся у тегов
	noteIds := [][]int{}

	for k, v := range tags {
		for _, id := range tagIds {
			if k == id {
				noteIds = append(noteIds, v)
			}
		}

	}

	// карта с k=noteId, v=количеству упоминаний
	temp := map[int]int{}

	for i := 0; i < len(noteIds); i++ {
		for _, id := range noteIds[i] {
			temp[id] = temp[id] + 1
		}
	}

	// результирующий срез с noteId, содержащиеся в искомом/искомых теге/тегах
	res := []int{}

	for k, v := range temp {
		if v == len(tagIds) {
			res = append(res, k)
		}
	}

	return res, nil
}

func (im *IndexManager) GetFilteredThemeNoteIds(themeId int) ([]int, error) {
	ids := []int{}
	themes := im.i.Themes

	if themeId < 0 {
		return ids, nil
	}

	for k, v := range themes {
		if k == themeId {
			ids = append(ids, v...)
			break
		}
	}

	return ids, nil
}
