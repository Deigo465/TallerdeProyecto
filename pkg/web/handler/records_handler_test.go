package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/open-wm/blockehr/pkg/domain/entities"
	usecase "github.com/open-wm/blockehr/pkg/domain/usecases"
	moc "github.com/open-wm/blockehr/pkg/mocks"
)

func getRecordHandler() *recordHandler {
	recordRepo := moc.NewInMemoryRecordRepository()
	profileRepo := moc.NewInMemoryProfileRepository()
	fileRepo := moc.NewInMemoryFileRepository()
	blockchain := moc.NewInMemoryBlockchain()
	ucRecord := usecase.NewRecordUsecase(recordRepo, profileRepo, fileRepo, blockchain)
	return &recordHandler{uc: ucRecord}
}

func TestAdd(t *testing.T) {
	// GIVEN
	recordHandler := getRecordHandler()
	// doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)

	// record
	newRecord := entities.NewRecord(1, "Juan Pérez acude a la consulta médica quejándose de dolor en el pecho y dificultad para respirar.",
		"01-06-2003", "29-04-2004", 2, 3)

	// Check if running from tests
	file := "./storage/test/file.txt"
	if strings.Contains(os.Args[0], ".test") {
		file = "../../../storage/test/file.txt"
	}
	values := map[string]io.Reader{
		"file[]": mustOpen(file), // lets assume its this file
		"body":   strings.NewReader(newRecord.Body),
	}
	b, contentType, _ := createFormData(values, "text/plain")

	req, err := http.NewRequest("POST", "/patients/1/records/", &b)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"patient_id": "1"})
	req.Header.Set("Content-Type", contentType)
	rw := httptest.NewRecorder()

	// WHEN
	recordHandler.Add(&doctorActor, rw, req)
	resp := rw.Result()

	// THEN
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code to be 201, got %v", resp.StatusCode)
		t.Fatal(rw.Body)
	}

	// Verify the response body
	var responseRecord entities.Record
	if err := json.NewDecoder(rw.Body).Decode(&responseRecord); err != nil {
		t.Fatal(err)
	}

}
func TestGetAllForPatientId(t *testing.T) {
	// GIVEN
	recordHandler := getRecordHandler()

	// doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)

	// when
	req, err := http.NewRequest("GET", "/records/patient/1", nil)
	req = mux.SetURLVars(req, map[string]string{"patient_id": "1"})
	if err != nil {
		t.Fatal(err)
	}

	// capture response
	rw := httptest.NewRecorder()

	// router
	router := mux.NewRouter()
	recordHandler.GetAllForPatientId(&doctorActor, rw, req)
	router.ServeHTTP(rw, req)
	resp := rw.Result()

	// THEN
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be 200, got %v", resp.StatusCode)
	}

	// var responseRecords []entities.Record

	// TODO: Add patient to repo to be able to test this
	// err = parseResponse(resp, &responseRecords)
	// assert.Nil(t, err)

	// for _, record := range responseRecords {
	// 	assert.Equal(t, 1, record.PatientId)
	// }
}
func TestGetByRecordId(t *testing.T) {
	// GIVEN
	recordHandler := getRecordHandler()

	// Doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)

	// create request with id 2
	req, err := http.NewRequest("GET", "/records/2", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	//capture response
	rw := httptest.NewRecorder()

	// router
	router := mux.NewRouter()
	router.HandleFunc("/records/{record_id}", func(rw http.ResponseWriter, req *http.Request) {
		recordHandler.GetByRecordId(&doctorActor, rw, req)
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

	var responseRecord entities.Record
	if err := json.Unmarshal([]byte(responseBody), &responseRecord); err != nil {
		t.Fatal(err)
	}

}
