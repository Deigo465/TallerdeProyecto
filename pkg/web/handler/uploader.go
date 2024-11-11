package handler

import (
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	usecase "github.com/open-wm/blockehr/pkg/domain/usecases"
)

// type Context interface {
// 	// JSON(code int, i interface{}) error
// 	// Bind(i interface{}) error
// 	BodyParser(out interface{}) error
// 	ParamsInt(key string, defaultValue ...int) (int, error)
// 	FormFile(key string) (*multipart.FileHeader, error)
// 	// JSONError(code int, error error) error
// 	// Render(code int, data Map, templateName string, shouldRedirect bool) error
// 	Params(key string, defaultValue ...string) string
// 	Redirect(path string, statusCode ...int) error
// 	Render(templateName string, data interface{}, layout ...string) error
// 	SaveFile(fileheader *multipart.FileHeader, path string) error
// 	Status(status int) Context
// 	Next() error
// 	SendFile(file string, compress ...bool) error
// 	Query(key string, defaultValue ...string) string
// 	Method(override ...string) string
// 	GetReqHeaders() map[string]string
// }

// Map is a shortcut for map[string]interface{}, useful for JSON returns
type Map map[string]interface{}

func Upload(rw http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10 MB
	files := r.MultipartForm.File["file"]
	file := files[0]
	// Save file
	if file != nil {
		_, err := StoreFile(file)
		if err != nil {
			log.Println(err)
			JSON(rw, http.StatusInternalServerError, err.Error(), nil)
		}
		if err != nil {
			log.Println(err)
			JSON(rw, http.StatusInternalServerError, err.Error(), nil)
		}
		JSON(rw, http.StatusCreated, "Media uploaded successfully", file)
		return
	}
	JSON(rw, http.StatusInternalServerError, "something weird happened", nil)
}

func StoreFile(input *multipart.FileHeader) (*entities.File, error) {
	contType := input.Header.Get("Content-Type")
	mediaType, _, _ := mime.ParseMediaType(contType)

	contDisposition := input.Header.Get("Content-Disposition")
	_, params, _ := mime.ParseMediaType(contDisposition)

	// Upload the file to specific dst.
	// NOTE: this could actually replace files if they collide (40**8 chance tho)
	extension := filepath.Ext(params["filename"])
	newFileName := "/user_uploads/" + usecase.GenerateRandomString(8) + extension

	basePathFormat := "storage%s"
	// Check if running from tests
	if strings.Contains(os.Args[0], ".test") {
		basePathFormat = "../../../storage%s"
	}
	// out, err := os.Create(fmt.Sprintf(basePathFormat, newFileName))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer out.Close()

	// save the original file with an ".og" at the end
	err := SaveMultipartFile(input, fmt.Sprintf(basePathFormat, newFileName))
	if err != nil {
		return nil, err
	}

	fileModel := entities.File{
		Url:      newFileName,
		Name:     params["filename"],
		FileSize: fmt.Sprintf("%d", input.Size),
		MimeType: mediaType,
	}

	return &fileModel, nil
}

func SaveMultipartFile(fh *multipart.FileHeader, path string) (err error) {
	var (
		f  multipart.File
		ff *os.File
	)
	f, err = fh.Open()
	if err != nil {
		return
	}

	var ok bool
	if ff, ok = f.(*os.File); ok {
		// Windows can't rename files that are opened.
		if err = f.Close(); err != nil {
			return
		}

		// If renaming fails we try the normal copying method.
		// Renaming could fail if the files are on different devices.
		if os.Rename(ff.Name(), path) == nil {
			return nil
		}

		// Reopen f for the code below.
		if f, err = fh.Open(); err != nil {
			return
		}
	}

	defer func() {
		e := f.Close()
		if err == nil {
			err = e
		}
	}()

	if ff, err = os.Create(path); err != nil {
		return
	}
	defer func() {
		e := ff.Close()
		if err == nil {
			err = e
		}
	}()
	_, err = copyZeroAlloc(ff, f)
	return
}

func copyZeroAlloc(w io.Writer, r io.Reader) (int64, error) {
	if wt, ok := r.(io.WriterTo); ok {
		return wt.WriteTo(w)
	}
	if rt, ok := w.(io.ReaderFrom); ok {
		return rt.ReadFrom(r)
	}
	vbuf := copyBufPool.Get()
	buf := vbuf.([]byte)
	n, err := io.CopyBuffer(w, r, buf)
	copyBufPool.Put(vbuf)
	return n, err
}

var copyBufPool = sync.Pool{
	New: func() any {
		return make([]byte, 4096)
	},
}

// "github.com/nfnt/resize"
// func SaveAndResize(input *multipart.FileHeader, out *os.File, extension string) {
// 	// this  s technically the "proper way, but I think due to the content disposition"
// 	// it fails
// 	// extension := ".jpg"
// 	// bytes, err := ioutil.ReadAll(f)
// 	// if err != nil {
// 	// 	log.Println(err)
// 	// }
// 	// contentType := http.DetectContentType(bytes)

// 	// switch contentType {
// 	// case "image/png":
// 	// 	extension = ".png"
// 	// case "image/jpeg":
// 	// 	extension = ".jpg"
// 	// }
// 	f, err := input.Open()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var img image.Image
// 	if extension == ".jpg" {
// 		// decode jpeg into image.Image
// 		img, err = jpeg.Decode(f)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// 	if extension == ".png" {
// 		// decode jpeg into image.Image
// 		img, err = png.Decode(f)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}
// 	f.Close()

// 	// resize to width 1000 using Lanczos resampling
// 	// and preserve aspect ratio
// 	// m := resize.Resize(1000, 0, img, resize.Lanczos3)

// 	// preserve original format, we can change this later too
// 	// if extension == ".jpg" {
// 	// 	// write new image to file
// 	// 	jpeg.Encode(out, m, nil)
// 	// }
// 	// if extension == ".png" {
// 	// 	// write new image to file
// 	// 	png.Encode(out, m)
// 	// }
// }
