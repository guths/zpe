package handlers

import (
	"github.com/guths/zpe/config"
	"github.com/guths/zpe/datatransfers"
	"github.com/guths/zpe/models"
	"gorm.io/gorm"
)

var Handler HandlerFunc

type HandlerFunc interface {
	AuthenticateUser(credentials datatransfers.UserLogin) (res datatransfers.Response)
	RegisterUser(credentials datatransfers.UserSignup) (res datatransfers.Response)
	DeleteUser(email string) (res datatransfers.Response)
	RetrieveUser(email string) (*models.User, error)
	ValidateUserRoles(userRoleLvl int, roles []string) bool
	// UpdateUser(id uint, user datatransfers.UserUpdate) (err error)
}

type Module struct {
	Db *dbEntity
}

type dbEntity struct {
	conn          *gorm.DB
	userOrmer     models.UserOrmer
	roleOrmer     models.RoleOrmer
	userRoleOrmer models.UserRoleOrmer
}

func InitializeHandler() (err error) {
	db := config.DBManager.GetDB()

	Handler = &Module{
		Db: &dbEntity{
			conn:          db,
			userOrmer:     models.NewUserOrmer(db),
			roleOrmer:     models.NewRoleOrmer(db),
			userRoleOrmer: models.NewUserRoleOrmer(db),
		},
	}

	return
}
