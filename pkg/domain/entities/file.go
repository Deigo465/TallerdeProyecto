package entities

type File struct {
	ID       int    `json:"id"`
	Url      string `json:"url"`
	Name     string `json:"name"`
	FileSize string `json:"file_size"`
	MimeType string `json:"mime_type"`
	RecordId int    `json:"record_id"`
}

func NewFile(id int, url, name, fileSize, mimeType string, recordId int) File {
	return File{id, url, name, fileSize, mimeType, recordId}
}

func NewFakeFile() File {
	return File{1, "https://medlineplus.gov/images/Xray_share.jpg", "Radiografía", "455 KB", ".jpg", 1}
}
