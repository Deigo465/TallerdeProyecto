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

func TestAddFile(t *testing.T) {
	// GIVEN
	fileRepo := moc.NewInMemoryFileRepository()
	ucFile := usecase.NewFileUsecase(fileRepo)
	recordRepo := moc.NewInMemoryRecordRepository()
	profileRepo := moc.NewInMemoryProfileRepository()
	blockchain := moc.NewInMemoryBlockchain()
	ucRecord := usecase.NewRecordUsecase(recordRepo, profileRepo, fileRepo, blockchain)

	// doctor
	doctorProfile := entities.NewFakeProfile()
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)
	doctorProfile.Role = entities.DOCTOR

	// handler
	fileHandler := &fileHandler{fileUsecase: ucFile, recordUsecase: ucRecord}

	// Check if running from tests
	file := "./storage/test/file.txt"
	if strings.Contains(os.Args[0], ".test") {
		file = "../../../storage/test/file.txt"
	}
	values := map[string]io.Reader{
		"file[]": mustOpen(file), // lets assume its this file
		"other":  strings.NewReader("hello world!"),
	}
	b, contentType, _ := createFormData(values, "text/plain")

	// when
	req, err := http.NewRequest("POST", "/files", &b)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", contentType)

	rw := httptest.NewRecorder()
	fileHandler.Add(&doctorActor, rw, req)
	resp := rw.Result()
	// THEN
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code to be 201, got %v", resp.StatusCode)
	}

	// Print the response body for debugging

	// Verify the response body
	var responseFile entities.File
	if err := json.NewDecoder(rw.Body).Decode(&responseFile); err != nil {
		t.Fatal(err)
	}

}
func TestGetByFileId(t *testing.T) {
	// GIVEN
	fileRepo := moc.NewInMemoryFileRepository()
	ucFile := usecase.NewFileUsecase(fileRepo)
	recordRepo := moc.NewInMemoryRecordRepository()
	profileRepo := moc.NewInMemoryProfileRepository()
	blockchain := moc.NewInMemoryBlockchain()
	ucRecord := usecase.NewRecordUsecase(recordRepo, profileRepo, fileRepo, blockchain)

	// doctor
	doctorProfile := entities.NewFakeProfile()
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)
	doctorProfile.Role = entities.DOCTOR

	// handler
	fileHandler := &fileHandler{fileUsecase: ucFile, recordUsecase: ucRecord}

	// create request with id 2
	req, err := http.NewRequest("GET", "/files/2", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	//capture response
	rw := httptest.NewRecorder()

	// router
	router := mux.NewRouter()
	router.HandleFunc("/files/{file_id}", func(rw http.ResponseWriter, req *http.Request) {
		fileHandler.GetProfileById(&doctorActor, rw, req)
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

	var responseFile entities.File
	if err := json.Unmarshal([]byte(responseBody), &responseFile); err != nil {
		t.Fatal(err)
	}

}

func TestGetFileByRecordId(t *testing.T) {
	// GIVEN
	fileRepo := moc.NewInMemoryFileRepository()
	ucFile := usecase.NewFileUsecase(fileRepo)
	recordRepo := moc.NewInMemoryRecordRepository()
	profileRepo := moc.NewInMemoryProfileRepository()
	blockchain := moc.NewInMemoryBlockchain()
	ucRecord := usecase.NewRecordUsecase(recordRepo, profileRepo, fileRepo, blockchain)

	// doctor
	doctorProfile := entities.NewFakeProfile()
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)
	doctorProfile.Role = entities.DOCTOR

	// handler
	fileHandler := &fileHandler{fileUsecase: ucFile, recordUsecase: ucRecord}

	// create request with id 3
	req, err := http.NewRequest("GET", "/records/3", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	//capture response
	rw := httptest.NewRecorder()

	// router
	router := mux.NewRouter()
	router.HandleFunc("/records/{record_id}", func(rw http.ResponseWriter, req *http.Request) {
		fileHandler.GetByRecordId(&doctorActor, rw, req)
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

	var responseFile entities.File
	if err := json.Unmarshal([]byte(responseBody), &responseFile); err != nil {
		t.Fatal(err)
	}

}
