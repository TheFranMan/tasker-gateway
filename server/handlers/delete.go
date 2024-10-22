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
	errInvalidID      = "invalid ID"
	errContentType    = "invalid content type"
)

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if "application/json" != r.Header.Get("Content-Type") {
		http.Error(w, errContentType, http.StatusUnsupportedMediaType)
		return
	}

	var deleteParams DeleteParams
	err := json.NewDecoder(r.Body).Decode(&deleteParams)
	if nil != err {
		log.WithError(err).Error(errDeseraliseJSON)
		http.Error(w, errDeseraliseJSON, http.StatusInternalServerError)
		return
	}

	if !common.ValidID(deleteParams.ID) {
		http.Error(w, errInvalidID, http.StatusBadRequest)
		return
	}

	token, err := h.app.Repo.NewDelete(r.Header.Get("Authorization"), deleteParams.ID)
	if nil != err {
		log.WithError(err).Error(errDeleteSave)
		http.Error(w, errDeleteSave, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(TokenResponse{token})
	if nil != err {
		log.WithError(err).Error(errSeraliseJSON)
		http.Error(w, errSeraliseJSON, http.StatusInternalServerError)
		return
	}
}
