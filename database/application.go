package database

import (
	"errors"

	"github.com/google/uuid"
)

type Id uuid.UUID

func (id Id) String() string {
	return uuid.UUID(id).String()
}

func NewId() Id {
	return Id(uuid.New())
}

type User struct {
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Biography string `json:"biography"`
}

type Application struct {
	data map[Id]User
}

func New() Application {
	return Application{
		data: make(map[Id]User),
	}
} 

func (application *Application) FindAll() map[string]User {
	users := make(map[string]User)
	for id, user := range application.data {
		users[id.String()] = user
	}
	return users
}

func (application *Application) FindById(id Id) *User {
	value, ok := application.data[id]
	if !ok {
		return nil
	}
	return &value
} 

func (application *Application) Insert(user User) {
	id := NewId();
	application.data[id] = user
}

func (application *Application) Update(id Id, user User) error {
	_, ok := application.data[id]

	if !ok {
		return errors.New("trying to update inexistent user")
	}

	application.data[id] = user

	return nil
}

func (application *Application) Delete(id Id) {
	delete(application.data, id)
}