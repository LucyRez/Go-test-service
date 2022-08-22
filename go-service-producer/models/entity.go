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

var entityList = []*Entity{
	{
		ID:           1,
		Name:         "First Entity",
		CreationTime: time.Now().UTC().String(),
	},
	{
		ID:           2,
		Name:         "Second Entity",
		CreationTime: time.Now().UTC().String(),
	},
}

func GetEntities() Entities {
	return entityList
}

func AddEntity(entity *Entity) {
	// Increase id by one and append new entity
	entity.ID = entityList[len(entityList)-1].ID + 1
	entity.CreationTime = time.Now().UTC().String()
	entityList = append(entityList, entity)
}

func GetEntityById(id int) *Entity {
	return entityList[id-1]
}
