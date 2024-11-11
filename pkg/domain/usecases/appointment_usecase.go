package usecase

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	"github.com/open-wm/blockehr/pkg/interfaces"
)

// Interface
type AppointmentUsecase interface {
	Add(actor *entities.User, appointment *entities.Appointment) error
	UpdateStatus(actor *entities.User, id int, appointment *entities.Appointment) error
	GetAll(actor *entities.User) ([]*AppointmentWithDoctor, error)
	GetById(actor *entities.User, id int) (*entities.Appointment, error)
	EndAppointment(actor *entities.User, id int) error
}

// impl
type appointmentUsecase struct {
	appointmentRepo interfaces.AppointmentRepository
	userRepo        interfaces.UserRepository
	profileRepo     interfaces.ProfileRepository
	client          interfaces.BlockchainClient
}

func NewAppointmentUsecase(
	repo interfaces.AppointmentRepository,
	userRepo interfaces.UserRepository,
	profileRepo interfaces.ProfileRepository,
	client interfaces.BlockchainClient,
) AppointmentUsecase {
	return &appointmentUsecase{
		appointmentRepo: repo,
		userRepo:        userRepo,
		profileRepo:     profileRepo,
		client:          client,
	}
}

func (uc *appointmentUsecase) Add(actor *entities.User, newAppointment *entities.Appointment) error {
	// check if role of actor is staff
	if actor.Profile.Role != entities.STAFF {
		return errors.New("only staff can add appointments")
	}

	if newAppointment == nil {
		return errors.New("expecting appointment to not be empty")
	}

	// check that patient can't have 2 concurrent appointments
	appointments, err := uc.appointmentRepo.GetByPatientID(newAppointment.PatientId)
	if err != nil {
		return errors.New("cant find patient appointments")
	}

	for _, appointment := range appointments {
		if appointment.StartsAt.Year() == newAppointment.StartsAt.Year() &&
			appointment.StartsAt.Month() == newAppointment.StartsAt.Month() &&
			appointment.StartsAt.Day() == newAppointment.StartsAt.Day() &&
			appointment.StartsAt.Hour() == newAppointment.StartsAt.Hour() {
			// cant create
			return errors.New("can't create concurrent appointments")
		}
	}

	// check if appointment is in the future or empty
	// if appointment.StartsAt.Before(time.Now()) {
	// 	return errors.New("appointment date must be in the future")
	// }

	//validate
	if newAppointment.Specialty == "" {
		return errors.New("expecting specialty to not be empty")
	}

	if newAppointment.Status < 0 || newAppointment.Status > 4 {
		return errors.New("expecting status to be valid")
	}

	if newAppointment.DoctorId == 0 {
		return errors.New("expecting doctor id to not be empty")
	}

	// get doctor user
	doctor, err := uc.userRepo.GetUserByProfileID(newAppointment.DoctorId)
	if err != nil {
		return err
	}
	// get patient user
	patient, err := uc.profileRepo.GetById(newAppointment.PatientId)
	if err != nil {
		log.Println("valid patient must be passed")
		return err
	}

	// check if doctor belongs to the same health center as the actor
	if actor.HealthCenterId != doctor.HealthCenterId {
		return errors.New("doctor does not belong to the same health center as the actor")
	}
	if newAppointment.Description == "" {
		return errors.New("expecting description to not be empty")
	}

	//if all ok add appointment
	newAppointment.DoctorId = doctor.ID
	if err := uc.appointmentRepo.Add(newAppointment); err != nil {
		return err
	}

	if newAppointment.Status == entities.STARTED {
		message := fmt.Sprintf("Otorgando permiso al doctor %s para ver la información del paciente: %d", doctor.Email, patient.ID)
		uc.addPermissionToBlockchain(newAppointment.DoctorId, patient.ID, "granted", message)
	}
	return nil

}

func (uc *appointmentUsecase) addPermissionToBlockchain(doctorId, patientId int, permissionType string, message string) {
	doctor := fmt.Sprintf("doctor%d", doctorId)
	patient := fmt.Sprintf("patient%d", patientId)
	uc.client.AddPermission(doctor, patient, permissionType, message)
}

