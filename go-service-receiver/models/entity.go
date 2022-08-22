package models

import (
	"encoding/json"
	"io"
	"time"
)

type Entity struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	CreationTime string `json:"creation-time"`
	ReceiveTime  string `json:"receive-time"`
}

func (e *Entity) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(e)
}

type Entities []*Entity

func (e *Entities) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(e)
}

func (e *Entity) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(e)
}

var entityList = []*Entity{}

func GetReceivedEntities() Entities {
	return entityList
}

func ReceiveEntity(entity *Entity) {
	// Increase id by one and append new entity
	entity.ReceiveTime = time.Now().UTC().String()
	entityList = append(entityList, entity)
}
