package main

import (
	"go-backend/models"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserGetAdmin(t *testing.T) {
	setup()
	defer teardown()

	token, _ := generateTestJWT(1)

	req, err := http.NewRequest("GET", "/api/admin", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(jwtMiddleware(s.GetUserAdminRoute))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		log.Printf("err %v", rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}
func TestUserGetAdminFailure(t *testing.T) {
	setup()
	defer teardown()

	token, _ := generateTestJWT(2)

	req, err := http.NewRequest("GET", "/api/admin", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(jwtMiddleware(s.GetUserAdminRoute))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		log.Printf("err %v", rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestGetUserSuccess(t *testing.T) {
	setup()
	defer teardown()

	token, _ := generateTestJWT(1)
	req, err := http.NewRequest("GET", "/api/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(jwtMiddleware(admin(s.GetUserRoute)))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	var user models.User
	parseJsonResponse(t, rr.Body.Bytes(), &user)
	if user.ID != 1 {
		t.Errorf("handler returned wrong user id, got %v want %v", user.ID, 1)
	}

}
func TestGetUserUnauthorized(t *testing.T) {
	setup()
	defer teardown()

	token, _ := generateTestJWT(2)
	req, err := http.NewRequest("GET", "/api/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(jwtMiddleware(admin(s.GetUserRoute)))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

}
func TestGetUserBadInput(t *testing.T) {
	setup()
	defer teardown()

	token, _ := generateTestJWT(1)
	req, err := http.NewRequest("GET", "/api/users/-1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.SetPathValue("id", "-1")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(jwtMiddleware(admin(s.GetUserRoute)))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGetCurrentUserSuccess(t *testing.T) {

	setup()
	defer teardown()

	token, _ := generateTestJWT(3)
	req, err := http.NewRequest("GET", "/api/current", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(jwtMiddleware(s.GetCurrentUserRoute))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	var user models.User
	parseJsonResponse(t, rr.Body.Bytes(), &user)
	if user.ID != 3 {
		t.Errorf("handler returned wrong user id, got %v want %v", user.ID, 3)
	}
}

func TestGetUserSubscriptionSuccess(t *testing.T) {

	setup()
	defer teardown()

	token, _ := generateTestJWT(1)
	req, err := http.NewRequest("GET", "/api/users/1/subscription", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(jwtMiddleware(s.GetUserSubscriptionRoute))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	var userSub models.UserSubscription
	parseJsonResponse(t, rr.Body.Bytes(), &userSub)
	if userSub.ID != 1 {
		t.Errorf("handler returned wrong user id, got %v want %v", userSub.ID, 1)
	}
}
func TestGetUserSubscriptionUnauthorized(t *testing.T) {
	setup()
	defer teardown()

	token, _ := generateTestJWT(3)
	req, err := http.NewRequest("GET", "/api/users/1/subscription", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(jwtMiddleware(s.GetUserSubscriptionRoute))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestGetUsersRouteSuccess(t *testing.T) {
	setup()
	defer teardown()

	token, _ := generateTestJWT(1)
	req, err := http.NewRequest("GET", "/api/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(jwtMiddleware(s.GetUsersRoute))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	var users []models.User
	parseJsonResponse(t, rr.Body.Bytes(), &users)
	if len(users) != 10 {
		t.Errorf("handler returned wrong number of users, got %v want %v", len(users), 10)
	}
}
