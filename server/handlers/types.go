package handlers

import "github.com/TheFranMan/tasker-common/types"

type DeleteParams struct {
	ID int `json:"id"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type pollResponse struct {
	Status types.RequestStatusString `json:"status"`
}
