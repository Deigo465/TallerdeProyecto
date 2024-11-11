package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	usecase "github.com/open-wm/blockehr/pkg/domain/usecases"
	moc "github.com/open-wm/blockehr/pkg/mocks"
)

func TestLoginForm(t *testing.T) {
	// GIVEN
	userRepo := moc.NewInMemoryUserRepository()
	sessRepo := moc.NewInMemorySessionRepository()
	profileRepo := moc.NewInMemoryProfileRepository()
	ucAuth := usecase.NewAuthUsecase(userRepo, sessRepo, profileRepo)
	authHandler := NewAuthHandler(ucAuth)
	SetBasePath("../../web/views/")

	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	rw := httptest.NewRecorder()

	// WHEN
	authHandler.LoginForm(rw, req)

	// THEN
	resp := rw.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be 200, got %v", resp.StatusCode)
	}

	body := rw.Body.String()

	if body == "" {
		t.Errorf("Expected body to not be empty")
	}

	if !strings.Contains(body, "Ingresa") {
		t.Errorf("Expected body to include 'Ingresa', got %s", body)
	}
}

func TestWithUser(t *testing.T) {

	// test that it returns a Handler function that when executed it has 3 cases

	// GIVEN
	userRepo := moc.NewInMemoryUserRepository()
	sessRepo := moc.NewInMemorySessionRepository()
	profileRepo := moc.NewInMemoryProfileRepository()
	ucAuth := usecase.NewAuthUsecase(userRepo, sessRepo, profileRepo)
	authHandler := NewAuthHandler(ucAuth)
	testHandler := func(user *entities.User, rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("OK"))
	}

	// add staffCookie
	staffCookie := http.Cookie{
		Name:  usecase.SESSION_ID,
		Value: "STAFF",
	}

	testCases := []struct {
		name           string
		cookie         http.Cookie
		statusCode     int
		expectRedirect bool
	}{
		{
			name:           "Valid",
			cookie:         staffCookie,
			statusCode:     http.StatusOK,
			expectRedirect: false,
		},
		{
			name:           "No cookie",
			cookie:         http.Cookie{},
			statusCode:     http.StatusFound,
			expectRedirect: true,
		},
		{
			name: "Expired cookie",
			cookie: http.Cookie{
				Name:    usecase.SESSION_ID,
				Value:   "EXPIRED",
				Expires: staffCookie.Expires.Add(-1),
			},
			statusCode:     http.StatusFound,
			expectRedirect: true,
		},
		{
			name: "No roles",
			cookie: http.Cookie{
				Name:  usecase.SESSION_ID,
				Value: "BADKEY",
			},
			statusCode:     http.StatusFound,
			expectRedirect: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// GIVEN
			req, err := http.NewRequest("GET", "/", nil)
			req.AddCookie(&tc.cookie)
			if err != nil {
				t.Fatal(err)
			}
			rw := httptest.NewRecorder()

			// WHEN
			authHandler.WithUser(testHandler)(rw, req)

			// THEN
			resp := rw.Result()
			if resp.StatusCode != tc.statusCode {
				t.Errorf("Expected status code to be %d, got %v", tc.statusCode, resp.StatusCode)
			}

			body := rw.Body.String()

			if body == "" {
				t.Errorf("Expected body to not be empty")
			}
			if tc.expectRedirect {
				if !strings.Contains(body, `<a href="/login?next=/">Found</a>`) {
					t.Errorf("Expected body to include 'Found', got %s", body)
				}
			} else {
				if !strings.Contains(body, "OK") {
					t.Errorf("Expected body to include 'OK', got %s", body)
				}
			}
		})
	}
}

func TestLoginJSON(t *testing.T)     {}
func TestLoginFormDate(t *testing.T) {}
func TestLogout(t *testing.T)        {}
