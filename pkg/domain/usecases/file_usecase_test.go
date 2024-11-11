package usecase_test

import (
	"testing"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	usecase "github.com/open-wm/blockehr/pkg/domain/usecases"
	moc "github.com/open-wm/blockehr/pkg/mocks"
)

func TestAddFile(t *testing.T) {
	//GIVEN
	fileRepo := moc.NewInMemoryFileRepository()
	uc := usecase.NewFileUsecase(fileRepo)

	//doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)
	//staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(2, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)

	//file
	newFile := entities.NewFile(1, "https://medlineplus.gov/images/Xray_share.jpg", "Radiografía", "455 KB", ".jpg", 2)
	//invalid
	invFileNoUrl := entities.NewFile(1, "", "Radiografía", "455 KB", ".jpg", 2)
	invFileNoName := entities.NewFile(1, "https://medlineplus.gov/images/Xray_share.jpg", "", "455 KB", ".jpg", 2)
	invFileNoSize := entities.NewFile(1, "https://medlineplus.gov/images/Xray_share.jpg", "Radiografía", "", ".jpg", 2)
	invFileNoMimeType := entities.NewFile(1, "https://medlineplus.gov/images/Xray_share.jpg", "Radiografía", "455 KB", "", 2)
	invFileNoRecord := entities.NewFile(1, "https://medlineplus.gov/images/Xray_share.jpg", "Radiografía", "455 KB", ".jpg", 0)

	testCases := []struct {
		name    string
		actor   entities.User
		file    *entities.File
		want    error
		wantErr bool
	}{
		{
			name:    "valid file with doctor",
			actor:   doctorActor,
			file:    &newFile,
			wantErr: false,
			want:    nil,
		},
		{
			name:    "invalid file -  file made by non-doctor",
			actor:   staffActor,
			file:    &newFile,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid file - url is empty",
			actor:   doctorActor,
			file:    &invFileNoUrl,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid file - name is empty",
			actor:   doctorActor,
			file:    &invFileNoName,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid file - size is empty",
			actor:   doctorActor,
			file:    &invFileNoSize,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid file - mime type is empty",
			actor:   doctorActor,
			file:    &invFileNoMimeType,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid file - record id is empty",
			actor:   doctorActor,
			file:    &invFileNoRecord,
			wantErr: true,
			want:    nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN
			err := uc.Add(&tc.actor, tc.file)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Error adding file: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Error adding file: %v", err)
				}

				_, err := fileRepo.GetById(tc.file.ID)
				if err != nil {
					t.Errorf("file not added to repository")
				}

			}
		})
	}

}

func TestGetFilesByRecordId(t *testing.T) {
	//GIVEN
	fileRepo := moc.NewInMemoryFileRepository()
	uc := usecase.NewFileUsecase(fileRepo)

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
		ID:        3,
		Body:      "Italo",
		CreatedAt: "01-06-2024",
		UpdatedAt: "01-06-2024",
		PatientId: 2,
	}

	testCases := []struct {
		name     string
		actor    entities.User
		recordId int
		want     error
		wantErr  bool
	}{
		{
			name:     "valid files with doctor",
			actor:    doctorActor,
			recordId: validRecord.ID,
			wantErr:  false,
			want:     nil,
		},
		{
			name:     "invalid file -  file made by non-doctor",
			actor:    staffActor,
			recordId: validRecord.ID,
			wantErr:  true,
			want:     nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN

			_, err := uc.GetByRecordId(&tc.actor, tc.recordId)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Error getting files: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Error getting files: %v", err)
				}
			}
		})
	}

}

func TestGetFileById(t *testing.T) {
	//GIVEN
	fileRepo := moc.NewInMemoryFileRepository()
	uc := usecase.NewFileUsecase(fileRepo)

	//doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)
	//staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(2, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)

	//file
	validFile := entities.File{ID: 2, Url: "https://medlineplus.gov/images/Xray_share.jpg", Name: "Radiografía", FileSize: "455 KB", MimeType: ".jpg", RecordId: 1}

	testCases := []struct {
		name    string
		actor   entities.User
		file    int
		want    error
		wantErr bool
	}{
		{
			name:    "valid file with doctor",
			actor:   doctorActor,
			file:    validFile.ID,
			wantErr: false,
			want:    nil,
		},
		{
			name:    "invalid file -  file made by non-doctor",
			actor:   staffActor,
			file:    validFile.ID,
			wantErr: true,
			want:    nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN

			_, err := uc.GetById(&tc.actor, tc.file)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Error getting file: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Error getting file: %v", err)
				}
			}
		})
	}

}
