package usecase

import (
	"errors"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	"github.com/open-wm/blockehr/pkg/interfaces"
)

type FileUsecase interface {
	Add(actor *entities.User, file *entities.File) error
	GetByRecordId(actor *entities.User, id int) ([]*entities.File, error)
	GetById(actor *entities.User, id int) (*entities.File, error)
}
type fileUsecase struct {
	fileRepo interfaces.FileRepository
}

func NewFileUsecase(repo interfaces.FileRepository) FileUsecase {
	return &fileUsecase{
		fileRepo: repo,
	}
}

func (uc *fileUsecase) GetByRecordId(actor *entities.User, id int) ([]*entities.File, error) {
	// check if role of actor is doctor

	if actor.Profile.Role != entities.DOCTOR {
		return nil, errors.New("only doctor can add file")
	}
	files, err := uc.fileRepo.GetByRecordId(id)
	if err != nil {
		return nil, err
	}

	return files, nil

}
func (uc *fileUsecase) GetById(actor *entities.User, id int) (*entities.File, error) {
	// check if role of actor is doctor

	if actor.Profile.Role != entities.DOCTOR {
		return nil, errors.New("only doctor can add file")
	}
	file, err := uc.fileRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	return file, nil

}
func (uc *fileUsecase) Add(actor *entities.User, file *entities.File) error {
	// check if role of actor is doctor
	if actor.Profile.Role != entities.DOCTOR {
		return errors.New("only doctor can add file")
	}
	// check if the fields are not empty

	if file.Url == "" {
		return errors.New("Expecting url to not be empty")
	}
	if file.Name == "" {
		return errors.New("Expecting name to not be empty")
	}
	if file.FileSize == "" {
		return errors.New("Expecting file size to not be empty")
	}
	if file.MimeType == "" {
		return errors.New("Expecting mime type to not be empty")
	}
	if file.RecordId == 0 {
		return errors.New("Expecting record id to not be empty")
	}

	if err := uc.fileRepo.Add(file); err != nil {
		return err
	}

	return nil
}
