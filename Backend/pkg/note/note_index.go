package note

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// нету функционала для очистки удаленных note в runtime, очистка только при старте приложения
type NoteIndex struct {
	Id          int       `json:"id"`
	Completed   bool      `json:"completed"`
	Deleted     bool      `json:"deleted"`
	ThemeId     int       `json:"theme_id"` // есть всегда, default = без темы
	TagIds      []int     `json:"tag_ids"`
	NoteColorId int       `json:"note_color_id"`
	Off         int64     `json:"off"`  // offset
	Size        int       `json:"size"` // writed bytes
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type NoteIndexManager struct {
	metadataManager IMetadataManager
	indexManager    IIndexManager
}

func NewNoteIndexManager(mm IMetadataManager, im IIndexManager) (*NoteIndexManager, error) {
	ni := &NoteIndexManager{
		metadataManager: mm,
		indexManager:    im,
	}

	err := ni.existOrCreate()
	if err != nil {
		return nil, fmt.Errorf("failed new note: %w", err)
	}

	ni.metadataManager = mm

	return ni, nil
}

type INoteIndexManager interface {
	AddNoteIndex(noteId, themeId, size, noteColorId int, off int64, tagIds ...int) error
}

// Проверка есть ли noteIndex file если нет создать пустым
func (ni NoteIndexManager) existOrCreate() error {
	isentry := false

	p := filepath.Join(ni.metadataManager.BasePath(), ni.metadataManager.IndexPath())

	files, _ := os.ReadDir(p)

	noteIndexFileName := ni.metadataManager.NoteIndexFileName()
	for _, file := range files {
		if !file.IsDir() && file.Name() == noteIndexFileName {
			isentry = true
			break
		}
	}

	p = filepath.Join(ni.metadataManager.BasePath(), ni.metadataManager.IndexPath(), noteIndexFileName)

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

func (ni NoteIndexManager) AddNoteIndex(noteId, themeId, noteColorId, size int, off int64, tagIds ...int) error {
	createdAt := time.Now()

	noteIndex := NoteIndex{
		Id:          noteId,
		Completed:   false,
		Deleted:     false,
		ThemeId:     themeId,
		TagIds:      tagIds,
		NoteColorId: noteColorId,
		Off:         off,
		Size:        size,
		UpdatedAt:   createdAt,
		CreatedAt:   createdAt,
	}

	indexPath := filepath.Join(ni.metadataManager.BasePath(), ni.metadataManager.IndexPath(), ni.metadataManager.NoteIndexFileName())
	indexFile, _ := os.OpenFile(indexPath, os.O_RDWR, 0666)
	defer indexFile.Close()
	b, err := json.Marshal(noteIndex)
	if err != nil {
		return err
	}

	WriteAt(b, indexFile)

	if err := ni.indexManager.AddNoteIndex(noteIndex); err != nil {

	}

	return nil
}
