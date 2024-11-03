package handlers

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/guths/zpe/config"
	"github.com/guths/zpe/constants"
	"github.com/guths/zpe/datatransfers"
	"github.com/guths/zpe/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (m *Module) AuthenticateUser(credentials datatransfers.UserLogin) (res datatransfers.Response) {
	var user models.User
	var err error
	if user, err = m.Db.userOrmer.GetOneByEmail(credentials.Email); err != nil {
		return datatransfers.Response{
			Code:  401,
			Error: "incorrect credentials",
		}
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return datatransfers.Response{
			Code:  401,
			Error: "incorrect credentials",
		}
	}

	token, err := generateToken(user)
	if err != nil {
		return datatransfers.Response{
			Code:  401,
			Error: "incorrect credentials",
		}
	}

	return datatransfers.Response{
		Code:    200,
		Message: "user authenticated",
		Data:    token,
	}
}

func generateToken(user models.User) (string, error) {
	now := time.Now()
	expiry := time.Now().Add(constants.AuthenticationTimeout)

	var roles []string

	for _, r := range user.Roles {
		roles = append(roles, r.Name)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":        user.ID,
		"Roles":     roles,
		"ExpiresAt": expiry.Unix(),
		"IssuedAt":  now.Unix(),
	})

	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

func (m *Module) RegisterUser(credentials datatransfers.UserSignup) datatransfers.Response {
	var hashedPassword []byte
	var err error

	if hashedPassword, err = bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost); err != nil {
		return datatransfers.Response{
			Code:  400,
			Error: "failed hashing password",
		}
	}

	_, err = m.Db.userOrmer.GetOneByEmail(credentials.Email)

	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return datatransfers.Response{
			Code:  400,
			Error: "error getting user",
		}
	}

	if err == nil {
		return datatransfers.Response{
			Code:  409,
			Error: "user already exists",
		}
	}

	roles, err := m.Db.roleOrmer.GetManyByName(credentials.Roles)

	if err != nil {
		return datatransfers.Response{
			Code:  400,
			Error: "the provided roles are invalid",
		}
	}

	if len(roles) == 0 {
		return datatransfers.Response{
			Code:  404,
			Error: "roles not found",
		}
	}

	user, err := m.Db.userOrmer.InsertUser(models.User{
		Username: credentials.Username,
		Email:    credentials.Email,
		Password: string(hashedPassword),
		Roles:    roles,
	})

	if err != nil {
		return datatransfers.Response{
			Code:  400,
			Error: err.Error(),
		}
	}

	return datatransfers.Response{
		Code:    201,
		Data:    user,
		Message: "user created successfully",
	}
}
