package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/open-wm/blockehr/pkg/domain/entities"
	usecase "github.com/open-wm/blockehr/pkg/domain/usecases"
)

type PatientHandler interface {
	List(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	GetAll(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	Add(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	Update(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	GetById(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	GetAllAppointmentsForPatient(actor *entities.User, rw http.ResponseWriter, r *http.Request)
}

type patientHandler struct {
	usecaseProfile usecase.ProfileUsecase
}

func NewPatientHandler(ucProfile usecase.ProfileUsecase) PatientHandler {
	return &patientHandler{
		usecaseProfile: ucProfile,
	}
}

func (h *patientHandler) List(actor *entities.User, rw http.ResponseWriter, r *http.Request) {
	patients, err := h.usecaseProfile.GetAllPatients(actor)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	ExecLayout(rw, "patients/index.html.tmpl", withDefaultData(actor, map[string]any{
		"patients": patients,
	}))
}

func (h *patientHandler) GetAll(actor *entities.User, rw http.ResponseWriter, r *http.Request) {

	patients, err := h.usecaseProfile.GetAllPatients(actor)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	responseData := []ProfileWithEmail{}
	for _, patient := range patients {
		patientWithEmail := ProfileWithEmail{
			Profile: entities.Profile{
				ID:             patient.ID,
				FirstName:      patient.FirstName,
				MotherLastName: patient.MotherLastName,
				FatherLastName: patient.FatherLastName,
				DocumentNumber: patient.DocumentNumber,
				Phone:          patient.Phone,
				DateOfBirth:    patient.DateOfBirth,
				Cmp:            patient.Cmp,
				Gender:         patient.Gender,
				Specialty:      patient.Specialty,
			},
			Email: patient.ContactEmail,
		}

		responseData = append(responseData, patientWithEmail)
	}

	//response
	JSON(rw, http.StatusOK, "ok", responseData)
}

func (h *patientHandler) Add(actor *entities.User, rw http.ResponseWriter, r *http.Request) {

	var patient entities.Profile
	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.usecaseProfile.AddPatient(actor, &patient); err != nil {
		log.Println(err)
		if errors.Is(err, usecase.ErrDocumentAlreadyExists) {
			JSON(rw, http.StatusBadRequest, err.Error(), nil)
		} else {
			JSON(rw, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}
	JSON(rw, http.StatusCreated, "", patient)

}

func (h *patientHandler) Update(actor *entities.User, rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	patient, err := strconv.Atoi(vars["patient_id"]) // Convertir el ID de la cita a un entero
	if err != nil {
		http.Error(rw, "Invalid profile ID", http.StatusBadRequest)
		return
	}

	var updatedPatient entities.Profile

	if err := json.NewDecoder(r.Body).Decode(&updatedPatient); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.usecaseProfile.Update(actor, patient, &updatedPatient)

	if err != nil {
		log.Println(err)
		http.Error(rw, "Failed to update status", http.StatusInternalServerError)
		return
	}
	JSON(rw, http.StatusOK, "Profile updated", updatedPatient)
}

func (h *patientHandler) GetById(actor *entities.User, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	patientID, err := strconv.Atoi(vars["patient_id"]) //TODO

	if err != nil {
		http.Error(rw, "Invalid patient ID", http.StatusBadRequest)
		return
	}
	patient, err := h.usecaseProfile.GetById(actor, patientID)
	if err != nil {
		http.Error(rw, "Failed to returned patient", http.StatusInternalServerError)
		return
	}
	JSON(rw, http.StatusOK, "patient returned", patient)
}

func (h *patientHandler) GetAllAppointmentsForPatient(actor *entities.User, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	patientID, err := strconv.Atoi(vars["patient_id"]) //TODO

	appointments, err := h.usecaseProfile.GetAllAppointmentsForPatient(actor, patientID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	//response
	JSON(rw, http.StatusOK, "ok", appointments)
}
