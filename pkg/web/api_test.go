package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/open-wm/blockehr/pkg/domain/entities"
	usecase "github.com/open-wm/blockehr/pkg/domain/usecases"
	"github.com/open-wm/blockehr/pkg/web/handler"
)

func setCookie(req *http.Request, staff string) {
	cookie := http.Cookie{
		Name:  usecase.SESSION_ID,
		Value: staff,
	}
	req.AddCookie(&cookie)
}

func getRouter() *mux.Router {
	router := mux.NewRouter()
	defineRoutes(router, true)
	return router
}

func TestGetAllDoctors(t *testing.T) {
	// GIVEN
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/doctors", nil)

	router := getRouter()
	setCookie(req, "STAFF")

	// WHEN
	router.ServeHTTP(w, req)

	// THEN
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	body := w.Body.String()
	if !strings.Contains(body, "data") {
		t.Errorf("Expected body to contain 'data' but got %v", body)
	}

	// check that the header contains application/json
	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected Content-Type to be application/json, but got %v", contentType)
	}

	type response struct {
		Data   []handler.ProfileWithEmail `json:"data"`
		Status int                        `json:"status"`
	}

	res := response{}
	// unmarshall json from the data

	err := json.Unmarshal([]byte(w.Body.Bytes()), &res)
	if err != nil {
		t.Errorf("Failed to unmarshall json: %v body: %v", err, w.Body.String())
	}

}

func TestAddDoctor(t *testing.T) {
	// GIVEN
	w := httptest.NewRecorder()
	validDoctor := handler.ProfileWithEmail{
		// Profile fields
		Profile: entities.Profile{
			FirstName:      "Rodrigo",
			MotherLastName: "Aguilar",
			FatherLastName: "Gonzales",
			Cmp:            "123546",
			DateOfBirth:    "2005-02-27",
			DocumentNumber: "71323285",
			Phone:          "987654321",
			Gender:         "Masculino",
			Specialty:      "Cardiologia",
		},
		Email: "rodrigo@example.org",
	}
	jsonDoctor, _ := json.Marshal(validDoctor)
	reader := bytes.NewReader(jsonDoctor)

	req, _ := http.NewRequest("POST", "/api/v1/doctors", reader)

	router := mux.NewRouter()
	defineRoutes(router, true)
	setCookie(req, "STAFF")

	// WHEN
	router.ServeHTTP(w, req)

	// THEN
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, resp.StatusCode)
	}

	body := w.Body.String()
	if !strings.Contains(body, "data") {
		t.Errorf("Expected body to contain 'data' but got %v", body)
	}

	// check that the header contains application/json
	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected Content-Type to be application/json, but got %v", contentType)
	}

	type response struct {
		Data   handler.ProfileWithEmail `json:"data"`
		Status int                      `json:"status"`
	}

	res := response{}
	// unmarshall json from the data

	err := json.Unmarshal([]byte(w.Body.Bytes()), &res)
	if err != nil {
		t.Errorf("Failed to unmarshall json: %v body: %v", err, w.Body.String())
	}

}

func TestUpdateDoctor(t *testing.T) {
	// GIVEN
	w := httptest.NewRecorder()
	validDoctor := handler.ProfileWithEmail{
		// Profile fields
		Email: "rodrigo@example.org",
		Profile: entities.Profile{
			FirstName:      "Rodrigo",
			MotherLastName: "Aguilar",
			FatherLastName: "Gonzales",
			Cmp:            "123546",
			DateOfBirth:    "2005-02-27",
			DocumentNumber: "71323285",
			Phone:          "987654321",
			Gender:         "Masculino",
			Specialty:      "Cardiologia",
		},
	}
	jsonDoctor, _ := json.Marshal(validDoctor)
	reader := bytes.NewReader(jsonDoctor)

	req, _ := http.NewRequest("PUT", "/api/v1/doctors/1", reader)

	router := getRouter()
	setCookie(req, "STAFF")

	// WHEN
	router.ServeHTTP(w, req)

	// THEN
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	body := w.Body.String()
	if !strings.Contains(body, "data") {
		t.Errorf("Expected body to contain 'data' but got %v", body)
	}

	// check that the header contains application/json
	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected Content-Type to be application/json, but got %v", contentType)
	}

	type response struct {
		Data   handler.ProfileWithEmail `json:"data"`
		Status int                      `json:"status"`
	}

	res := response{}
	// unmarshall json from the data

	err := json.Unmarshal([]byte(w.Body.Bytes()), &res)
	if err != nil {
		t.Errorf("Failed to unmarshall json: %v body: %v", err, w.Body.String())
	}

}

