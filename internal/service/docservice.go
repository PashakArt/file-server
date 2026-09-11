package service

import (
	"context"
	"encoding/json"
	"io"

	"github.com/PashakArt/file-server/internal/transport/http/types"
)

type DocService struct {
}

func NewDocService() *DocService {
	return &DocService{}
}

func (s *DocService) Upload(
	ctx context.Context,
	file io.Reader,
	meta *types.DocMeta,
	jsonField *json.RawMessage,
) error {
	return nil
}
