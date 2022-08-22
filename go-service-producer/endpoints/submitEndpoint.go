package endpoints

import (
	"context"
	"encoding/json"
	"go-service-producer/models"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
)

type Submit struct {
}

func NewSubmitEndpoint() *Submit {
	return &Submit{}
}

func (submit *Submit) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		submit.submitEntities(rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (submit *Submit) submitEntities(rw http.ResponseWriter, r *http.Request) {
	connection, _ := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "topic_test", 0)
	connection.SetWriteDeadline(time.Now().Local().Add(time.Second * 10))

	entity := &models.Entity{}
	err := entity.FromJSON(r.Body)

	res, err := json.Marshal(entity)

	if err != nil {
		http.Error(rw, "Unable to encode json", http.StatusInternalServerError)
	}

	connection.WriteMessages(kafka.Message{Value: res})
}
