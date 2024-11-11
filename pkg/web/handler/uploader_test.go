package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"strings"
	"testing"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	usecase "github.com/open-wm/blockehr/pkg/domain/usecases"
	"github.com/stretchr/testify/assert"
)

type HandlerTestSuite struct {
	cookies []*http.Cookie
	user    entities.User
	db      *sql.DB
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

func createFormData(values map[string]io.Reader, contentType string) (bytes.Buffer, string, error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		var err error
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = createFormFile(w, key, x.Name(), contentType); err != nil {
				return b, "", err
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return b, "", err
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return b, "", err
		}
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()
	return b, w.FormDataContentType(), nil
}

func createFormFile(w *multipart.Writer, fieldname, filename, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s";  filename="%s"`, fieldname, filename))
	h.Set("Content-Type", contentType)
	return w.CreatePart(h)
}

func TestCreateMediaFile(t *testing.T) {
	// Given

	// Check if running from tests
	file := "./storage/test/file.txt"
	if strings.Contains(os.Args[0], ".test") {
		file = "../../../storage/test/file.txt"
	}
	values := map[string]io.Reader{
		"file":  mustOpen(file), // lets assume its this file
		"other": strings.NewReader("hello world!"),
	}
	b, contentType, _ := createFormData(values, "text/plain")

	// When

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", "/uploads", &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", contentType)
	cookie := http.Cookie{
		Name:  usecase.SESSION_ID,
		Value: "STAFF",
	}
	req.AddCookie(&cookie)

	// Submit the request
	SetBasePath("./views/")
	w := httptest.NewRecorder()

	// WHEN
	Upload(w, req)

	// Then
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code %d, but got %d", http.StatusCreated, resp.StatusCode)
	}
	jsonObject := struct {
		Status  int           `json:"status"`
		Data    entities.File `json:"data"`
		Message string        `json:"message"`
	}{}
	err = parseResponse(resp, &jsonObject)

	assert.Nil(t, err)
	assert.Equal(t, "Media uploaded successfully", jsonObject.Message)
	assert.NotNil(t, jsonObject.Data.Url)
	assert.Equal(t, http.StatusCreated, jsonObject.Status)
}

// Helper function that parses an HTTP Response into a predefined struct
func parseResponse(res *http.Response, v interface{}) error {
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	// Uncomment to see the body for debugging
	// fmt.Printf("%v --- \n", string(bytes))

	if res.StatusCode == 404 {
		return errors.New(res.Status)
	}

	err = json.Unmarshal(bytes, &v)
	if err != nil {
		return err
	}
	return nil
}
