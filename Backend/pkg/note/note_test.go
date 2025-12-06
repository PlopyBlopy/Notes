package note

import (
	"fmt"
	"testing"
	"time"
)

func TestGetFilteredNotes(t *testing.T) {

	noteLen := 10
	notes := make([]Note, 0, noteLen)

	for i := 0; i < noteLen; i++ {
		notes = append(notes, Note{
			Id:          i,
			Title:       fmt.Sprintf("Title %d", i),
			Description: "Some Description for this note",
		})
	}

	noteIndexes := make([]NoteIndex, 0, noteLen)

	noteIndexes = append(noteIndexes, NoteIndex{
		Id:        0,
		Completed: false,
		Deleted:   false,
		ThemeId:   0,
		TagIds:    []int{0},
		CreatedAt: time.Now(),
	})

	noteIndexes = append(noteIndexes, NoteIndex{
		Id:        1,
		Completed: false,
		Deleted:   false,
		ThemeId:   0,
		TagIds:    []int{0},
		CreatedAt: time.Now(),
	})

	noteIndexes = append(noteIndexes, NoteIndex{
		Id:        2,
		Completed: false,
		Deleted:   false,
		ThemeId:   1,
		TagIds:    []int{0, 1},
		CreatedAt: time.Now(),
	})

	noteIndexes = append(noteIndexes, NoteIndex{
		Id:        3,
		Completed: false,
		Deleted:   false,
		ThemeId:   1,
		TagIds:    []int{1},
		CreatedAt: time.Now(),
	})

	noteIndexes = append(noteIndexes, NoteIndex{
		Id:        4,
		Completed: false,
		Deleted:   false,
		ThemeId:   2,
		TagIds:    []int{2},
		CreatedAt: time.Now(),
	})

	noteIndexes = append(noteIndexes, NoteIndex{
		Id:        5,
		Completed: true,
		Deleted:   false,
		ThemeId:   0,
		TagIds:    []int{0},
		CreatedAt: time.Now(),
	})

	noteIndexes = append(noteIndexes, NoteIndex{
		Id:        6,
		Completed: true,
		Deleted:   false,
		ThemeId:   0,
		TagIds:    []int{0},
		CreatedAt: time.Now(),
	})

	noteIndexes = append(noteIndexes, NoteIndex{
		Id:        7,
		Completed: true,
		Deleted:   false,
		ThemeId:   1,
		TagIds:    []int{0, 1},
		CreatedAt: time.Now(),
	})

	noteIndexes = append(noteIndexes, NoteIndex{
		Id:        8,
		Completed: true,
		Deleted:   false,
		ThemeId:   1,
		TagIds:    []int{1},
		CreatedAt: time.Now(),
	})

	noteIndexes = append(noteIndexes, NoteIndex{
		Id:        9,
		Completed: true,
		Deleted:   false,
		ThemeId:   2,
		TagIds:    []int{2},
		CreatedAt: time.Now(),
	})

	indexManager := IndexManager{
		i: Index{
			NoteTitles: make(map[string]int, noteLen),
		},
	}

	noteCards := make([]NoteCard, 0, noteLen)

	for i := 0; i < noteLen; i++ {
		noteCards = append(noteCards, NoteCard{
			Note:        notes[i],
			Completed:   noteIndexes[i].Completed,
			ThemeId:     noteIndexes[i].ThemeId,
			TagsId:      noteIndexes[i].TagIds,
			NoteColorId: noteIndexes[i].NoteColorId,
			CreatedAt:   noteIndexes[i].CreatedAt,
		})
	}

	for i := 0; i < noteLen; i++ {
		if noteIndexes[i].Completed {
			indexManager.i.CompletedNotes = append(indexManager.i.CompletedNotes, notes[i])
		} else {
			indexManager.i.UncompletedNotes = append(indexManager.i.UncompletedNotes, notes[i])
		}
	}

	themes := map[int][]int{0: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 1: {2, 3, 7, 8}, 2: {4, 9}}
	tags := map[int][]int{0: {0, 1, 2, 5, 6, 7}, 1: {2, 3, 7, 8}, 2: {4, 9}}

	indexManager.i.NoteIndexes = noteIndexes
	indexManager.i.Tags = tags
	indexManager.i.Themes = themes

	indexManager.scanNoteTitle()

	tests := []struct {
		name      string
		completed bool
		search    string
		limit     int
		cursor    int
		themeId   int
		tagIds    []int
		expected  []NoteCard
	}{
		{"uncompleted", false, "", 20, 0, 0, []int{}, []NoteCard{noteCards[0], noteCards[1], noteCards[2], noteCards[3], noteCards[4]}},
		{"uncompleted,theme#1", false, "", 20, 0, 2, []int{}, []NoteCard{noteCards[4]}},
		{"uncompleted,theme#2", false, "", 20, 0, 1, []int{}, []NoteCard{noteCards[2], noteCards[3]}},
		{"uncompleted,search,theme#1", false, "Tit", 20, 0, 0, []int{}, []NoteCard{noteCards[0], noteCards[1], noteCards[2], noteCards[3], noteCards[4]}},
		{"uncompleted,search,theme#2", false, "Tit", 20, 0, 1, []int{}, []NoteCard{noteCards[2], noteCards[3]}},
		{"uncompleted,search,theme#3", false, "Title 2", 20, 0, 1, []int{}, []NoteCard{noteCards[2]}},
		{"uncompleted,theme,tag", false, "", 20, 0, 0, []int{0}, []NoteCard{noteCards[0], noteCards[1], noteCards[2]}},
		{"uncompleted,theme,tags#1", false, "", 20, 0, 0, []int{0, 1}, []NoteCard{noteCards[2]}},
		{"uncompleted,theme,tags#2", false, "", 20, 0, 1, []int{0, 1}, []NoteCard{noteCards[2]}},

		{"completed,empty", true, "", 20, 0, 0, []int{}, []NoteCard{noteCards[9]}},
		{"completed,theme#1", true, "", 20, 0, 2, []int{}, []NoteCard{noteCards[9]}},
		{"completed,theme#2", true, "", 20, 0, 1, []int{}, []NoteCard{noteCards[7], noteCards[8]}},
		{"completed,search,theme#1", true, "Tit", 20, 0, 0, []int{}, []NoteCard{noteCards[5], noteCards[6], noteCards[7], noteCards[8], noteCards[9]}},
		{"completed,search,theme#2", true, "Tit", 20, 0, 1, []int{}, []NoteCard{noteCards[7], noteCards[8]}},
		{"completed,search,theme#3", true, "Title 2", 20, 0, 1, []int{}, []NoteCard{noteCards[7]}},
		{"completed,theme,tag", true, "", 20, 0, 0, []int{0}, []NoteCard{noteCards[5], noteCards[6], noteCards[7]}},
		{"completed,theme,tags#1", true, "", 20, 0, 0, []int{0, 1}, []NoteCard{noteCards[7]}},
		{"completed,theme,tags#2", true, "", 20, 0, 1, []int{0, 1}, []NoteCard{noteCards[7]}},
	}

	noteManager := NoteManager{indexManager: &indexManager}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _, err := noteManager.GetFilteredNoteCards(tt.completed, tt.search, tt.limit, tt.cursor, tt.themeId, tt.tagIds...)
			if err != nil {
				t.Errorf("failed test: %s", err)
			}

			if len(res) == 0 {
				return
			}

			// expectedCursor := tt.expected[len(tt.expected)-1].Note.Id
			// if c != expectedCursor {
			// 	t.Errorf("the resulting cursor is not equal to expected cursor")
			// }

			resCount := 0

			for i := 0; i < len(res); i++ {
				equal := true
				card := res[i]
				expect := tt.expected[i]

				resCount++

				// if card.Note != expect.Note || card.Completed != expect.Completed || card.ThemeId != expect.ThemeId || card.NoteColorId != expect.NoteColorId || card.CreatedAt != expect.CreatedAt {
				// 	equal = false
				// }

				if card.Note != expect.Note {
					equal = false
					t.Errorf("card.Note is not equal expect.Note: \ncard: %+v\n expect: %+v", card.Note, expect.Note)
				}

				if card.Completed != expect.Completed {
					equal = false
					t.Errorf("card.Completed is not equal expect.Completed: \ncard: %+v\n expect: %+v", card.Completed, expect.Completed)
				}

				if card.ThemeId != expect.ThemeId {
					equal = false
				}

				if card.NoteColorId != expect.NoteColorId {
					equal = false
				}

				if card.CreatedAt != expect.CreatedAt {
					equal = false
				}

				if len(card.TagsId) != len(expect.TagsId) {
					equal = false
				} else {
					for j := range card.TagsId {
						if card.TagsId[j] != expect.TagsId[j] {
							equal = false
							break
						}
					}
				}

				if !equal {
					t.Errorf("\nGetFilteredNotes(%s, %d, %d, %d)=%+v;\n expected %+v", tt.search, tt.limit, tt.themeId, tt.tagIds, res, tt.expected[i])
				}

			}

			if resCount < len(tt.expected) {
				t.Errorf("less values passed the check than they should have passed")
			}
		})
	}
}
