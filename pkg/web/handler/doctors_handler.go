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

type DoctorHandler interface {
	List(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	GetAll(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	Add(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	Update(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	GetById(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	ResetDoctorPassword(actor *entities.User, rw http.ResponseWriter, r *http.Request)
}

type doctorHandler struct {
	usecaseProfile usecase.ProfileUsecase
	usecaseUser    usecase.UserUsecase
}

func NewDoctorHandler(ucProfile usecase.ProfileUsecase, ucUser usecase.UserUsecase) DoctorHandler {
	return &doctorHandler{
		usecaseProfile: ucProfile,
		usecaseUser:    ucUser,
	}
}

func (h *doctorHandler) List(actor *entities.User, rw http.ResponseWriter, r *http.Request) {
	ExecLayout(rw, "doctors/index.html.tmpl", withDefaultData(actor, nil))
}

// Perfil con Email (vamos a reutilizarlo)
type ProfileWithEmail struct {
	entities.Profile
	Email string `json:"email"`
}

func (h *doctorHandler) GetAll(actor *entities.User, rw http.ResponseWriter, r *http.Request) {
	doctors, err := h.usecaseProfile.GetAllDoctors(actor)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := []ProfileWithEmail{}
	for _, doctor := range doctors {
		doctorWithEmail := ProfileWithEmail{
			Profile: entities.Profile{
				ID:             doctor.ID,
				FirstName:      doctor.FirstName,
				MotherLastName: doctor.MotherLastName,
				FatherLastName: doctor.FatherLastName,
				DocumentNumber: doctor.DocumentNumber,
				Phone:          doctor.Phone,
				Gender:         doctor.Gender,
				DateOfBirth:    doctor.DateOfBirth,
				Cmp:            doctor.Cmp,
				Specialty:      doctor.Specialty,
			},
			Email: doctor.ContactEmail,
		}

		responseData = append(responseData, doctorWithEmail)
	}

	//response
	JSON(rw, http.StatusOK, "ok", responseData)
}
func (h *doctorHandler) Add(actor *entities.User, rw http.ResponseWriter, r *http.Request) {
	//TODO: add user and password?

	var doctor entities.Profile
	if err := json.NewDecoder(r.Body).Decode(&doctor); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.usecaseProfile.AddDoctor(actor, &doctor); err != nil {
		log.Println(err)
		if errors.Is(err, usecase.ErrBadDoctorData) {
			JSON(rw, http.StatusBadRequest, err.Error(), nil)
		} else if errors.Is(err, usecase.ErrDocumentAlreadyExists) {
			JSON(rw, http.StatusBadRequest, err.Error(), nil)
		} else {
			JSON(rw, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}
	JSON(rw, http.StatusCreated, "", doctor)

}
func (h *doctorHandler) Update(actor *entities.User, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	doctorId, err := strconv.Atoi(vars["doctor_id"]) // Convertir el ID de la cita a un entero
	if err != nil {
		http.Error(rw, "Invalid profile ID", http.StatusBadRequest)
		return
	}

	var updatedDoctor entities.Profile

	if err := json.NewDecoder(r.Body).Decode(&updatedDoctor); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.usecaseProfile.Update(actor, doctorId, &updatedDoctor)

	if err != nil {
		log.Println(err)
		http.Error(rw, "Failed to update profile", http.StatusInternalServerError)
		return
	}
	JSON(rw, http.StatusOK, "Doctor updated", updatedDoctor)
}
func (h *doctorHandler) GetById(actor *entities.User, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	doctorID, err := strconv.Atoi(vars["doctor_id"])

	if err != nil {
		http.Error(rw, "Invalid doctor ID", http.StatusBadRequest)
		return
	}
	doctor, err := h.usecaseProfile.GetById(actor, doctorID)
	if err != nil {
		http.Error(rw, "Failed to returned doctor", http.StatusInternalServerError)
		return
	}
	JSON(rw, http.StatusOK, "doctor returned", doctor)
}

func (h *doctorHandler) ResetDoctorPassword(actor *entities.User, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	doctorID, err := strconv.Atoi(vars["doctor_id"])

	if err != nil {
		http.Error(rw, "Invalid doctor ID", http.StatusBadRequest)
		return
	}
	doctor, err := h.usecaseUser.ResetDoctorPassword(actor, doctorID)
	if err != nil {
		http.Error(rw, "Failed to reset password for doctor", http.StatusInternalServerError)
		return
	}
	JSON(rw, http.StatusOK, "doctor password updated", map[string]any{
		"password": doctor,
	})
}
