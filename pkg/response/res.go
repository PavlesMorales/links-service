package res

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(statusCode int, w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)

}
