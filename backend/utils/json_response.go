package utils

import (
	"encoding/json"
	"net/http"
)

// JSONResponse serializa `data` a JSON, establece los headers y escribe la respuesta.
// Si ocurre un error en la serialización, devuelve un error HTTP 500.
func JSONResponse(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, `{"error":"Failed to serialize response"}`, http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonData)
}
