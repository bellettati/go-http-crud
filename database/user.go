package database

import (
	"errors"

	"github.com/google/uuid"
)

type Id uuid.UUID

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

var application = Application{
	data: make(map[Id]User),
} 

func FindAll() map[Id]User {
	return application.data
}

func FindById(id Id) *User {
	value, ok := application.data[id]
	if !ok {
		return nil
	}
	return &value
} 

func Insert(user User) {
	id := NewId();
	application.data[id] = user
}

func Update(id Id, user User) error {
	_, ok := application.data[id]

	if !ok {
		return errors.New("trying to update inexistent user")
	}

	application.data[id] = user

	return nil
}

func Delete(id Id) {
	delete(application.data, id)
}