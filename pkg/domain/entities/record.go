package entities

type Record struct {
	ID        int      `json:"id"`
	Body      string   `json:"body"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	PatientId int      `json:"profile_id"`
	DoctorId  int      `json:"doctor_id"`
	Doctor    *Profile `json:"doctor"`
	Patient   *Profile `json:"patient"`
	Files     []*File  `json:"files"`
	Specialty string   `json:"specialty"`
}

func NewRecord(id int, body, createdAt, UpdatedAt string, patientId, doctorId int) Record {
	return Record{
		ID:        id,
		Body:      body,
		CreatedAt: createdAt,
		UpdatedAt: UpdatedAt,
		PatientId: patientId,
		DoctorId:  doctorId,
		Specialty: "TBD",
	}
}

func NewFakeRecord() Record {
	doctor := NewFakeProfile()
	patient := NewFakeProfile()
	file := NewFakeFile()
	return Record{
		ID:        1,
		Body:      "Juan Pérez acude a la consulta médica quejándose de dolor en el pecho y dificultad para respirar.",
		CreatedAt: "01-06-2003",
		UpdatedAt: "29-04-2004",
		PatientId: 3,
		DoctorId:  2,
		Doctor:    &doctor,
		Patient:   &patient,
		Files:     []*File{&file},
		Specialty: "Cardiología",
	}
}
