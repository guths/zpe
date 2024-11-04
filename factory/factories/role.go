package factory

import (
	"github.com/go-faker/faker/v4"
	"github.com/guths/zpe/models"
)

type roleFactory struct{}

func NewRoleFactory() *roleFactory {
	return &roleFactory{}
}

func (f *roleFactory) Create() (*models.Role, error) {
	r := models.Role{}

	if err := faker.FakeData(&r); err != nil {
		return nil, err
	}

	return &r, nil
}
