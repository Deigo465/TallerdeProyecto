package usecase_test

import (
	"testing"
	"time"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	usecase "github.com/open-wm/blockehr/pkg/domain/usecases"
	moc "github.com/open-wm/blockehr/pkg/mocks"
)

func TestPasswordHash(t *testing.T) {
}
func TestLogin(t *testing.T) {
	// Given
	userRepo := moc.NewInMemoryUserRepository()
	sessionRepo := moc.NewInMemorySessionRepository()
	profileRepo := moc.NewInMemoryProfileRepository()
	uc := usecase.NewAuthUsecase(userRepo, sessionRepo, profileRepo)

	// valid credential
	credsDoctor := usecase.LoginStruct{
		Email:    "doctor@blockehr.pe",
		Password: "perroLoco",
	}
	credStaff := usecase.LoginStruct{
		Email:    "staff@blockehr.pe",
		Password: "perroLoco",
	}

	// invalid credential
	credPatient := usecase.LoginStruct{
		Email:    "patient@blockehr.pe",
		Password: "perroLoco",
	}

	testCases := []struct {
		name    string
		creds   usecase.LoginStruct
		wantErr bool
	}{
		{
			name:    "valid login with doctor",
			creds:   credsDoctor,
			wantErr: false,
		},
		{
			name:    "valid login with staff",
			creds:   credStaff,
			wantErr: false,
		},
		{
			name:    "invalid login with patient",
			creds:   credPatient,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			if r, ok := userRepo.(*moc.InMemoryUserRepository); ok {
				r.Reset()
			}
			// WHEN
			session, err := uc.Login(tc.creds)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				userFound, err := userRepo.GetUser(tc.creds.Email, tc.creds.Password)
				if err != nil {
					t.Errorf("Error getting user: %v", err)
				}

				if userFound == nil {
					t.Errorf("User not found in repository")
				}
				if session == nil {
					t.Errorf("Session was not created")
				}
				// log.Printf("User email: %s, Role: %s", userFound.Email, userFound.Profile.Role)
				// log.Printf("Session Token: %s", session.Token)

			}
		})
	}
}

func TestLogout(t *testing.T) {
	// Given
	userRepo := moc.NewInMemoryUserRepository()
	sessionRepo := moc.NewInMemorySessionRepository()
	profileRepo := moc.NewInMemoryProfileRepository()
	uc := usecase.NewAuthUsecase(userRepo, sessionRepo, profileRepo)

	// valid cred
	credsDoctor := usecase.LoginStruct{
		Email:    "doctor@blockehr.pe",
		Password: "perroLoco",
	}
	sess, err := uc.Login(credsDoctor)
	if err != nil {
		t.Fatalf("Error logging in: %s", err)
	}

	// WHEN
	err = uc.Logout(sess)
	if err != nil {
		t.Fatalf("Error logging out: %s", err)
	}
	//Then
	if sess.User != nil {
		t.Fatal("User was expected to be null after logout but it wasn't")
	}

}
func TestGetSession(t *testing.T) {

	// Given
	profile := entities.NewFakeProfile()
	profileRepo := moc.NewInMemoryProfileRepository()
	profileRepo.Add(&profile)

	userRepo := moc.NewInMemoryUserRepository()

	user := entities.NewFakeUser()
	user.ID = 1
	user.ProfileId = profile.ID

	userRepo.SaveUser(&user)

	sessionRepo := moc.NewInMemorySessionRepository()
	sessionRepo.SaveSession(&entities.Session{
		UserID:    1,
		Token:     "abc",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	uc := usecase.NewAuthUsecase(userRepo, sessionRepo, profileRepo)
	sessionToken := "abc"
	//when
	session := uc.GetSession(sessionToken)

	if session == nil {
		t.Error("The session should not be null")
		return
	}

	if session.User == nil {
		t.Error("The session should have a user")
		return
	}
	if session.User.ID != 1 {
		t.Fatalf("Expecting user to be 1 got, %d", session.User.ID)
	}
}
func TestDeleteSession(t *testing.T) {
	// Given
	userRepo := moc.NewInMemoryUserRepository()
	sessionRepo := moc.NewInMemorySessionRepository()
	profileRepo := moc.NewInMemoryProfileRepository()
	uc := usecase.NewAuthUsecase(userRepo, sessionRepo, profileRepo)

	// validCred
	credsDoctor := usecase.LoginStruct{
		Email:    "doctor@blockehr.pe",
		Password: "perroLoco",
	}

	sess, err := uc.Login(credsDoctor)
	if err != nil {
		t.Fatalf("Error logging in: %s", err)
	}

	// When
	err = uc.DeleteSession(sess.Token)

	// Then
	if err != nil {
		t.Fatalf("Error deleting session: %v", err)
	}

	// Verify that the session is deleted
	deletedSession := sessionRepo.GetSession(sess.Token)
	if deletedSession != nil {
		t.Errorf("Session was not deleted")
	}
}
func TestGenerateRandomString(t *testing.T) {
	// Given
	length := 20

	// When
	str := usecase.GenerateRandomString(length)

	// Then
	if str == "" {
		t.Errorf("Expected random string to be generated, but got empty string")
	}

	if len(str) != length {
		t.Errorf("Expected random string to have length 20, but got %v", len(str))
	}

	str2 := usecase.GenerateRandomString(length)

	if str == str2 {
		t.Errorf("Expected str1 to str2 to be different but they were the same %v", str)
	}
}
