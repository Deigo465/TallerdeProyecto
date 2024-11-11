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
	moc "github.com/open-wm/blockehr/pkg/mocks"
)

func getProfileUC() usecase.ProfileUsecase {
	profileRepo := moc.NewInMemoryProfileRepository()
	appointmentsRepo := moc.NewInMemoryAppointmentRepository()
	userRepo := moc.NewInMemoryUserRepository()
	uc := usecase.NewProfileUsecase(profileRepo, appointmentsRepo, userRepo)
	return uc
}

func TestAddDoctor(t *testing.T) {
	// GIVEN
	ucProfile := getProfileUC()

	userRepo := moc.NewInMemoryUserRepository()
	ucUser := usecase.NewUserUsecase(userRepo)

	// staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)

	// handlers
	doctorHandler := &doctorHandler{usecaseProfile: ucProfile, usecaseUser: ucUser}

	// doctor

	newDoctor := entities.NewFakeProfile()
	newDoctor.Role = entities.DOCTOR

	// when
	doctorJson, _ := json.Marshal(newDoctor)

	req, err := http.NewRequest("POST", "/doctors", bytes.NewBuffer(doctorJson))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	doctorHandler.Add(&staffActor, rw, req)
	resp := rw.Result()

	// THEN
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code to be 201, got %v", resp.StatusCode)
	}

	// Verify the response body
	var responseDoctor entities.Profile
	if err := json.NewDecoder(rw.Body).Decode(&responseDoctor); err != nil {
		t.Fatal(err)
	}

}

func TestGetAllDoctors(t *testing.T) {
	// GIVEN
	ucProfile := getProfileUC()

	userRepo := moc.NewInMemoryUserRepository()
	ucUser := usecase.NewUserUsecase(userRepo)

	// staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)

	// handlers
	doctorHandler := &doctorHandler{usecaseProfile: ucProfile, usecaseUser: ucUser}

	// when
	req, err := http.NewRequest("GET", "/doctors", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	doctorHandler.GetAll(&staffActor, rw, req)
	resp := rw.Result()

	// THEN
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be 200, got %v", resp.StatusCode)
	}

}
func TestUpdateDoctor(t *testing.T) {
	// GIVEN
	ucProfile := getProfileUC()

	userRepo := moc.NewInMemoryUserRepository()
	ucUser := usecase.NewUserUsecase(userRepo)

	// staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)

	// handlers
	doctorHandler := &doctorHandler{usecaseProfile: ucProfile, usecaseUser: ucUser}

	// profile DOCTOR ID:2
	updateProfile := entities.NewFakeProfile()
	updateProfile.ID = 2
	updateProfile.Role = entities.DOCTOR
	updateProfile.FatherLastName = "BARRA"

	updateBody, err := json.Marshal(updateProfile)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/doctors/2", bytes.NewBuffer(updateBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	//capture response

	rw := httptest.NewRecorder()

	// router
	router := mux.NewRouter()
	router.HandleFunc("/doctors/{doctor_id}", func(rw http.ResponseWriter, req *http.Request) {
		doctorHandler.Update(&staffActor, rw, req)
	}).Methods("PUT")

	//Handle request using router
	router.ServeHTTP(rw, req)
	resp := rw.Result()

	// THEN
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be 200, got %v", resp.StatusCode)
	}
}
func TestGetDoctorById(t *testing.T) {
	// GIVEN
	ucProfile := getProfileUC()

	userRepo := moc.NewInMemoryUserRepository()
	ucUser := usecase.NewUserUsecase(userRepo)

	// staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)

	// handlers
	doctorHandler := &doctorHandler{usecaseProfile: ucProfile, usecaseUser: ucUser}
	// create request with id 2
	req, err := http.NewRequest("GET", "/doctors/2", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	//capture response
	rw := httptest.NewRecorder()

	// router
	router := mux.NewRouter()
	router.HandleFunc("/doctors/{doctor_id}", func(rw http.ResponseWriter, req *http.Request) {
		doctorHandler.GetById(&staffActor, rw, req)
	}).Methods("GET")

	//Handle request using router
	router.ServeHTTP(rw, req)
	resp := rw.Result()

	// THEN
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be 200, got %v", resp.StatusCode)
	}

	var responseDoctor entities.Profile
	if err := json.Unmarshal(rw.Body.Bytes(), &responseDoctor); err != nil {
		t.Fatal(err)
	}

}
