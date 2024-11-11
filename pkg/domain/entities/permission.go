package entities

type Permission struct {
	CreatedAt string `json:"created_at"`
	Type      string `json:"type"` // type can be "granted" or "revoked"
}
