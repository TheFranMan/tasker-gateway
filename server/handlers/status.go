package handlers

import (
	"encoding/json"
	"gateway/common"
	"net/http"

	"github.com/TheFranMan/tasker-common/types"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	errInvalidToken   = "invalid token"
	errGetStatus      = "cannot retrieve status from token"
	errStatusResponse = "cannot marshall token status reponse"
	errCacheSave      = "cannot save to the cache"
	errCacheGet       = "cannot retrieve from the cache"
)

type statusResponse struct {
	Status string `json:"status"`
}

func (h *Handlers) Status(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]

	if !common.ValidToken(token) {
		http.Error(w, errInvalidToken, http.StatusBadRequest)
		return
	}

	l := log.WithField("token", token)

	status, err := h.app.Cache.GetKey(token)
	if nil != err {
		l.WithError(err).Error(errCacheGet)
		http.Error(w, errCacheGet, http.StatusInternalServerError)
		return
	}

	if "" != status {
		l.WithField("status", status).Debug("Cache hit")

		err = sendResponse(w, status)
		if nil != err {
			http.Error(w, errStatusResponse, http.StatusInternalServerError)
		}

		return
	}

	responseStatus, err := h.app.Repo.GetStatus(token)
	if nil != err {
		l.WithError(err).Error(errGetStatus)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	status = string(getRequestStatusString(responseStatus))

	err = h.app.Cache.SetKey(token, status)
	if nil != err {
		l.WithError(err).Error(errCacheSave)
		http.Error(w, errCacheSave, http.StatusInternalServerError)
		return
	}

	err = sendResponse(w, status)
	if nil != err {
		l.WithError(err).Error(errStatusResponse)
		http.Error(w, errStatusResponse, http.StatusInternalServerError)
	}
}

func getRequestStatusString(responseStatus types.RequestStatus) types.RequestStatusString {
	var status types.RequestStatusString

	switch responseStatus {
	case 0:
		status = types.RequestStatusStringNew
	case 1:
		status = types.RequestStatusStringInProgress
	case 2:
		status = types.RequestStatusStringCompleted
	case 3:
		status = types.RequestStatusStringFailed
	}

	return status
}

func sendResponse(w http.ResponseWriter, status string) error {
	w.Header().Add("Content-type", "application/json")
	return json.NewEncoder(w).Encode(statusResponse{
		Status: status,
	})
}
