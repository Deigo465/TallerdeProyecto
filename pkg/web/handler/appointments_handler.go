package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/open-wm/blockehr/pkg/domain/entities"
	usecase "github.com/open-wm/blockehr/pkg/domain/usecases"
)

type AppointmentHandler interface {
	Calendar(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	Add(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	Update(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	GetAll(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	GetById(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	EndAppointment(actor *entities.User, rw http.ResponseWriter, r *http.Request)
}

type appointmentHandler struct {
	uc usecase.AppointmentUsecase
}

func NewAppointmentHandler(ucAppointment usecase.AppointmentUsecase) AppointmentHandler {
	return &appointmentHandler{
		uc: ucAppointment,
	}
}

func withDefaultData(actor *entities.User, data map[string]interface{}) map[string]interface{} {
	if data == nil {
		data = make(map[string]interface{})
	}
	data["currentUser"] = actor
	return data
}

func (h *appointmentHandler) Calendar(actor *entities.User, rw http.ResponseWriter, r *http.Request) {
	if actor.Profile.Role == entities.DOCTOR {
		ExecLayout(rw, "appointments-doctor/index.html.tmpl", withDefaultData(actor, nil))
	}

	if actor.Profile.Role == entities.STAFF {
		ExecLayout(rw, "appointments/index.html.tmpl", withDefaultData(actor, nil))
	}
}

func (h *appointmentHandler) Add(actor *entities.User, rw http.ResponseWriter, r *http.Request) {

	var appointment entities.Appointment
	if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.uc.Add(actor, &appointment); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	JSON(rw, http.StatusCreated, "", appointment)

}

func (h *appointmentHandler) Update(actor *entities.User, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	appointmentID, err := strconv.Atoi(vars["appointment_id"]) // Convertir el ID de la cita a un entero
	if err != nil {
		http.Error(rw, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	var updatedAppointment entities.Appointment

	if err := json.NewDecoder(r.Body).Decode(&updatedAppointment); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.uc.UpdateStatus(actor, appointmentID, &updatedAppointment)

	if err != nil {
		log.Println(err)
		http.Error(rw, "Failed to update status", http.StatusInternalServerError)
		return
	}
	JSON(rw, http.StatusOK, "status updated", updatedAppointment)
}
func (h *appointmentHandler) GetById(actor *entities.User, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	appointmentID, err := strconv.Atoi(vars["appointment_id"])

	if err != nil {
		http.Error(rw, "Invalid appointment ID", http.StatusBadRequest)
		return
	}
	appointment, err := h.uc.GetById(actor, appointmentID)
	if err != nil {
		http.Error(rw, "Failed to returned appointment", http.StatusInternalServerError)
		return
	}
	JSON(rw, http.StatusOK, "appointment returned", appointment)
}

func (h *appointmentHandler) GetAll(actor *entities.User, rw http.ResponseWriter, r *http.Request) {
	appointments, err := h.uc.GetAll(actor)
	if err != nil {
		log.Println(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(rw, http.StatusOK, "success", appointments)
}

// FinishedAppointment implements AppointmentHandler.
func (h *appointmentHandler) EndAppointment(actor *entities.User, rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	appointmentID, err := strconv.Atoi(vars["appointment_id"]) // Convertir el ID de la cita a un entero
	if err != nil {
		http.Error(rw, "Invalid appointment ID", http.StatusBadRequest)
		return
	}
	err = h.uc.EndAppointment(actor, appointmentID)

	if err != nil {
		log.Println(err)
		http.Error(rw, "Failed to update status", http.StatusInternalServerError)
		return
	}
	JSON(rw, http.StatusOK, "status updated", nil)

}