// /////PATIENT////////
func TestGetAllPatients(t *testing.T) {
	// GIVEN
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/patients", nil)
	setCookie(req, "STAFF")

	router := getRouter()

	// WHEN
	router.ServeHTTP(w, req)

	// THEN
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	} else {
		fmt.Printf("Status code: %d\n", resp.StatusCode)
	}

	body := w.Body.String()

	if !strings.Contains(body, "data") {
		t.Errorf("Expected body to contain 'data' but got %v", body)
	}

	// check that the header contains application/json
	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected Content-Type to be application/json, but got %v", contentType)
	}

	type response struct {
		Data   []handler.ProfileWithEmail `json:"data"`
		Status int                        `json:"status"`
	}

	res := response{}
	// unmarshall json from the data

	err := json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Errorf("Failed to unmarshall json: %v body: %v", err, body)
	}

}

func TestAddPatient(t *testing.T) {
	// GIVEN
	w := httptest.NewRecorder()
	validPatient := handler.ProfileWithEmail{
		// Profile fields
		Email: "rodrigo@example.org",
		Profile: entities.Profile{
			FirstName:      "Rodrigo",
			MotherLastName: "Aguilar",
			FatherLastName: "Gonzales",
			Cmp:            "",
			DateOfBirth:    "2005-02-27",
			DocumentNumber: "71323285",
			Phone:          "987654321",
			Gender:         "Masculino",
			Specialty:      "",
		},
	}
	jsonPatient, _ := json.Marshal(validPatient)
	reader := bytes.NewReader(jsonPatient)
	req, _ := http.NewRequest("POST", "/api/v1/patients", reader)
	setCookie(req, "STAFF")

	router := getRouter()

	// WHEN
	router.ServeHTTP(w, req)

	// THEN
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, resp.StatusCode)
	}

	body := w.Body.String()
	if !strings.Contains(body, "data") {
		t.Errorf("Expected body to contain 'data' but got %v", body)
	}

	// check that the header contains application/json
	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected Content-Type to be application/json, but got %v", contentType)
	}

	type response struct {
		Data   handler.ProfileWithEmail `json:"data"`
		Status int                      `json:"status"`
	}

	res := response{}
	// unmarshall json from the data

	err := json.Unmarshal([]byte(w.Body.Bytes()), &res)
	if err != nil {
		t.Errorf("Failed to unmarshall json: %v body: %v", err, w.Body.String())
	}

}

func TestUpdatePatient(t *testing.T) {
	// GIVEN
	w := httptest.NewRecorder()
	validPatient := handler.ProfileWithEmail{
		// Profile fields
		Email: "rodrigo@example.org",
		Profile: entities.Profile{
			FirstName:      "Rodrigo",
			MotherLastName: "Aguilar",
			FatherLastName: "Gonzales",
			DateOfBirth:    "2005-02-27",
			Gender:         "Masculino",
			DocumentNumber: "71323285",
			Phone:          "987654321",
		},
	}
	jsonPatient, _ := json.Marshal(validPatient)
	reader := bytes.NewReader(jsonPatient)
	// id = 3 is patient
	req, _ := http.NewRequest("PUT", "/api/v1/patients/3", reader)
	setCookie(req, "STAFF")

	router := getRouter()

	// WHEN
	router.ServeHTTP(w, req)
	// THEN
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	body := w.Body.String()
	if !strings.Contains(body, "data") {
		t.Errorf("Expected body to contain 'data' but got %v", body)
	}

	// check that the header contains application/json
	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected Content-Type to be application/json, but got %v", contentType)
	}

	type response struct {
		Data   handler.ProfileWithEmail `json:"data"`
		Status int                      `json:"status"`
	}

	res := response{}
	// unmarshall json from the data

	err := json.Unmarshal([]byte(w.Body.Bytes()), &res)
	if err != nil {
		t.Errorf("Failed to unmarshall json: %v body: %v", err, w.Body.String())
	}

}

func TestUpdatePatientFails(t *testing.T) {
	// GIVEN
	w := httptest.NewRecorder()
	validPatient := handler.ProfileWithEmail{
		Email: "rodrigo@example.org",
		Profile: entities.Profile{
			// Profile fields
			FirstName:      "Rodrigo",
			MotherLastName: "Aguilar",
			FatherLastName: "Gonzales",
			DateOfBirth:    "2005-02-27",
			DocumentNumber: "71323285",
			Phone:          "987654321",
		},
	}
	jsonPatient, _ := json.Marshal(validPatient)
	reader := bytes.NewReader(jsonPatient)
	// id = 3 is patient
	req, _ := http.NewRequest("PUT", "/api/v1/patients/3", reader)
	setCookie(req, "DOCTOR")

	router := getRouter()

	// WHEN
	router.ServeHTTP(w, req)
	// THEN
	resp := w.Result()
	if resp.StatusCode != http.StatusFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	// check if redirects to /login
	if location := resp.Header.Get("Location"); !strings.Contains(location, "/login") {
		t.Errorf("Expected Location to be /login, but got %v", location)
	}
}
