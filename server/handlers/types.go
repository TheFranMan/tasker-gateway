package handlers

type DeleteParams struct {
	ID int `json:"id"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
