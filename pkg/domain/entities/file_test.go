package entities

import "testing"

func TestNewFile(t *testing.T) {
	//given
	url := "New Url"
	name := "New Name"
	fileSize := "New File Size"
	mimeType := "New Mime Type"
	recordId := 1

	//when
	newFile := NewFile(1, url, name, fileSize, mimeType, recordId)

	//then
	if newFile.Url != url {
		t.Fatalf("Expecting url to be %s, got %s", url, newFile.Url)
	}
	if newFile.Name != name {
		t.Fatalf("Expecting name to be %s, got %s", name, newFile.Name)
	}
	if newFile.FileSize != fileSize {
		t.Fatalf("Expecting file size to be %s, got %s", fileSize, newFile.FileSize)
	}
	if newFile.MimeType != mimeType {
		t.Fatalf("Expecting mime type size to be %s, got %s", mimeType, newFile.MimeType)
	}
	if newFile.RecordId != recordId {
		t.Fatalf("Expecting record id to be %d, got %d", recordId, newFile.RecordId)
	}
}

func TestNewFakeFile(t *testing.T) {

	newFile := NewFakeFile()

	if newFile.Url == "" {
		t.Fatalf("Expecting url to not be empty")
	}
	if newFile.Name == "" {
		t.Fatalf("Expecting name to not be empty")
	}
	if newFile.FileSize == "" {
		t.Fatalf("Expecting file size to not be empty")
	}
	if newFile.MimeType == "" {
		t.Fatalf("Expecting mime type to not be empty")
	}
	if newFile.RecordId == 0 {
		t.Fatalf("Expecting record id to not be empty")
	}

}
