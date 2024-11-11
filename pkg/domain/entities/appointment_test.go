package entities

import (
	"testing"
	"time"
)

func TestNewAppointment(t *testing.T) {
	//given
	specialty := "New Specialty"
	status := PAID
	tomorrow := time.Now().AddDate(0, 0, 1)
	doctor := NewFakeUser()
	patient := NewFakeProfile()
	description := "toi malito"

	//when
	appointment := NewAppointment(1, specialty, status, tomorrow, doctor.ID, doctor, patient.ID, patient, description)

	//then
	if appointment.Specialty != specialty {
		t.Fatalf("Expecting specialty to be %s, got %s", specialty, appointment.Specialty)
	}
	if appointment.Status != status {
		t.Fatalf("Expecting status to be %d, got %d", status, appointment.Status)
	}
	if appointment.StartsAt != tomorrow {
		t.Fatalf("Expecting date appointment to be %s, got %s", &tomorrow, appointment.StartsAt)

	}
	if appointment.DoctorId != doctor.ID {
		t.Fatalf("Expecting doctor id to be %d, got %d", doctor.ID, appointment.DoctorId)
	}
	if appointment.PatientId != patient.ID {
		t.Fatalf("Expecting patient id to be %d, got %d", patient.ID, appointment.PatientId)
	}
	if appointment.Description != description {
		t.Fatalf("Expecting description  to be %s, got %s", description, appointment.Description)
	}
}

func TestNewFakeAppointment(t *testing.T) {
	appointment := NewFakeAppointment()

	if appointment.Specialty == "" {
		t.Fatalf("Expecting specialty to not be empty")
	}
	if appointment.Status == 0 {
		t.Fatalf("Expecting status to not be empty")
	}

	if appointment.DoctorId == 0 {
		t.Fatalf("Expecting doctor id to not be empty")
	}
}
