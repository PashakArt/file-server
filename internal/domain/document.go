package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID       uuid.UUID       `json:"id" db:"id"`
	OwnerID  uuid.UUID       `json:"owner_id" db:"owner_id"`
	Name     string          `json:"name" db:"name"`
	Mime     *string         `json:"mime" db:"mime"`
	HasFile  bool            `json:"has_file" db:"has_file"`
	IsPublic bool            `json:"public" db:"is_public"`
	FilePath *string         `json:"file_path,omitempty" db:"file_path"`
	JSONData json.RawMessage `json:"json,omitempty" db:"json_data"`
	Created  time.Time       `json:"created" db:"created_at"`
}
