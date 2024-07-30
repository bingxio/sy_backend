package api

import (
	"encoding/json"
	"net/http"
)

func ServerError(w http.ResponseWriter, err error) {
	w.Write([]byte(err.Error()))
	w.WriteHeader(http.StatusInternalServerError)
}

// func BadRequest(w http.ResponseWriter) {
// 	w.WriteHeader(http.StatusBadRequest)
// }

//	func NotFound(w http.ResponseWriter) {
//		w.WriteHeader(http.StatusNotFound)
//	}
func JSON(w http.ResponseWriter, model any) {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.Encode(map[string]any{"data": model})
}
