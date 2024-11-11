package usecase

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	"github.com/open-wm/blockehr/pkg/interfaces"
)

type AuthUseCase interface {
	Login(input LoginStruct) (*entities.Session, error)
	Logout(session *entities.Session) error

	GetSession(sessionToken string) *entities.Session
	DeleteSession(sessionToken string) error
}

type SignupStruct struct {
	// User info
	Email    string `json:"email"`
	Password string `json:"password"`
	// Profile info
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	DocumentNumber string `json:"document_number"`
	Birthdate      string `json:"birthdate"`
	Phone          string `json:"phone"`
	DocumentTypeId int    `json:"document_type_id"`
	Avatar         string `json:"avatar"`
	// Gender         string `json:"gender"`
}

type LoginStruct struct {
	// User info
	Email    string `form:"email" json:"email"`
	Password string `form:"email" json:"password"`
}

type authUsecase struct {
	userRepo    interfaces.UserRepository
	sessRepo    interfaces.SessionRepository
	profileRepo interfaces.ProfileRepository
}

func NewAuthUsecase(r interfaces.UserRepository, sessRepo interfaces.SessionRepository, profileRepo interfaces.ProfileRepository) AuthUseCase {
	return &authUsecase{userRepo: r, sessRepo: sessRepo, profileRepo: profileRepo}
}

var (
	errInvalidCredentials = errors.New("invalid credentials")
	ErrNotAuthorized      = errors.New("not authorized")
)

const (
	SESSION_ID         = "sess_id"
	SESSION_KEY_LENGTH = 100
)

// Generate a random string
func GenerateRandomString(length int) string {
	byteLength := (length * 6) / 8
	b := make([]byte, byteLength)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)[:length]
}

// SESSION_KEY_LENGTH
func (uc *authUsecase) Login(input LoginStruct) (*entities.Session, error) {
	// check if a user with the creds exist
	user, err := uc.userRepo.GetUser(input.Email, input.Password)
	if err != nil {
		log.Println("Error getting user")
		return nil, err
	}

	profile, err := uc.profileRepo.GetById(user.ProfileId)
	if err != nil {
		log.Println("Error getting profile")
		return nil, err
	}
	user.Profile = *profile

	if user.Profile.Role != "STAFF" && user.Profile.Role != "DOCTOR" {
		return nil, errors.New("only staff or doctor can log in")
	}
	// create the session

	sess := &entities.Session{
		UserID:    user.ID,
		User:      user,
		Token:     GenerateRandomString(SESSION_KEY_LENGTH),
		CreatedAt: time.Now(),
	}
	// store the session in the db
	sess, err = uc.sessRepo.SaveSession(sess)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

const (
	DNI = iota // 0
	RUC        // 1
	CE         // 2
)

func (uc *authUsecase) Logout(session *entities.Session) error {
	session.User = nil
	return nil
}

func (uc *authUsecase) GetSession(sessionToken string) *entities.Session {
	sess := uc.sessRepo.GetSession(sessionToken)
	if sess != nil {
		user, err := uc.userRepo.GetUserByID(sess.UserID)
		if err != nil {
			log.Printf("no user found %d\n", sess.UserID)
			return nil
		}

		profile, err := uc.profileRepo.GetById(user.ProfileId)
		if err != nil {
			log.Printf("no profile found %d\n", user.ProfileId)
			return nil
		}
		sess.User = user
		sess.User.Profile = *profile
	}
	return sess
}

func (uc *authUsecase) DeleteSession(sessionToken string) error {
	return uc.sessRepo.DeleteSession(sessionToken)
}
