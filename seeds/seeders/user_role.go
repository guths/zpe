package seeders

import (
	"github.com/guths/zpe/models"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserRoleSeeder(db *gorm.DB) {
	role1 := models.Role{
		Name:  "admin",
		Level: 1,
	}

	db.Create(&role1)

	role2 := models.Role{
		Name:  "modifier",
		Level: 2,
	}

	db.Create(&role2)

	role3 := models.Role{
		Name:  "watcher",
		Level: 3,
	}

	db.Create(&role3)

	var defaultPass []byte
	var err error

	defaultPass, err = bcrypt.GenerateFromPassword([]byte(viper.GetString("DEFAULT_PASS")), bcrypt.DefaultCost)
	if err != nil {
		panic("failed to generate passwords")
	}

	adminUser := models.User{
		Username: "admin",
		Email:    "admin@admin.com",
		Password: string(defaultPass),
	}

	adminUser.Roles = []models.Role{role1, role2, role3}

	modifierUser := models.User{
		Username: "modifier",
		Email:    "modifier@modifier.com",
		Password: string(defaultPass),
	}

	modifierUser.Roles = []models.Role{role2, role3}

	watcherUser := models.User{
		Username: "watcher",
		Email:    "watcher@watcher.com",
		Password: string(defaultPass),
	}

	watcherUser.Roles = []models.Role{role3}

	db.Create(&adminUser)
	db.Create(&modifierUser)
	db.Create(&watcherUser)
}
