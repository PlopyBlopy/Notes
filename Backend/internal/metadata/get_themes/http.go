package getthemes

import (
	"encoding/json"
	"net/http"
)

func NewHttpHandler(usecase func() (output, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-API-Version", "1.0")

		output, err := usecase()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)

		b, err := json.Marshal(output.Themes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		w.Write(b)
	}
}
