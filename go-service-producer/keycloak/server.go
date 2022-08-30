package keycloak

import (
	"context"
	"fmt"
	"go-service-producer/endpoints"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type serverHTTP struct {
	server *http.Server
}

func NewServer(host, port string, keycloak *keycloak) *serverHTTP {

	router := mux.NewRouter()

	noAuthRouter := router.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return r.Header.Get("Authorization") == ""
	}).Subrouter()

	authRouter := router.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return true
	}).Subrouter()

	authRouterAdmin := router.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return true
	}).Subrouter()

	controller := newController(keycloak)

	noAuthRouter.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controller.login(w, r)
	}).Methods("POST")

	entityEndpoint := endpoints.NewEntityEndpoint()

	authRouter.HandleFunc("/entity", func(w http.ResponseWriter, r *http.Request) {
		entityEndpoint.ServeHTTP(w, r)
	}).Methods("GET")

	authRouterAdmin.HandleFunc("/entity", func(w http.ResponseWriter, r *http.Request) {
		entityEndpoint.ServeHTTP(w, r)
	}).Methods("GET", "POST")

	middleware := newMiddleware(keycloak)
	authRouter.Use(middleware.verifyToken)
	authRouterAdmin.Use(middleware.verifyTokenAdmin)

	log := log.New(os.Stdout, "server-producer", log.LstdFlags)
	server := &serverHTTP{
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", host, port),
			Handler:      router,
			ErrorLog:     log,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 8 * time.Second,
			IdleTimeout:  30 * time.Minute,
		},
	}

	return server
}

func (server *serverHTTP) Listen() error {
	return server.server.ListenAndServe()
}

func (server *serverHTTP) Shutdown() {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.server.Shutdown(ctx)
}
