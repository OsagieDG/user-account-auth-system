package handlers

import (
	"encoding/json"
	"net/http"
)

type loginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func invalidCredentials(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	resp := loginResp{
		Type: "error",
		Msg:  "invalid credentials",
	}
	err := json.NewEncoder(w).Encode(resp)
	return err
}
