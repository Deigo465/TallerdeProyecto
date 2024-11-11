package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/open-wm/blockehr/pkg/domain/entities"
	usecase "github.com/open-wm/blockehr/pkg/domain/usecases"
)

type RecordHandler interface {
	Add(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	Update(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	GetAllForPatientId(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	GetByRecordId(actor *entities.User, rw http.ResponseWriter, r *http.Request)
}
type recordHandler struct {
	uc usecase.RecordUsecase
}

func NewRecordHandler(ucRecord usecase.RecordUsecase) RecordHandler {
	return &recordHandler{
		uc: ucRecord,
	}
}

func (h *recordHandler) Add(actor *entities.User, rw http.ResponseWriter, r *http.Request) {
	// get patient id
	vars := mux.Vars(r)
	patientIDstr := vars["patient_id"]
	patientId, err := strconv.Atoi(patientIDstr) // convert string to int
	if err != nil {
		http.Error(rw, "Invalid patient ID: "+patientIDstr, http.StatusBadRequest)
		return
	}

	r.ParseMultipartForm(10 << 20) // read 10 MB
	multipart := r.MultipartForm
	if multipart == nil {
		log.Println("No file found in multipart")
		JSON(rw, http.StatusBadRequest, "No file found", nil)
		return
	}
	files := multipart.File["file[]"]
	if len(files) == 0 {
		log.Println("No file found in file[]")
		JSON(rw, http.StatusBadRequest, "No file found", nil)
		return
	}

	fileModels := []*entities.File{}
	for _, file := range files {
		// Save file
		fileModel, err := StoreFile(file)
		if err != nil {
			log.Println("error storing", err)
			JSON(rw, http.StatusInternalServerError, err.Error(), nil)
		}
		fileModels = append(fileModels, fileModel)
	}

	// get record from request MultiPartForm
	body := multipart.Value["body"]
	if len(body) == 0 {
		log.Println("No body found in multipart", err)
		JSON(rw, http.StatusBadRequest, "No body found", nil)
		return
	}
	var record entities.Record = entities.NewRecord(
		0, body[0],
		time.Now().Local().String(), time.Now().Local().String(),
		patientId, actor.Profile.ID,
	)
	record.Files = fileModels
	err = h.uc.Add(actor, &record)
	if err != nil {
		log.Println("Error addinng record", err)
		JSON(rw, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	JSON(rw, http.StatusCreated, "", record)
}
func (h *recordHandler) Update(actor *entities.User, rw http.ResponseWriter, r *http.Request) {

	var record entities.Record
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.uc.Add(actor, &record); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	JSON(rw, http.StatusCreated, "", record)

}

func (h *recordHandler) GetAllForPatientId(actor *entities.User, rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	patientID, err := strconv.Atoi(vars["patient_id"])
	if err != nil {
		http.Error(rw, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	// Obtener registros para el paciente
	records, err := h.uc.GetAllForPatient(actor, patientID)
	if err != nil {
		JSON(rw, http.StatusUnauthorized, err.Error(), nil)
		return
	}
	JSON(rw, http.StatusOK, "records returned", records)
}

func (h *recordHandler) GetByRecordId(actor *entities.User, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	recordId, err := strconv.Atoi(vars["record_id"]) //TODO

	if err != nil {
		http.Error(rw, "Invalid record ID", http.StatusBadRequest)
		return
	}
	record, err := h.uc.GetById(actor, recordId)
	if err != nil {
		http.Error(rw, "Failed to returned record", http.StatusInternalServerError)
		return
	}
	JSON(rw, http.StatusOK, "record returned", record)
}
