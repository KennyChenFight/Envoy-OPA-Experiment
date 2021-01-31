package store

import (
	model2 "github.com/KennyChenFight/Envoy-OPA-Experiment/k8s-envoy-opa/jwt-example/service/model"
	"github.com/gofiber/fiber/v2/utils"
	"sync"
)

type MemoryStore struct {
	people sync.Map
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) Save(person model2.Person) string {
	id := utils.UUIDv4()
	m.people.Store(id, person)
	return id
}

func (m *MemoryStore) GetAll() []model2.Person {
	var people []model2.Person
	m.people.Range(func(key, value interface{}) bool {
		id := key.(string)
		p := value.(model2.Person)
		p.ID = id
		people = append(people, p)
		return true
	})
	return people
}
