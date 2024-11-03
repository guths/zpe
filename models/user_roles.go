package models

import "gorm.io/gorm"

type UserRole struct {
	UserID uint `gorm:"primaryKey;constraint:OnDelete:CASCADE;"`
	RoleID uint `gorm:"primaryKey;constraint:OnDelete:CASCADE;"`
	User   User `gorm:"constraint:OnDelete:CASCADE;"`
	Role   Role `gorm:"constraint:OnDelete:CASCADE;"`
}

type userRoleOrm struct {
	db *gorm.DB
}

type UserRoleOrmer interface{}

func NewUserRoleOrmer(db *gorm.DB) UserRoleOrmer {

	err := db.AutoMigrate(&UserRole{})

	if err != nil {
		panic(err)
	}
	return &userRoleOrm{db}
}
