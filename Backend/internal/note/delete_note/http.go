package deletenote

import (
	"net/http"
	"strconv"
)

func NewHttpHandler(usecase func(i input) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-API-Version", "1.0")

		query := r.URL.Query()

		input := input{}

		if idStr := query.Get("id"); idStr != "" {
			if id, err := strconv.Atoi(idStr); err == nil {
				input.Id = id
			} else {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		err := usecase(input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
