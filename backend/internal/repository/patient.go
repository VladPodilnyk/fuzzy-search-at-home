package repository

import (
	"context"
	"dicomviewer/internal/models"
	"fmt"

	"github.com/google/uuid"
)

type PatientRepository interface {
	Insert(ctx context.Context, name, birthDate string) (uint, error)
	IsExist(ctx context.Context, name string) (bool, error)
	Get(ctx context.Context, id uint) (models.UserData, error)
}

type InMemPatientRepository struct {
	store map[uint]models.UserData
}

func NewInMemPatientRepository() *InMemPatientRepository {
	return &InMemPatientRepository{store: make(map[uint]models.UserData)}
}

func (impl *InMemPatientRepository) Insert(ctx context.Context, name, birthDate string) (uint, error) {
	newId := uint(uuid.New().ID())
	impl.store[newId] = models.UserData{Name: name, BirthDate: birthDate}
	return newId, nil
}

func (impl *InMemPatientRepository) IsExist(ctx context.Context, name string) (bool, error) {
	for _, value := range impl.store {
		if value.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func (impl *InMemPatientRepository) Get(ctx context.Context, id uint) (models.UserData, error) {
	value, isOk := impl.store[id]
	if !isOk {
		return models.UserData{}, fmt.Errorf("cannot find a user with id %d", id)
	}
	return value, nil
}
