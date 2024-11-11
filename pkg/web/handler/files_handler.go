package handler

import (
	"strconv"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/open-wm/blockehr/pkg/domain/entities"
	usecase "github.com/open-wm/blockehr/pkg/domain/usecases"
)

type FileHandler interface {
	Add(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	GetProfileById(actor *entities.User, rw http.ResponseWriter, r *http.Request)
	GetByRecordId(actor *entities.User, rw http.ResponseWriter, r *http.Request)
}

type fileHandler struct {
	fileUsecase   usecase.FileUsecase
	recordUsecase usecase.RecordUsecase
}

func NewFiletHandler(ucFile usecase.FileUsecase, ucRecord usecase.RecordUsecase) FileHandler {
	return &fileHandler{
		fileUsecase:   ucFile,
		recordUsecase: ucRecord,
	}
}

func (h *fileHandler) Add(actor *entities.User, rw http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10 MB
	multipart := r.MultipartForm
	if multipart == nil {
		JSON(rw, http.StatusBadRequest, "No file found", nil)
		return
	}
	files := multipart.File["file[]"]
	if len(files) == 0 {
		JSON(rw, http.StatusBadRequest, "No file found", nil)
		return
	}
	file := files[0]
	// Save file
	if file != nil {
		fileModel, err := StoreFile(file)
		fileModel.RecordId = 1
		if err != nil {
			log.Println(err)
			JSON(rw, http.StatusInternalServerError, err.Error(), nil)
		}
		err = h.fileUsecase.Add(actor, fileModel)
		if err != nil {
			log.Println(err)
			JSON(rw, http.StatusInternalServerError, err.Error(), nil)
		}
		JSON(rw, http.StatusCreated, "Media uploaded successfully", fileModel)
	}
	JSON(rw, http.StatusInternalServerError, "something weird happened", nil)
}

func (h *fileHandler) GetProfileById(actor *entities.User, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	fileID, err := strconv.Atoi(vars["file_id"]) //TODO

	if err != nil {
		http.Error(rw, "Invalid file ID", http.StatusBadRequest)
		return
	}
	file, err := h.fileUsecase.GetById(actor, fileID)
	if err != nil {
		http.Error(rw, "Failed to returned file", http.StatusInternalServerError)
		return
	}
	JSON(rw, http.StatusOK, "file returned", file)
}

func (h *fileHandler) GetByRecordId(actor *entities.User, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	recordId, err := strconv.Atoi(vars["record_id"]) //TODO

	if err != nil {
		http.Error(rw, "Invalid record ID", http.StatusBadRequest)
		return
	}
	files, err := h.fileUsecase.GetByRecordId(actor, recordId)
	if err != nil {
		http.Error(rw, "Failed to returned file", http.StatusInternalServerError)
		return
	}
	JSON(rw, http.StatusOK, "file returned", files)
}
