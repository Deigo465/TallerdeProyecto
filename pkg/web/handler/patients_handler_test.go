package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/open-wm/blockehr/pkg/domain/entities"
)

func TestAddPatient(t *testing.T) {
	// GIVEN
	ucProfile := getProfileUC()

	// staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)

	// handlers
	patientHandler := &patientHandler{usecaseProfile: ucProfile}

	// patient

	newPatient := entities.NewFakeProfile()
	newPatient.Role = entities.PATIENT

	// when
	patientJson, _ := json.Marshal(newPatient)

	req, err := http.NewRequest("POST", "/patients", bytes.NewBuffer(patientJson))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	patientHandler.Add(&staffActor, rw, req)
	resp := rw.Result()

	// THEN
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code to be 201, got %v", resp.StatusCode)
	}

	// Print the response body for debugging

	// Verify the response body
	var responsePatient entities.Profile
	if err := json.NewDecoder(rw.Body).Decode(&responsePatient); err != nil {
		t.Fatal(err)
	}

}
func TestGetAllPatients(t *testing.T) {
	// GIVEN
	ucProfile := getProfileUC()

	// staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)

	// handlers
	patientHandler := &patientHandler{usecaseProfile: ucProfile}

	// when
	req, err := http.NewRequest("GET", "/patients", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	patientHandler.GetAll(&staffActor, rw, req)
	resp := rw.Result()

	// THEN
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be 200, got %v", resp.StatusCode)
	}
}

func TestUpdatePatient(t *testing.T) {
	// GIVEN
	ucProfile := getProfileUC()

	// staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.ID = 3 // 3 is patient
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)

	// handlers
	patientHandler := &patientHandler{usecaseProfile: ucProfile}

	// profile PATIENT ID:9
	updateProfile := entities.NewFakeProfile()
	updateProfile.ID = 3 // id 3 is patient
	updateProfile.Role = entities.PATIENT
	updateProfile.FatherLastName = "BARRA"
	updateProfile.Gender = "Masculino"

	updateBody, err := json.Marshal(updateProfile)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/patients/3", bytes.NewBuffer(updateBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	//capture response

	rw := httptest.NewRecorder()

	// router
	router := mux.NewRouter()
	router.HandleFunc("/patients/{patient_id}", func(rw http.ResponseWriter, req *http.Request) {
		patientHandler.Update(&staffActor, rw, req)
	}).Methods("PUT")

	//Handle request using router
	router.ServeHTTP(rw, req)
	resp := rw.Result()

	// THEN
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be 200, got %v", resp.StatusCode)
	}
}

func TestGetPatientById(t *testing.T) {
	// GIVEN
	ucProfile := getProfileUC()

	// staff
	patientProfile := entities.NewFakeProfile()
	patientProfile.ID = 3
	patientProfile.Role = entities.STAFF
	staffActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, patientProfile.ID, patientProfile)

	// handlers
	patientHandler := &patientHandler{usecaseProfile: ucProfile}

	// create request with id 4
	req, err := http.NewRequest("GET", "/patients/3", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	//capture response
	rw := httptest.NewRecorder()

	// router
	router := mux.NewRouter()
	router.HandleFunc("/patients/{patient_id}", func(rw http.ResponseWriter, req *http.Request) {
		patientHandler.GetById(&staffActor, rw, req)
	}).Methods("GET")

	//Handle request using router
	router.ServeHTTP(rw, req)
	resp := rw.Result()

	// THEN
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be 200, got %v", resp.StatusCode)
	}

	var responsePatient entities.Profile
	if err := json.Unmarshal(rw.Body.Bytes(), &responsePatient); err != nil {
		t.Fatal(err)
	}

}
