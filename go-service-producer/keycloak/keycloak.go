package keycloak

import "github.com/Nerzal/gocloak/v11"

type keycloak struct {
	gocloak      gocloak.GoCloak
	clientId     string
	clientSecret string
	realm        string
}

func NewKeycloak() *keycloak {
	return &keycloak{
		gocloak: gocloak.NewClient(
			"http://localhost:8080",
			gocloak.SetAuthAdminRealms("admin/realms"),
			gocloak.SetAuthRealms("realms")),
		clientId:     "test-go",
		clientSecret: "mJwvh9boCtMf0hMDXAhnwsQ7fIg2rL0h",
		realm:        "test-go",
	}
}
