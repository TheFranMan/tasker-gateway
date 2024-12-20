package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"gateway/common"
)

var (
	errMsgInvalidToken   = "Invalid token"
	errMsgStatusGet      = "cannot retrieve status from repo"
	errMsgResponseStatus = "cannot marshall token status reponse"
	errMsgCacheSave      = "cannot save status to the cache"
	errMsgCacheGet       = "cannot retrieve status from the cache"
)

func (h *Handlers) Poll(w http.ResponseWriter, r *http.Request) {
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

		h.app.Monitor.StatusCacheHit()

		err = json.NewEncoder(w).Encode(PollResponse{
			Status: *status,
		})
		if nil != err {
			l.WithError(err).Error(errMsgResponseStatus)
			http.Error(w, errMsgResponseStatus, http.StatusInternalServerError)
		}

		return
	}

	h.app.Monitor.StatusCacheMiss()

	status, err = h.app.Repo.GetStatus(token)
	if nil != err {
		l.WithError(err).Error(errMsgStatusGet)
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

	err = json.NewEncoder(w).Encode(PollResponse{
		Status: *status,
	})
	if nil != err {
		l.WithError(err).Error(errMsgResponseStatus)
		http.Error(w, errMsgResponseStatus, http.StatusInternalServerError)
	}
}
