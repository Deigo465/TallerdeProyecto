package entities

import "testing"

func TestNewRecord(t *testing.T) {
	//given
	body := "New Body"
	createdAt := "New Created"
	UpdatedAt := "New Update"
	patientId := 3
	doctorId := 2
	//when
	record := NewRecord(1, body, createdAt, UpdatedAt, patientId, doctorId)

	//then
	if record.Body != body {
		t.Fatalf("Expecting body to be %s, got %s", body, record.Body)
	}
	if record.CreatedAt != createdAt {
		t.Fatalf("Expecting created at to be %s, got %s", createdAt, record.CreatedAt)
	}
	if record.UpdatedAt != UpdatedAt {
		t.Fatalf("Expecting update at to be %s, got %s", UpdatedAt, record.UpdatedAt)
	}
	if record.PatientId != patientId {
		t.Fatalf("Expecting patient id to be %d, got %d", patientId, record.PatientId)
	}
	if record.DoctorId != doctorId {
		t.Fatalf("Expecting doctor id to be %d, got %d", doctorId, record.DoctorId)
	}
}

func TestNewFakeRecord(t *testing.T) {

	record := NewFakeRecord()

	if record.Body == "" {
		t.Fatalf("Expecting body to not be empty")
	}
	if record.CreatedAt == "" {
		t.Fatalf("Expecting created at to not be empty")
	}
	if record.UpdatedAt == "" {
		t.Fatalf("Expecting update at to not be empty")
	}
	if record.PatientId == 0 {
		t.Fatalf("Expecting patient id to not be empty")
	}
	if record.DoctorId == 0 {
		t.Fatalf("Expecting doctor id to not be empty")
	}
}
