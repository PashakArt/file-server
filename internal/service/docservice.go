package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"os"
	"path/filepath"

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
		return fmt.Errorf("DocService:Upload:uuid.Parse - %w", err)
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
		if meta.Mime == "" {
			detectedMime := mime.TypeByExtension(filepath.Ext(meta.Name))
			if detectedMime != "" {
				meta.Mime = detectedMime
			} else {
				meta.Mime = "application/octet-stream"
			}
		}

		if err := os.MkdirAll(s.cfg.UploadDir, 0755); err != nil {
			return fmt.Errorf("DocService:Upload:os.MkdirAll - %w", err)
		}

		filePath = filepath.Join(s.cfg.UploadDir, docID.String())
		dst, err := os.Create(filePath)

		if err != nil {
			return fmt.Errorf("Failed to create file - %w", err)
		}

		_, err = io.Copy(dst, file)
		dst.Close()

		if err != nil {
			os.Remove(filePath)
			return fmt.Errorf("DocService:Upload:io.Copy - %w", err)
		}
	}

	document := domain.Document{
		ID:       docID,
		Name:     meta.Name,
		Mime:     &meta.Mime,
		FilePath: &filePath,
		IsPublic: meta.Public,
		HasFile:  meta.File,
		JSONData: *jsonField,
		OwnerID:  ownerID,
	}

	if len(userIds) > 0 {
		err = s.repo.SaveDocumentWithGrants(ctx, &document, userIds)
		if err != nil {
			os.Remove(filePath)
			return fmt.Errorf("DocService:Upload:repo.SaveDocumentWithGrants - %w", err)
		}
		return nil
	}

	err = s.repo.SaveDocument(ctx, &document)
	if err != nil {
		os.Remove(filePath)
		return fmt.Errorf("DocService:Upload:repo.SaveDocument - %w", err)
	}
	return nil
}

func (s *DocService) DeleteById(ctx context.Context, idStr string, ownerIDStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("DocService:DeleteById:uuid.Parse id - %w", err)
	}

	ownerID, err := uuid.Parse(ownerIDStr)
	if err != nil {
		return fmt.Errorf("DocService:DeleteById:uuid.Parse ownerID - %w", err)
	}

	doc, err := s.repo.DeleteById(ctx, id, ownerID)
	if err != nil {
		return fmt.Errorf("DocService:DeleteById:s.repo.DeleteById - %w", err)
	}

	if doc.HasFile {
		err = os.Remove(*doc.FilePath)
		if err != nil {
			log.Printf("DocService:DeleteById:os.Remove - %w", err)
		}
	}

	return nil
}

func (s *DocService) GetList(
	ctx context.Context,
	userIDStr string,
	params GetListParams,
) ([]DocResponse, error) {
	repoParams := repository.GetDocsParams{
		Limit:     params.Limit,
		FilterKey: params.FilterKey,
		FilterVal: params.FilterVal,
	}

	if params.TargetLogin == "" {
		uid, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("DocService.GetList: invalid uuid: %w", err)
		}
		repoParams.OwnerID = &uid
	} else {
		repoParams.TargetLogin = params.TargetLogin
	}

	docs, err := s.repo.GetDocuments(ctx, repoParams)
	if err != nil {
		return nil, fmt.Errorf("DocService:GetList:repo.GetDocuments: %w", err)
	}

	if len(docs) == 0 {
		return []DocResponse{}, nil
	}

	docIDs := make([]uuid.UUID, len(docs))
	for i, d := range docs {
		docIDs[i] = d.ID
	}

	grantsMap, err := s.repo.GetGrantsByDocIDs(ctx, docIDs)
	if err != nil {
		return nil, fmt.Errorf("DocService.GetList: grants: %w", err)
	}

	result := make([]DocResponse, 0, len(docs))
	for _, doc := range docs {
		grants := grantsMap[doc.ID]
		if grants == nil {
			grants = []string{}
		}

		result = append(result, DocResponse{
			ID:      doc.ID.String(),
			Name:    doc.Name,
			Mime:    doc.Mime,
			File:    doc.HasFile,
			Public:  doc.IsPublic,
			Created: doc.Created.Format("2006-01-02 15:04:05"),
			Grant:   grants,
		})
	}

	return result, nil
}

func (s *DocService) GetByID(ctx context.Context, rawDocID, rawUserID string) (*domain.Document, error) {
	docID, err := uuid.Parse(rawDocID)
	if err != nil {
		return nil, fmt.Errorf("DocService:GetByID:uuid.Parse id - %w", err)
	}

	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		return nil, fmt.Errorf("DocService:GetByID:uuid.Parse userID - %w", err)
	}

	doc, err := s.repo.GetByID(ctx, docID, userID)
	if err != nil {
		return nil, fmt.Errorf("DocService.GetByID: %w", err)
	}

	return doc, nil
}
