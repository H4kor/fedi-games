package web

import (
	"encoding/json"
	"net/http"
)

func writeJson(w http.ResponseWriter, data interface{}) error {
	s, _ := json.Marshal(data)
	w.Write(s)
	return nil
}
