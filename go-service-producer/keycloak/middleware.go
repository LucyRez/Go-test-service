package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type keycloakMiddleware struct {
	keycloak *keycloak
}

type Token struct {
	Claims Claims `json:"Claims,omitempty"`
}

type Claims struct {
	ResourceAccess client `json:"resource_access,omitempty"`
	JTI            string `json:"jti,omitempty"`
}

type client struct {
	TestGo clientRoles `json:"test-go,omitempty"`
}

type clientRoles struct {
	Roles []string `json:"roles,omitempty"`
}

func newMiddleware(keycloak *keycloak) *keycloakMiddleware {
	return &keycloakMiddleware{keycloak: keycloak}
}

func (auth *keycloakMiddleware) getBearerToken(token string) string {
	return strings.Replace(token, "Bearer ", "", 1)
}

func (auth *keycloakMiddleware) verifyUser(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {

		if !auth.verifyToken(w, r, "") {
			http.Error(w, "Could not authorize token for user", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}

func (auth *keycloakMiddleware) verifyAdmin(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {

		if !auth.verifyToken(w, r, "ADMIN") {
			http.Error(w, "Could not authorize token for user", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}

func (auth *keycloakMiddleware) verifyToken(w http.ResponseWriter, r *http.Request, role string) bool {
	token := r.Header.Get("Authorization")

	if token == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return false
	}

	token = auth.getBearerToken(token)

	if token == "" {
		http.Error(w, "Bearer Token missing", http.StatusUnauthorized)
		return false
	}

	// verifying access token
	result, err := auth.keycloak.gocloak.RetrospectToken(context.Background(), token, auth.keycloak.clientId, auth.keycloak.clientSecret, auth.keycloak.realm)

	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid or malformed token: %s", err.Error()), http.StatusUnauthorized)
		return false
	}

	jwt, _, err := auth.keycloak.gocloak.DecodeAccessToken(context.Background(), token, auth.keycloak.realm)

	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid or malformed token: %s", err.Error()), http.StatusUnauthorized)
		return false
	}

	jwtj, _ := json.Marshal(jwt)
	// fmt.Printf("token: %v\n", string(jwtj))

	var claims Token
	json.Unmarshal(jwtj, &claims)
	// fmt.Printf("roles:  %v\n", claims.Claims.ResourceAccess.TestGo.Roles)

	if !*result.Active {
		http.Error(w, "Invalid or expired Token", http.StatusUnauthorized)
		return false
	}

	if role != "" {
		var foundRole = false

		for _, userRole := range claims.Claims.ResourceAccess.TestGo.Roles {
			if userRole == role {
				foundRole = true
			}
		}

		if !foundRole {
			http.Error(w, "User doesn't have the right role to authorize", http.StatusUnauthorized)
			return false
		}
	}

	return true
}
