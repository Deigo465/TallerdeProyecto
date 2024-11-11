package entities

import "time"

type AppointmentStatus int

const (
	PENDING AppointmentStatus = iota
	PAID
	STARTED
	DONE
	CANCELED
)

type Appointment struct {
	ID          int               `json:"id"`
	Specialty   string            `json:"specialty"`
	Status      AppointmentStatus `json:"status"`
	StartsAt    time.Time         `json:"starts_at"`
	DoctorId    int               `json:"doctor_id"`
	Doctor      User              `json:"doctor"`
	PatientId   int               `json:"patient_id"`
	Patient     Profile           `json:"patient"`
	Description string            `json:"description"`
}

func NewAppointment(id int, specialty string, status AppointmentStatus, startsAt time.Time, doctorId int, doctor User, patientId int, patient Profile, description string) Appointment {
	return Appointment{id, specialty, status, startsAt, doctorId, doctor, patientId, patient, description}
}
func NewFakeAppointment() Appointment {
	tomorrow := time.Now().AddDate(0, 0, 1)
	return Appointment{1, "Psicólogo", PAID, tomorrow, 1, NewFakeUser(), 1, NewFakeProfile(), "toi malito"}

}
