package endpoints

import (
	model "go-service-producer/models"
	"net/http"
	"regexp"
	"strconv"
)

type Entity struct {
}

func NewEntityEndpoint() *Entity {
	return &Entity{}
}

func (entity *Entity) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		regex := regexp.MustCompile("/([0-9]+)")
		g := regex.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) == 0 {
			entity.getEntities(rw, r)
			return
		}

		// We can have only single id in URL
		if len(g) != 1 {
			http.Error(rw, "Invalid URL", http.StatusBadRequest)
			return
		}

		stringId := g[0][1]
		id, err := strconv.Atoi(stringId)

		if err != nil {
			http.Error(rw, "Invalid URL", http.StatusBadRequest)
			return
		}

		entity.getEntityById(id, rw, r)
		return
	}

	if r.Method == http.MethodPost {
		entity.addEntity(rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (e *Entity) getEntities(rw http.ResponseWriter, r *http.Request) {
	list := model.GetEntities()

	err := list.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to encode json", http.StatusInternalServerError)
	}
}

func (e *Entity) getEntityById(id int, rw http.ResponseWriter, r *http.Request) {
	entity := model.GetEntityById(id)
	err := entity.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to encode json", http.StatusInternalServerError)
	}
}

func (e *Entity) addEntity(rw http.ResponseWriter, r *http.Request) {
	entity := &model.Entity{}
	err := entity.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unable to decode json", http.StatusInternalServerError)
	}

	model.AddEntity(entity)
}
