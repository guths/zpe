package handlers

import (
	"fmt"

	"github.com/guths/zpe/datatransfers"
	"github.com/guths/zpe/models"
)

func (m *Module) ValidateUserRoles(userRoleLvl int, roles []string) bool {
	if len(roles) == 0 {
		return true
	}

	mRoles, err := m.Db.roleOrmer.GetManyByName(roles)

	if err != nil {
		return false
	}

	maxLvlRole := models.GetMaxRoleLvl(mRoles)

	return maxLvlRole >= userRoleLvl
}

func (m *Module) RetrieveUser(email string) (*models.User, error) {
	user, err := m.Db.userOrmer.GetOneByEmail(email)

	if err != nil {
		return nil, err
	}

	return &user, err
}

func (m *Module) DeleteUser(email string) datatransfers.Response {
	err := m.Db.userOrmer.DeleteOneByEmail(email)

	if err != nil {
		return datatransfers.Response{
			Code:  404,
			Error: "can not find user",
		}
	}

	return datatransfers.Response{
		Message: "user sucessfully deleted",
		Code:    200,
	}
}

func (m *Module) UpdateUser(id uint, user datatransfers.UserUpdate) (*models.User, error) {
	roles, err := m.Db.roleOrmer.GetManyByName(user.Roles)

	if err != nil {
		return nil, fmt.Errorf("error getting the roles")
	}

	updatedUser, err := m.Db.userOrmer.UpdateUser(models.User{
		ID:       id,
		Username: user.Username,
		Roles:    roles,
	})

	if err != nil {
		return nil, fmt.Errorf("error updating the user")
	}

	return &updatedUser, nil
}
