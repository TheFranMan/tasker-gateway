package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/TheFranMan/tasker-common/types"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"gateway/common"
)

var (
	errMsgInvalidToken = "Invalid token"
	errGetStatus       = "cannot retrieve status from token"
	errStatusResponse  = "cannot marshall token status reponse"
	errMsgCacheSave    = "cannot save to the cache"
	errMsgCacheGet     = "cannot retrieve from the cache"
)

type statusResponse struct {
	Status types.RequestStatusString `json:"status"`
}

func (h *Handlers) Status(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]

	if !common.ValidToken(token) {
		http.Error(w, errMsgInvalidToken, http.StatusBadRequest)
		return
	}

	l := log.WithField("token", token)

	status, err := h.app.Cache.GetKey(token)
	if nil != err {
		l.WithError(err).Error(errMsgCacheGet)
		http.Error(w, errMsgCacheGet, http.StatusInternalServerError)
		return
	}

	if nil != status {
		l.WithField("status", *status).Debug("Cache hit")

		h.app.Monitor.PathStatusCached()

		err = sendResponse(w, *status)
		if nil != err {
			http.Error(w, errStatusResponse, http.StatusInternalServerError)
		}

		return
	}

	status, err = h.app.Repo.GetStatus(token)
	if nil != err {
		l.WithError(err).Error(errGetStatus)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if nil == status {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = h.app.Cache.SetKey(token, *status)
	if nil != err {
		l.WithError(err).Error(errMsgCacheSave)
		http.Error(w, errMsgCacheSave, http.StatusInternalServerError)
		return
	}

	err = sendResponse(w, *status)
	if nil != err {
		l.WithError(err).Error(errStatusResponse)
		http.Error(w, errStatusResponse, http.StatusInternalServerError)
	}
}

func sendResponse(w http.ResponseWriter, status types.RequestStatusString) error {
	w.Header().Add("Content-type", "application/json")
	return json.NewEncoder(w).Encode(statusResponse{
		Status: status,
	})
}
