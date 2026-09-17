package repository

import (
	"context"
	_ "embed"
	"errors"

	"github.com/PashakArt/file-server/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	//go:embed query/create_document.sql
	createDocumentQuery string

	//go:embed query/create_grants.sql
	createGrantsQuery string

	//go:embed query/delete_document_by_id.sql
	deleteByIdQuery string
)

type DocRepository struct {
	dbPool *pgxpool.Pool
}

func NewDocRepository(
	dbPool *pgxpool.Pool,
) *DocRepository {
	return &DocRepository{
		dbPool: dbPool,
	}
}

func (r *DocRepository) SaveDocumentWithGrants(
	ctx context.Context,
	doc *domain.Document,
	userIds []uuid.UUID,
) error {
	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		createDocumentQuery,
		doc.ID,
		doc.OwnerID,
		doc.Name,
		doc.Mime,
		doc.HasFile,
		doc.IsPublic,
		doc.JSONData,
		doc.FilePath,
	)
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}
	for _, userID := range userIds {
		batch.Queue(createGrantsQuery, doc.ID, userID)
	}
	res := tx.SendBatch(ctx, batch)
	err = res.Close()
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *DocRepository) SaveDocument(
	ctx context.Context,
	doc *domain.Document,
) error {
	_, err := r.dbPool.Exec(
		ctx,
		createDocumentQuery,
		doc.ID,
		doc.OwnerID,
		doc.Name,
		doc.Mime,
		doc.HasFile,
		doc.IsPublic,
		doc.JSONData,
		doc.FilePath,
	)

	return err
}

func (r *DocRepository) DeleteById(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) (*domain.Document, error) {
	var document domain.Document

	err := r.dbPool.QueryRow(
		ctx,
		deleteByIdQuery,
		id,
		ownerID,
	).Scan(&document.HasFile, &document.FilePath)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDocNotFound
		}

		return nil, err
	}

	return &document, err
}
