package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type userOrm struct {
	db *gorm.DB
}

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(100);size:100" json:"username" faker:"username"`
	Email     string    `gorm:"unique" json:"email" faker:"email"`
	Password  string    `json:"-"`
	Roles     []Role    `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE;" faker:"-" json:"roles"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type UserOrmer interface {
	GetOneByID(id uint) (user User, err error)
	GetOneByEmail(email string) (User, error)
	InsertUser(user User) (u User, err error)
	UpdateUser(user User) (User, error)
	DeleteOneByEmail(email string) error
}

func NewUserOrmer(db *gorm.DB) UserOrmer {

	err := db.AutoMigrate(&User{})

	if err != nil {
		panic(err)
	}
	return &userOrm{db}
}

func (o *userOrm) GetOneByID(id uint) (User, error) {
	var user User
	result := o.db.Model(&User{}).Preload("Roles").Where("id = ?", id).First(&user)
	return user, result.Error
}

func (o *userOrm) GetOneByEmail(email string) (User, error) {
	var user User
	result := o.db.Model(&User{}).Preload("Roles").Where("email = ?", email).First(&user)
	return user, result.Error
}

func (o *userOrm) InsertUser(user User) (u User, err error) {
	result := o.db.Model(&User{}).Create(&user)
	return user, result.Error
}

func (o *userOrm) UpdateUser(user User) (User, error) {
	fmt.Println(user.Roles)
	if len(user.Roles) > 0 {
		o.db.Model(&User{}).Model(&user).Association("Roles").Replace(user.Roles)
	}

	result := o.db.Model(&User{}).Model(&user).Updates(&user)
	return user, result.Error
}

func (o *userOrm) DeleteOneByEmail(email string) error {
	result := o.db.Model(&User{}).Where("email = ?", email).Delete(&email)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
