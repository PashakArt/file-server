package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/PashakArt/file-server/internal/config"
	"github.com/PashakArt/file-server/internal/db/repository"
	"github.com/PashakArt/file-server/internal/domain"
	"github.com/PashakArt/file-server/internal/transport/http/types"
	"github.com/google/uuid"
)

type DocService struct {
	cfg      *config.Config
	repo     *repository.DocRepository
	userRepo *repository.UserRepository
}

func NewDocService(
	cfg *config.Config,
	repo *repository.DocRepository,
	userRepo *repository.UserRepository,
) *DocService {
	return &DocService{
		cfg:      cfg,
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *DocService) Upload(
	ctx context.Context,
	ownerIDStr string,
	file io.Reader,
	meta *types.DocMeta,
	jsonField *json.RawMessage,
) error {
	ownerID, err := uuid.Parse(ownerIDStr)
	if err != nil {
		return fmt.Errorf("DocService:Upload:uuid.Parse", err)
	}

	var userIds []uuid.UUID
	if len(meta.Grant) > 0 {
		userIds, err = s.userRepo.GetIDsByLogins(ctx, meta.Grant)
		if err != nil {
			return fmt.Errorf("DocService:Upload:userRepo.GetIDsByLogins - %w", err)
		}
	}

	docID := uuid.New()

	var filePath string
	if meta.File {
		filePath = path.Join(s.cfg.UploadDir, docID.String())
		dst, err := os.Create(filePath)

		if err != nil {
			return fmt.Errorf("Failed to create file - %w", err)
		}

		_, err = io.Copy(dst, file)
		dst.Close()

		if err != nil {
			os.Remove(filePath)
			return fmt.Errorf("DocService:Upload:io.Copy", err)
		}
	}

	document := domain.Document{
		ID:       docID,
		Name:     meta.Name,
		Mime:     meta.Mime,
		FilePath: filePath,
		IsPublic: meta.Public,
		HasFile:  meta.File,
		JSONData: *jsonField,
		OwnerID:  ownerID,
	}

	if len(userIds) > 0 {
		err = s.repo.SaveDocumentWithGrants(ctx, &document, userIds)
		if err != nil {
			return fmt.Errorf("DocService:Upload:repo.SaveDocumentWithGrants", err)
		}
		return nil
	}

	err = s.repo.SaveDocument(ctx, &document)
	if err != nil {
		return fmt.Errorf("DocService:Upload:repo.SaveDocument", err)
	}
	return nil
}
