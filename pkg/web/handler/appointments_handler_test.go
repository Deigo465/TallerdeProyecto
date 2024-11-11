package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/open-wm/blockehr/pkg/domain/entities"
	usecase "github.com/open-wm/blockehr/pkg/domain/usecases"
	mock "github.com/open-wm/blockehr/pkg/mocks"
)

func TestAddAppointment(t *testing.T) {
	// GIVEN
	appointmentRepo := mock.NewInMemoryAppointmentRepository()
	userRepo := mock.NewInMemoryUserRepository()
	profileRepo := mock.NewInMemoryProfileRepository()
	blockchain := mock.NewInMemoryBlockchain()
	uc := usecase.NewAppointmentUsecase(appointmentRepo, userRepo, profileRepo, blockchain)

	// staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)

	// handler
	appointmentHandler := &appointmentHandler{uc: uc}

	// appointment

	newAppointment := entities.NewFakeAppointment()

	// when
	appointmentJson, _ := json.Marshal(newAppointment)

	req, err := http.NewRequest("POST", "/appointments", bytes.NewBuffer(appointmentJson))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	appointmentHandler.Add(&staffActor, rw, req)
	resp := rw.Result()

	// THEN
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code to be 201, got %v", resp.StatusCode)
	}

	// Print the response body for debugging
	// Verify the response body
	var responseAppointment entities.Appointment
	if err := json.NewDecoder(rw.Body).Decode(&responseAppointment); err != nil {
		t.Fatal(err)
	}

}

func TestGetAllAppointments(t *testing.T) {
	// GIVEN
	appointmentRepo := mock.NewInMemoryAppointmentRepository()
	userRepo := mock.NewInMemoryUserRepository()
	profileRepo := mock.NewInMemoryProfileRepository()
	blockchain := mock.NewInMemoryBlockchain()
	uc := usecase.NewAppointmentUsecase(appointmentRepo, userRepo, profileRepo, blockchain)

	// staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)

	// handler
	appointmentHandler := &appointmentHandler{uc: uc}

	// when
	req, err := http.NewRequest("GET", "/appointmets", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	appointmentHandler.GetAll(&staffActor, rw, req)
	resp := rw.Result()

	// THEN
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be 200, got %v", resp.StatusCode)
	}
}

func TestGetByAppointmentId(t *testing.T) {
	// GIVEN
	appointmentRepo := mock.NewInMemoryAppointmentRepository()
	userRepo := mock.NewInMemoryUserRepository()
	profileRepo := mock.NewInMemoryProfileRepository()
	blockchain := mock.NewInMemoryBlockchain()
	uc := usecase.NewAppointmentUsecase(appointmentRepo, userRepo, profileRepo, blockchain)

	// staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)

	// handler
	appointmentHandler := &appointmentHandler{uc: uc}
	// create request with id 5
	req, err := http.NewRequest("GET", "/appointments/5", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	//capture response
	rw := httptest.NewRecorder()

	// router
	router := mux.NewRouter()
	router.HandleFunc("/appointments/{appointment_id}", func(rw http.ResponseWriter, req *http.Request) {
		appointmentHandler.GetById(&staffActor, rw, req)
	}).Methods("GET")

	//Handle request using router
	router.ServeHTTP(rw, req)
	resp := rw.Result()

	// THEN
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be 200, got %v", resp.StatusCode)
	}

	//check the body
	responseBody := rw.Body.String()
	var responseAppointment entities.Appointment
	if err := json.Unmarshal([]byte(responseBody), &responseAppointment); err != nil {
		t.Fatal(err)
	}

}

func TestUpdateAppointment(t *testing.T) {
	// GIVEN
	appointmentRepo := mock.NewInMemoryAppointmentRepository()
	userRepo := mock.NewInMemoryUserRepository()
	profileRepo := mock.NewInMemoryProfileRepository()
	blockchain := mock.NewInMemoryBlockchain()
	uc := usecase.NewAppointmentUsecase(appointmentRepo, userRepo, profileRepo, blockchain)

	// staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)

	//usecase
	appointmentHandler := &appointmentHandler{uc: uc}

	// appointment ID:7

	updateAppointment := entities.NewFakeAppointment()
	updateAppointment.ID = 7
	updateAppointment.Status = entities.CANCELED

	updateBody, err := json.Marshal(updateAppointment)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/appointments/7", bytes.NewBuffer(updateBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	//capture response

	rw := httptest.NewRecorder()

	// router
	router := mux.NewRouter()
	router.HandleFunc("/appointments/{appointment_id}", func(rw http.ResponseWriter, req *http.Request) {
		appointmentHandler.Update(&staffActor, rw, req)
	}).Methods("PUT")

	//Handle request using router
	router.ServeHTTP(rw, req)
	resp := rw.Result()

	// THEN
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be 200, got %v", resp.StatusCode)
	}
}
