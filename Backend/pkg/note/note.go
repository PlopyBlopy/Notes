package note

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type NoteManager struct {
	metadataManager  IMetadataManager
	noteIndexManager INoteIndexManager
	indexManager     IIndexManager
}

func NewNoteManager(mm IMetadataManager, ni INoteIndexManager, im IIndexManager) (*NoteManager, error) {
	nm := &NoteManager{
		metadataManager:  mm,
		noteIndexManager: ni,
		indexManager:     im,
	}

	err := nm.existOrCreate()
	if err != nil {
		return nil, fmt.Errorf("failed new note: %w", err)
	}

	nm.metadataManager = mm

	return nm, nil
}

// Проверка есть ли note file если нет создать пустым
func (nm NoteManager) existOrCreate() error {
	isentry := false

	p := filepath.Join(nm.metadataManager.BasePath(), nm.metadataManager.NotePath())

	files, _ := os.ReadDir(p)

	noteFileName := nm.metadataManager.NoteFileName()
	for _, file := range files {
		if !file.IsDir() && file.Name() == noteFileName {
			isentry = true
			break
		}
	}

	p = filepath.Join(nm.metadataManager.BasePath(), nm.metadataManager.NotePath(), nm.metadataManager.NoteFileName())

	empty := []interface{}{}

	if !isentry {
		os.Create(p)
		b, _ := json.Marshal(empty)
		f, _ := os.OpenFile(p, os.O_RDWR, 0666)
		n, err := f.Write(b)
		if n == 0 {

		}
		if err != nil {

		}
	}

	return nil
}

func (nm *NoteManager) AddNote(createNote CreateNote) error {
	noteId := nm.metadataManager.GetNoteId()

	note := Note{
		Id:          noteId,
		Title:       createNote.Title,
		Description: createNote.Description,
	}

	notePath := filepath.Join(nm.metadataManager.BasePath(), nm.metadataManager.NotePath(), nm.metadataManager.NoteFileName())
	noteFile, _ := os.OpenFile(notePath, os.O_RDWR, 0666)
	defer noteFile.Close()

	b, _ := json.Marshal(note)
	off, size := WriteAt(b, noteFile)

	// add to indexes
	err := nm.noteIndexManager.AddNoteIndex(noteId, createNote.ThemeId, createNote.NoteColorId, size, off, createNote.TagIds...)
	if err != nil {
		// nm.RemoveLastNote()
		return err
	}

	if err := nm.indexManager.AddNote(note); err != nil {

	}

	return nil
}

func (nm NoteManager) GetFilteredNoteCards(completed bool, search string, limit, cursor, themeId int, tagIds ...int) ([]NoteCard, int, error) {
	filter := 0
	notes := map[int]int{}

	if len(search) != 0 {
		searchNotes, err := nm.indexManager.GetFilteredTitleNoteIds(search)
		if err != nil {

		}
		if len(searchNotes) != 0 {
			for _, i := range searchNotes {
				notes[i]++
			}
		}
		filter++
	}

	if len(tagIds) != 0 {
		tagNotes, err := nm.indexManager.GetFilteredTagNoteIds(tagIds...)
		if err != nil {

		}

		if len(tagNotes) != 0 {
			for _, i := range tagNotes {
				notes[i]++
			}
		}
		filter++
	}

	themeNotes, err := nm.indexManager.GetFilteredThemeNoteIds(themeId)
	if err != nil {

	}

	if len(themeNotes) != 0 {
		for _, i := range themeNotes {
			notes[i]++
		}
	}
	filter++

	res := []NoteCard{}

	if filter != 0 && len(notes) != 0 {
		noteIds := []int{}

		for k, v := range notes {
			if v == filter {
				noteIds = append(noteIds, k)
			}
		}

		sort.Ints(noteIds)

		filteredNotes, c, err := nm.indexManager.GetNotes(completed, cursor, limit, noteIds...)
		if err != nil {
			return nil, cursor, err
		}

		cursor = c

		noteIndexes, err := nm.indexManager.GetNoteIndexesFilteredNoteIds(noteIds...)
		if err != nil {

		}

		for i := 0; i < len(filteredNotes); i++ {
			note := filteredNotes[i]
			noteIndex := noteIndexes[i]

			res = append(res, NoteCard{
				Note:        note,
				Completed:   noteIndex.Completed,
				ThemeId:     noteIndex.ThemeId,
				TagsId:      noteIndex.TagIds,
				NoteColorId: noteIndex.NoteColorId,
				CreatedAt:   noteIndex.CreatedAt,
			})
		}
	}

	return res, cursor, nil
}

func (nm NoteManager) DeleteNote(id int) error {
	err := nm.indexManager.DeleteNote(id)
	if err != nil {
		return err
	}

	return nil
}
