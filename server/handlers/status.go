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

	status, err := h.app.Repo.GetStatus(token)
	if nil != err {
		l.WithError(err).Error(errGetStatus)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(statusResponse{
		Status: string(getRequestStatusString(status)),
	})
	if nil != err {
		http.Error(w, errStatusResponse, http.StatusInternalServerError)
	}
}

func getRequestStatusString(status types.RequestStatus) types.RequestStatusString {
	var statusString types.RequestStatusString

	switch status {
	case 0:
		statusString = types.RequestStatusStringNew
	case 1:
		statusString = types.RequestStatusStringInProgress
	case 2:
		statusString = types.RequestStatusStringCompleted
	case 3:
		statusString = types.RequestStatusStringFailed
	}

	return statusString
}
