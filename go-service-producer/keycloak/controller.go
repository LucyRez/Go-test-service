package keycloak

import (
	"context"
	"encoding/json"
	"net/http"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

type controller struct {
	keycloak *keycloak
}

func newController(keycloak *keycloak) *controller {
	return &controller{
		keycloak: keycloak,
	}
}

func (ctr *controller) login(w http.ResponseWriter, r *http.Request) {
	request := &loginRequest{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jwt, err := ctr.keycloak.gocloak.Login(context.Background(),
		ctr.keycloak.clientId,
		ctr.keycloak.clientSecret,
		ctr.keycloak.realm,
		request.Username,
		request.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	response := &loginResponse{
		AccessToken:  jwt.AccessToken,
		RefreshToken: jwt.RefreshToken,
		ExpiresIn:    jwt.ExpiresIn,
	}

	responseEncoded, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseEncoded)
}
