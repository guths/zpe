package setup

import (
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guths/zpe/config"
	factory "github.com/guths/zpe/factory/factories"
	"github.com/guths/zpe/handlers"
	"github.com/guths/zpe/models"
	"github.com/guths/zpe/router"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var RouterTest *gin.Engine

var setupOnce sync.Once

func setupTest() {
	setupOnce.Do(func() {

		config.InitializeAppConfig()
		RouterTest = router.InitializeRouter()
		_ = handlers.InitializeHandler()
	})
}

func init() {
	setupTest()
}

func WithTransaction(t *testing.T, testFunc func(tx *gorm.DB) error) {
	db := config.DBManager.GetDB() // Ensure this returns a *gorm.DB instance
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			t.Fatalf("test panicked: %v", r)
		} else {
			tx.Rollback()
		}
	}()

	if err := testFunc(tx); err != nil {
		tx.Rollback()
		t.Fatalf("test failed: %v", err)
	} else {
		tx.Rollback() // Ensure rollback even if test passes
	}
}

func CreateAuthUser(roles []models.Role) models.User {
	f := factory.NewUserFactory()
	f.Roles = roles
	user, _ := f.Create()
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	orm := models.NewUserOrmer(config.DBManager.GetDB())
	u, _ := orm.InsertUser(*user)
	return u
}
