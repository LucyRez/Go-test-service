package keycloak

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v11"
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

type registrationRequest struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Enabled   bool   `json:"enabled"`
	Username  string `json:"username"`
	Password  string `json:"password"`
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

func (ctr *controller) register(w http.ResponseWriter, r *http.Request) {
	request := &registrationRequest{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token := r.Header.Get("Authorization")

	if token == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	token = ctr.getBearerToken(token)

	user := gocloak.User{
		FirstName: gocloak.StringP(request.FirstName),
		LastName:  gocloak.StringP(request.LastName),
		Email:     gocloak.StringP(request.Email),
		Enabled:   gocloak.BoolP(request.Enabled),
		Username:  gocloak.StringP(request.Username),
	}

	userId, err := ctr.keycloak.gocloak.CreateUser(context.Background(), token, "test-go", user)

	if err != nil {
		http.Error(w, "Error: Could not create user", http.StatusBadRequest)
		return
	}

	passwordError := ctr.keycloak.gocloak.SetPassword(context.Background(), token, userId, "test-go", request.Password, false)

	if passwordError != nil {
		http.Error(w, "Error: Could not set new user's password", http.StatusBadRequest)
		return
	}

}

func (ctr *controller) getBearerToken(token string) string {
	return strings.Replace(token, "Bearer ", "", 1)
}
