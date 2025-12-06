package getfilterednotecards

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func NewHttpHandler(usecase func(input) (output, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-API-Version", "1.0")

		query := r.URL.Query()

		input := input{}

		if completedStr := query.Get("completed"); completedStr != "" {
			if completed, err := strconv.ParseBool(completedStr); err == nil {
				input.Completed = completed
			} else {
				input.Completed = false
			}
		}

		input.Search = query.Get("search")

		if limitStr := query.Get("limit"); limitStr != "" {
			if limit, err := strconv.Atoi(limitStr); err == nil {
				input.Limit = limit
			} else {
				input.Limit = 20
			}
		}

		if cursorStr := query.Get("cursor"); cursorStr != "" {
			if cursor, err := strconv.Atoi(cursorStr); err == nil {
				input.Cursor = cursor
			} else {
				input.Cursor = 0
			}
		}

		if themeIdStr := query.Get("themeId"); themeIdStr != "" {
			if themeId, err := strconv.Atoi(themeIdStr); err == nil {
				input.ThemeId = themeId
			} else {
				input.ThemeId = 0
			}
		}

		if tagIdsStr := query.Get("tagIds"); tagIdsStr != "" {
			tagIds := strings.Split(tagIdsStr, ",")
			for _, tagIdStr := range tagIds {
				if tagId, err := strconv.Atoi(tagIdStr); err == nil {
					input.TagIds = append(input.TagIds, tagId)
				}
			}
		}

		output, err := usecase(input)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)

		b, err := json.Marshal(output)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		w.Write(b)
	}
}
