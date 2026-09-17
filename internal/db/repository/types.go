package repository

import (
	"github.com/PashakArt/file-server/internal/domain"
	"github.com/google/uuid"
)

type GetDocsParams struct {
	TargetLogin string
	OwnerID     *uuid.UUID
	FilterKey   string
	FilterVal   string
	Limit       int
}

type DocWithGrants struct {
	Document domain.Document
	Grants   []string
}
