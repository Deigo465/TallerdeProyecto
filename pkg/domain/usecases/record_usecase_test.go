package usecase_test

import (
	"log"
	"testing"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	usecase "github.com/open-wm/blockehr/pkg/domain/usecases"
	"github.com/open-wm/blockehr/pkg/interfaces"
	moc "github.com/open-wm/blockehr/pkg/mocks"
)

func getRecordUC() (interfaces.RecordRepository, usecase.RecordUsecase) {
	recordRepo := moc.NewInMemoryRecordRepository()
	profileRepo := moc.NewInMemoryProfileRepository()
	fileRepo := moc.NewInMemoryFileRepository()
	blockchain := moc.NewInMemoryBlockchain()
	uc := usecase.NewRecordUsecase(recordRepo, profileRepo, fileRepo, blockchain)
	return recordRepo, uc
}

func TestAddRecord(t *testing.T) {
	//GIVEN
	recordRepo, uc := getRecordUC()

	//doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)
	//staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(2, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)
	//Record
	newRecord := entities.NewRecord(1, "Juan Pérez acude a la consulta médica quejándose de dolor en el pecho y dificultad para respirar.",
		"01-06-2003", "29-04-2004", 1, 2)
	//invalid
	invRecordNoBody := entities.NewRecord(1, "", "01-06-2003", "29-04-2004", 1, 2)
	invRecordNoCreatedAt := entities.NewRecord(1, "Juan Pérez acude a la consulta médica quejándose de dolor en el pecho y dificultad para respirar.",
		"", "29-04-2004", 1, 2)
	invRecordNoUpdatedAt := entities.NewRecord(1, "Juan Pérez acude a la consulta médica quejándose de dolor en el pecho y dificultad para respirar.",
		"01-06-2003", "", 1, 2)
	invRecordNoPatientId := entities.NewRecord(1, "Juan Pérez acude a la consulta médica quejándose de dolor en el pecho y dificultad para respirar.",
		"01-06-2003", "29-04-2004", 0, 2)
	invRecordNoDoctorId := entities.NewRecord(1, "Juan Pérez acude a la consulta médica quejándose de dolor en el pecho y dificultad para respirar.",
		"01-06-2003", "29-04-2004", 1, 0)

	testCases := []struct {
		name    string
		actor   entities.User
		record  *entities.Record
		want    error
		wantErr bool
	}{
		{
			name:    "valid Record with doctor",
			actor:   doctorActor,
			record:  &newRecord,
			wantErr: false,
			want:    nil,
		},
		{
			name: "invalid Record - record made by  non-doctor",

			actor:   staffActor,
			record:  &newRecord,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid Record - body is empty",
			actor:   doctorActor,
			record:  &invRecordNoBody,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid Record - created at  is empty",
			actor:   doctorActor,
			record:  &invRecordNoCreatedAt,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid Record - updated at  is empty",
			actor:   doctorActor,
			record:  &invRecordNoUpdatedAt,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid Record - Patient id  is empty",
			actor:   doctorActor,
			record:  &invRecordNoPatientId,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid Record - Doctor id  is empty",
			actor:   doctorActor,
			record:  &invRecordNoDoctorId,
			wantErr: true,
			want:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN

			err := uc.Add(&tc.actor, tc.record)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Error adding record: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Error adding record: %v", err)
				}

				records, _ := recordRepo.GetAll()
				found := false

				for _, record := range records {
					if record.ID == tc.record.ID {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Record not added to repository")
				}

			}
		})
	}

}

func TestGetAllRecords(t *testing.T) {
	//GIVEN
	_, uc := getRecordUC()

	//doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)
	//staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(2, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)
	//record
	validRecord := entities.Record{
		ID:        7,
		Body:      "Rodrigo",
		CreatedAt: "01-06-2024",
		UpdatedAt: "01-06-2024",
		PatientId: 3,
	}

	testCases := []struct {
		name    string
		actor   entities.User
		record  *entities.Record
		want    error
		wantErr bool
	}{
		{
			name:   "valid Record with doctor",
			actor:  doctorActor,
			record: &validRecord,

			wantErr: false,
			want:    nil,
		},
		{
			name:    "invalid Record - record made by  non-doctor",
			actor:   staffActor,
			record:  &validRecord,
			wantErr: true,
			want:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN

			xd, err := uc.GetAllForPatient(&tc.actor, tc.record.PatientId)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Error getting records: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Error getting records: %v", err)
				}

			}
			log.Print(len(xd))
		})
	}

}

func TestGetRecordId(t *testing.T) {

	//GIVEN
	_, uc := getRecordUC()

	//doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)
	//staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(2, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)
	//Record
	validRecord := entities.Record{
		ID:        2,
		Body:      "Rodrigo",
		CreatedAt: "01-06-2024",
		UpdatedAt: "01-06-2024",
	}

	testCases := []struct {
		name    string
		actor   entities.User
		record  *entities.Record
		want    error
		wantErr bool
	}{
		{
			name:    "valid Record with doctor",
			actor:   doctorActor,
			record:  &validRecord,
			wantErr: false,
			want:    nil,
		},
		{
			name:    "invalid Record - record made by  non-staff",
			actor:   staffActor,
			record:  &validRecord,
			wantErr: true,
			want:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN

			_, err := uc.GetById(&tc.actor, tc.record.ID)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Error getting record: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Error getting record: %v", err)
				}

			}

		})

	}

}

func TestUpdatedByPatientId(t *testing.T) {

	//GIVEN
	_, uc := getRecordUC()

	//doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)
	//staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(2, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)
	//Record
	validRecord := entities.Record{
		ID:        2,
		Body:      "newBody",
		CreatedAt: "01-06-2024",
		UpdatedAt: "01-06-2024",
		PatientId: 3,
	}

	testCases := []struct {
		name    string
		actor   entities.User
		record  *entities.Record
		want    error
		wantErr bool
	}{
		{
			name:    "valid Record with doctor",
			actor:   doctorActor,
			record:  &validRecord,
			wantErr: false,
			want:    nil,
		},
		{
			name:    "invalid Record - record made by  non-staff",
			actor:   staffActor,
			record:  &validRecord,
			wantErr: true,
			want:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN

			err := uc.UpdateByPatientId(&tc.actor, tc.record.PatientId, tc.record.Body)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Error updating record: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Error updating record: %v", err)
				}

			}

		})

	}

}
