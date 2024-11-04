package factory

import (
	"github.com/go-faker/faker/v4"
	"github.com/guths/zpe/models"
)

type userFactory struct {
	Roles []models.Role
}

func NewUserFactory() *userFactory {
	return &userFactory{}
}

func (f *userFactory) Create() (*models.User, error) {
	u := models.User{
		Roles: f.Roles,
	}

	if err := faker.FakeData(&u); err != nil {
		return nil, err
	}

	return &u, nil
}
