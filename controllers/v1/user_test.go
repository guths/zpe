package v1_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/guths/zpe/config/setup"
	factory "github.com/guths/zpe/factory/factories"
	"github.com/guths/zpe/models"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*setup.TestDB, func()) {
	testDB := setup.GetTestDB()
	tx := testDB.DB.Begin()

	cleanup := func() {
		tx.Rollback()
		testDB.CleanupTestDB(t)
	}

	return testDB, cleanup
}

func createRoleAndUser(t *testing.T, db *gorm.DB, password string, roles []models.Role) models.User {
	model := models.NewRoleOrmer(db)
	var rolesToInsert []models.Role

	for _, r := range roles {
		role, err := model.InsertRole(r)
		require.NoError(t, err)
		rolesToInsert = append(rolesToInsert, role)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	f := factory.NewUserFactory()
	f.Roles = rolesToInsert
	u, _ := f.Create()
	u.Password = string(hashedPassword)

	orm := models.NewUserOrmer(db)
	createdUser, err := orm.InsertUser(*u)
	require.NoError(t, err)

	return createdUser
}

func performLoginRequest(email, password string) (*httptest.ResponseRecorder, map[string]interface{}) {
	reqBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	setup.RouterTest.ServeHTTP(w, req)

	var body map[string]interface{}
	json.NewDecoder(w.Body).Decode(&body)
	return w, body
}

func TestPOSTLogin(t *testing.T) {
	testDB, cleanup := setupTestDB(t)
	defer cleanup()

	tests := []struct {
		name          string
		password      string
		expectedCode  int
		expectedField string
		expectedValue string
	}{
		{
			name:          "successful login",
			password:      "password",
			expectedCode:  200,
			expectedField: "message",
			expectedValue: "user authenticated",
		},
		{
			name:          "failed login",
			password:      "wrongpassword",
			expectedCode:  401,
			expectedField: "error",
			expectedValue: "incorrect credentials",
		},
	}

	roleAdmin := models.Role{
		Level: 1,
		Name:  "admin",
	}

	user := createRoleAndUser(t, testDB.DB, "password", []models.Role{roleAdmin})

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w, body := performLoginRequest(user.Email, tc.password)

			require.Equal(t, tc.expectedCode, w.Code)
			require.Contains(t, body, tc.expectedField)
			require.Equal(t, tc.expectedValue, body[tc.expectedField])
		})
	}
}

func TestGETUser(t *testing.T) {
	testDB, cleanup := setupTestDB(t)
	defer cleanup()

	roleAdmin := models.Role{
		Level: 1,
		Name:  "admin",
	}

	user := createRoleAndUser(t, testDB.DB, "password", []models.Role{roleAdmin})

	w, body := performLoginRequest(user.Email, "password")

	require.Equal(t, 200, w.Code)

	token, _ := body["data"].(string)

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/user/%s", user.Email), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	w = httptest.NewRecorder()
	setup.RouterTest.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)

	var b map[string]interface{}
	json.NewDecoder(w.Body).Decode(&b)

	data := b["data"].(map[string]interface{})
	require.Equal(t, user.Email, data["email"])
	require.Equal(t, user.Username, data["username"])
	require.NotEmpty(t, data["created_at"])
	require.NotEmpty(t, data["updated_at"])
	require.Contains(t, data, "roles")

	roles := data["roles"].([]interface{})
	require.Len(t, roles, 1)
	role := roles[0].(map[string]interface{})
	require.Equal(t, "admin", role["Name"])
	require.Equal(t, float64(1), role["Level"])
}

func TestGETUserWithInvalidToken(t *testing.T) {
	testDB, cleanup := setupTestDB(t)
	defer cleanup()

	roleAdmin := models.Role{
		Level: 1,
		Name:  "admin",
	}

	user := createRoleAndUser(t, testDB.DB, "password", []models.Role{roleAdmin})

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/user/%s", user.Email), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "xxxxxx"))

	w := httptest.NewRecorder()
	setup.RouterTest.ServeHTTP(w, req)

	require.Equal(t, 401, w.Code)
}

func TestDELETEUser(t *testing.T) {
	testDB, cleanup := setupTestDB(t)
	defer cleanup()

	roleAdmin := models.Role{
		Level: 1,
		Name:  "admin",
	}

	roleModifier := models.Role{
		Level: 2,
		Name:  "modifier",
	}

	roleWatcher := models.Role{
		Level: 2,
		Name:  "watcher",
	}

	user := createRoleAndUser(t, testDB.DB, "password", []models.Role{roleAdmin})
	userToDelete := createRoleAndUser(t, testDB.DB, "password", []models.Role{roleModifier, roleWatcher})

	w, body := performLoginRequest(user.Email, "password")

	require.Equal(t, 200, w.Code)

	token, _ := body["data"].(string)

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/user/%s", userToDelete.Email), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	w = httptest.NewRecorder()
	setup.RouterTest.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)

	model := models.NewUserOrmer(testDB.DB)
	_, err := model.GetOneByEmail(userToDelete.Email)

	require.Error(t, err)
}

func TestDELETEUserWithInvalidRole(t *testing.T) {
	testDB, cleanup := setupTestDB(t)
	defer cleanup()

	roleAdmin := models.Role{
		Level: 1,
		Name:  "admin",
	}

	roleModifier := models.Role{
		Level: 2,
		Name:  "modifier",
	}

	user := createRoleAndUser(t, testDB.DB, "password", []models.Role{roleModifier})
	userToDelete := createRoleAndUser(t, testDB.DB, "password", []models.Role{roleAdmin})

	w, body := performLoginRequest(user.Email, "password")

	require.Equal(t, 200, w.Code)

	token, _ := body["data"].(string)

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/user/%s", userToDelete.Email), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	w = httptest.NewRecorder()
	setup.RouterTest.ServeHTTP(w, req)

	require.Equal(t, 401, w.Code)

	var b map[string]interface{}
	json.NewDecoder(w.Body).Decode(&b)

	message, _ := b["error"].(string)

	require.Equal(t, "permission denied", message)
}
