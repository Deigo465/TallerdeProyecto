package usecase

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	"github.com/open-wm/blockehr/pkg/interfaces"
)

type UserUsecase interface {
	SaveUser(actor *entities.User, user *entities.User) error
	ResetDoctorPassword(actor *entities.User, id int) (string, error)
}

type userUsecase struct {
	userRepo interfaces.UserRepository
}

func NewUserUsecase(repo interfaces.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: repo,
	}
}
func (uc *userUsecase) SaveUser(actor *entities.User, user *entities.User) error {
	// check if role of actor is staff

	if actor.Profile.Role != "STAFF" {
		return errors.New("only staff can add users")
	}
	// check if the fields are not empty
	if user.Email == "" {
		return errors.New("Expecting email to not be empty")
	}
	if user.Password == "" {
		return errors.New("Expecting password to not be empty")
	}
	if user.HealthCenterId == 0 {
		return errors.New("Expecting healthcenter id to not be empty")
	}
	if user.ProfileId == 0 {
		return errors.New("Expecting profile id to not be empty")
	}
	if err := uc.userRepo.SaveUser(user); err != nil {
		return err
	}

	return nil
}

func (uc *userUsecase) ResetDoctorPassword(actor *entities.User, id int) (string, error) {
	if actor.Profile.Role != "STAFF" {
		return "", errors.New("only staff can reset users' passwords")
	}

	newPassword, err := generateRandomPassword(8)
	if err != nil {
		return "", err
	}

	if err := uc.userRepo.UpdatePassword(id, newPassword); err != nil {
		return "", err
	}

	return newPassword, nil
}

func generateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var password string
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password += string(charset[randomIndex.Int64()])
	}
	return password, nil
}
