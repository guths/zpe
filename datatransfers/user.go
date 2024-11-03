package datatransfers

import "time"

type UserSignup struct {
	Username string   `json:"username" binding:"required"`
	Email    string   `json:"email" binding:"required"`
	Password string   `json:"password" binding:"required"`
	Roles    []string `json:"roles" binding:"required"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserInfo struct {
	Username  string    `json:"username"`
	Email     string    `json:"email" uri:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserUpdate struct {
	Username string   `json:"username" binding:"-"`
	Email    string   `json:"email" uri:"email"`
	Roles    []string `json:"roles" binding:"-"`
}
