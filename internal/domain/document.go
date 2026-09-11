package domain

import (
	"encoding/json"
	"time"
)

type Document struct {
	ID         string          `json:"id" db:"id"`
	Name       string          `json:"name" db:"name"`
	Mime       string          `json:"mime" db:"mime"`
	IsFile     bool            `json:"file" db:"is_file"`
	IsPublic   bool            `json:"public" db:"is_public"`
	FilePath   string          `json:"file_path,omitempty" db:"file_path"`
	JSONData   json.RawMessage `json:"json,omitempty" db:"json_data"`
	OwnerLogin string          `json:"owner_login" db:"owner_login"`
	Grant      []string        `json:"grant" db:"grant"`
	Created    time.Time       `json:"created" db:"created_at"`
}
