package endpoints

import (
	model "go-service-receiver/models"
	"net/http"
)

type Received struct {
}

func NewReceivedEndpoint() *Received {
	return &Received{}
}

func (received *Received) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		received.getReceivedEntities(rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (received *Received) getReceivedEntities(rw http.ResponseWriter, r *http.Request) {
	list := model.GetReceivedEntities()

	err := list.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to encode json", http.StatusInternalServerError)
	}
}