func (uc *appointmentUsecase) UpdateStatus(actor *entities.User, id int, appointment *entities.Appointment) error {
	// check if role of actor is staff
	if actor.Profile.Role != "STAFF" && actor.Profile.Role != "DOCTOR" {
		return errors.New("only staff or doctor can update appointments")
	}

	yesterday := time.Now().AddDate(0, 0, -1)

	if appointment.StartsAt.Before(yesterday) {
		return errors.New("appointment date must be in the future")
	}

	// check that patient can't have 2 concurrent appointments
	appointments, err := uc.appointmentRepo.GetByPatientID(appointment.PatientId)
	if err != nil {
		return errors.New("cant find patient appointments")
	}
	for _, _appointment := range appointments {
		if _appointment.StartsAt.Year() == appointment.StartsAt.Year() &&
			_appointment.StartsAt.Month() == appointment.StartsAt.Month() &&
			_appointment.StartsAt.Day() == appointment.StartsAt.Day() &&
			_appointment.StartsAt.Hour() == appointment.StartsAt.Hour() {
			// cant create
			return errors.New("can't update concurrent appointments")
		}
	}

	if _, err := uc.appointmentRepo.GetById(id); err != nil {
		return err
	}

	appointment.ID = id
	//if all ok updated appointment
	if err := uc.appointmentRepo.Update(appointment); err != nil {
		return err
	}

	if appointment.Status == entities.STARTED {
		message := fmt.Sprintf("%s otorga permiso al doctor %d para ver la información del paciente: %d", actor.Profile.FirstName+" "+actor.Profile.FatherLastName, appointment.DoctorId, appointment.PatientId)
		uc.addPermissionToBlockchain(appointment.DoctorId, appointment.PatientId, "granted", message)
	} else {
		message := fmt.Sprintf("%s revoca permiso al doctor %d para ver la información del paciente: %d", actor.Profile.FirstName+" "+actor.Profile.FatherLastName, appointment.DoctorId, appointment.PatientId)
		uc.addPermissionToBlockchain(appointment.DoctorId, appointment.PatientId, "revoked", message)
	}
	return nil
}

type AppointmentWithDoctor struct {
	entities.Appointment
	Doctor entities.Profile `json:"doctor"`
}

func (uc *appointmentUsecase) GetAll(actor *entities.User) ([]*AppointmentWithDoctor, error) {
	var appointments []*entities.Appointment = []*entities.Appointment{}
	var err error
	// check if role of actor is staff
	if actor.Profile.Role == entities.STAFF {
		appointments, err = uc.appointmentRepo.GetAll()
		if err != nil {
			log.Println("Error getting all appointments (staff)" + err.Error())
			return nil, err
		}
	}

	if actor.Profile.Role == entities.DOCTOR {
		appointments, err = uc.appointmentRepo.GetByDoctorID(actor.ID)
		if err != nil {
			return nil, err
		}
		if len(appointments) == 0 {
			appointments = []*entities.Appointment{}
		}
	}

	var appointmentsWithDoctor []*AppointmentWithDoctor = []*AppointmentWithDoctor{}
	// add doctor and patient data
	for _, appointment := range appointments {
		doctor, err := uc.userRepo.GetUserByID(appointment.DoctorId)
		if err != nil {
			log.Println("Couldnt find Doctor with ID", appointment.DoctorId)
			continue
			// return nil, err
		}
		appointment.Doctor = *doctor

		// this shouldnt even be valid
		if appointment.PatientId != 0 {
			patient, err := uc.profileRepo.GetById(appointment.PatientId)
			if err != nil {
				log.Println("ERROR: error getting patient")
			}
			appointment.Patient = *patient

			doctorProfile, err := uc.profileRepo.GetById(doctor.ProfileId)
			if err != nil {
				log.Println("ERROR: error getting doctor")
			}
			doctor.Profile = *doctorProfile
		}
		appointmentsWithDoctor = append(appointmentsWithDoctor, &AppointmentWithDoctor{
			Appointment: *appointment,
			Doctor:      doctor.Profile,
		})
	}
	return appointmentsWithDoctor, nil
	// return nil, errors.New("only staff or doctors can get appointments")
}

func (uc *appointmentUsecase) GetById(actor *entities.User, id int) (*entities.Appointment, error) {
	// check if role of actor is doctor

	if actor.Profile.Role != "STAFF" && actor.Profile.Role != "DOCTOR" {
		return nil, errors.New("only staff or doctor can update appointments")
	}
	appointment, err := uc.appointmentRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	return appointment, nil

}

// FinishedAppointment implements AppointmentUsecase.
func (uc *appointmentUsecase) EndAppointment(actor *entities.User, id int) error {
	if actor.Profile.Role != "DOCTOR" {
		return errors.New("only staff or doctor can update appointments")
	}

	appointment, err := uc.appointmentRepo.GetById(id)
	if err != nil {
		return err
	}

	if appointment.Status != entities.STARTED {
		return errors.New("appointment must be started")
	}

	appointment.Status = entities.DONE
	if err := uc.appointmentRepo.Update(appointment); err != nil {
		return err
	}

	message := fmt.Sprintf("%s finaliza cita, revoca permiso para ver la información del paciente: %d", actor.Profile.FirstName+" "+actor.Profile.FatherLastName, appointment.PatientId)
	uc.addPermissionToBlockchain(appointment.DoctorId, appointment.PatientId, "revoked", message)
	return nil
}
