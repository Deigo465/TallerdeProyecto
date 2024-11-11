package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/open-wm/blockehr/pkg/domain/entities"
	"github.com/open-wm/blockehr/pkg/web/handler"
)

func SetAPIRoutes(router *mux.Router, authHandler handler.AuthHandler, fileHandler handler.FileHandler,
	patientHandler handler.PatientHandler, doctorHandler handler.DoctorHandler, recordHandler handler.RecordHandler,
	appointmentsHandler handler.AppointmentHandler) {
	sr := roleRouter{router: router, middleware: authHandler.WithUser, roles: []entities.Role{entities.STAFF}}
	dr := roleRouter{router: router, middleware: authHandler.WithUser, roles: []entities.Role{entities.DOCTOR}}
	sdr := roleRouter{router: router, middleware: authHandler.WithUser, roles: []entities.Role{entities.STAFF, entities.DOCTOR}}

	sdr.reg(http.MethodGet, "/appointments", appointmentsHandler.GetAll)
	sr.reg(http.MethodPost, "/appointments", appointmentsHandler.Add)
	sr.reg(http.MethodPut, "/appointments/{appointment_id:[0-9]+}", appointmentsHandler.Update)

	//staff routes
	sr.reg(http.MethodGet, "/doctors", doctorHandler.GetAll)
	sr.reg(http.MethodGet, "/doctors/{doctor_id:[0-9]+}", doctorHandler.GetById)
	sr.reg(http.MethodPost, "/doctors", doctorHandler.Add)
	sr.reg(http.MethodPut, "/doctors/{doctor_id:[0-9]+}", doctorHandler.Update)
	sr.reg(http.MethodPost, "/doctors/{doctor_id:[0-9]+}/reset-password", doctorHandler.ResetDoctorPassword)

	// patients
	sr.reg(http.MethodGet, "/patients", patientHandler.GetAll)
	sr.reg(http.MethodGet, "/patients/{patient_id:[0-9]+}", patientHandler.GetById)
	sr.reg(http.MethodPost, "/patients", patientHandler.Add)
	sr.reg(http.MethodPut, "/patients/{patient_id:[0-9]+}", patientHandler.Update)
	sr.reg(http.MethodGet, "/patients/{patient_id:[0-9]+}/appointments", patientHandler.GetAllAppointmentsForPatient)

	// Doctors
	dr.reg(http.MethodPost, "/patients/{patient_id:[0-9]+}/records", recordHandler.Add)
	dr.reg(http.MethodPut, "/patients/{patient_id:[0-9]+}/records", recordHandler.Update)
	dr.reg(http.MethodGet, "/patients/{patient_id:[0-9]+}/records", recordHandler.GetAllForPatientId)
	dr.reg(http.MethodPost, "/appointments/{appointment_id:[0-9]+}/end", appointmentsHandler.EndAppointment)

	// files
	dr.reg(http.MethodPost, "/files", fileHandler.Add)
	// dr.reg(http.MethodPut, "/files", fileHandler.Update)
	// dr.reg(http.MethodDelete, "/files", fileHandler.Delete)
}
