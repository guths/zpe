package models

import (
	"math"

	"gorm.io/gorm"
)

type Role struct {
	ID    uint   `gorm:"primaryKey" json:"-"`
	Name  string `gorm:"unique;not null"`
	Level uint   `gorm:"not null"`
}

type roleOrm struct {
	db *gorm.DB
}

type RoleOrmer interface {
	GetManyByName(names []string) ([]Role, error)
}

func NewRoleOrmer(db *gorm.DB) RoleOrmer {

	err := db.AutoMigrate(&Role{})

	if err != nil {
		panic(err)
	}
	return &roleOrm{db}
}

func (o *roleOrm) GetManyByName(names []string) ([]Role, error) {
	var roles []Role

	result := o.db.Where("name IN ?", names).Find(&roles)

	return roles, result.Error
}

func GetMaxRoleLvl(roles []Role) int {
	maxLvl := math.MaxInt

	for _, r := range roles {
		if r.Level < uint(maxLvl) {
			maxLvl = int(r.Level)
		}
	}

	return maxLvl
}
