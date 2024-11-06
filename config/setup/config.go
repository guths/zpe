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

//This file is used in tests to rollback all the data manipulated in DB
//this helps to avoid inconsistency between the tests
//The TestDB have a mutex to avoid a test that truncate a table in the same time that other test is using

var (
	RouterTest     *gin.Engine
	setupOnce      sync.Once
	testDBInstance *TestDB
	testDBOnce     sync.Once
)

type TestDB struct {
	DB *gorm.DB
	mu sync.Mutex
}

func GetTestDB() *TestDB {
	testDBOnce.Do(func() {
		testDBInstance = &TestDB{
			DB: config.DBManager.GetDB(),
		}
	})
	return testDBInstance
}

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

func (tdb *TestDB) CleanupTestDB(t *testing.T) {
	t.Helper()
	tdb.mu.Lock()
	defer tdb.mu.Unlock()

	tx := tdb.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	tx.Exec("SET FOREIGN_KEY_CHECKS = 0")
	defer tx.Exec("SET FOREIGN_KEY_CHECKS = 1")

	tables := []string{"users", "roles", "user_roles"}
	for _, table := range tables {
		if err := tx.Exec("TRUNCATE TABLE " + table).Error; err != nil {
			tx.Rollback()
			t.Fatalf("failed to truncate table %s: %v", table, err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		t.Fatalf("failed to commit cleanup: %v", err)
	}
}

func CreateAuthUser(roles []models.Role, tx *gorm.DB) models.User {
	f := factory.NewUserFactory()
	f.Roles = roles
	user, _ := f.Create()
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	orm := models.NewUserOrmer(tx)
	u, _ := orm.InsertUser(*user)
	return u
}
