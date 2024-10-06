package handlers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"gateway/common"
)

var (
	errSeraliseJSON   = "cannot seralise JSON response"
	errDeleteSave     = "cannot save new deletion request"
	errDeseraliseJSON = "cannot deseralise JSON body"
)

type DeleteParams struct {
	ID int `json:"id"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func (h *Handlers) UserDelete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var deleteParams DeleteParams
	err := json.NewDecoder(r.Body).Decode(&deleteParams)
	if nil != err {
		log.WithError(err).Error(errDeseraliseJSON)
		http.Error(w, errDeseraliseJSON, http.StatusInternalServerError)
		return
	}

	if !common.ValidID(deleteParams.ID) {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	authToken := "test-token-delete"
	token, err := h.app.Repo.NewDelete(authToken, deleteParams.ID)
	if nil != err {
		log.WithError(err).Error(errSeraliseJSON)
		http.Error(w, errSeraliseJSON, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(TokenResponse{token})
	if nil != err {
		log.WithError(err).Error(errDeleteSave)
		http.Error(w, errDeleteSave, http.StatusInternalServerError)
		return
	}
}
