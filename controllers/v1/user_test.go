package v1_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/guths/zpe/config"
	"github.com/guths/zpe/config/setup"
	"github.com/guths/zpe/models"
)

func TestPOSTLogin(t *testing.T) {
	// Begin a transaction
	tx := config.DBManager.GetDB().Begin()

	// Ensure the transaction is rolled back at the end of the test
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			t.Fatalf("test panicked: %v", r)
		} else {
			tx.Rollback()
		}
	}()

	// Perform database operations within the transaction
	r := models.NewRoleOrmer(tx)
	role, err := r.InsertRole(models.Role{
		Name:  "admin",
		Level: 1,
	})
	if err != nil {
		t.Fatalf("failed to insert role: %v", err)
	}

	u := setup.CreateAuthUser([]models.Role{role})
	reqBody := fmt.Sprintf(`{"email": "%s", "password": "password"}`, u.Email)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	setup.RouterTest.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("wrong status code: got %v want %v", w.Code, 200)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	message, ok := body["message"].(string)
	if !ok {
		t.Fatalf("message not found in response body")
	}

	if message != "user authenticated" {
		t.Fatalf("unexpected message: got %v want %v", message, "user authenticated")
	}
}
