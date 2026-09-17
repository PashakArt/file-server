package repository

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

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

	//go:embed query/get_grants_by_doc_id.sql
	getGrantsByDocIDsQuery string
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

func (r *DocRepository) GetGrantsByDocIDs(ctx context.Context, docIDs []uuid.UUID) (map[uuid.UUID][]string, error) {
	grantsMap := make(map[uuid.UUID][]string)
	rows, err := r.dbPool.Query(ctx, getGrantsByDocIDsQuery, docIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var docID uuid.UUID
		var login string
		if err := rows.Scan(&docID, &login); err != nil {
			return nil, err
		}
		grantsMap[docID] = append(grantsMap[docID], login)
	}

	return grantsMap, nil
}

func (r *DocRepository) GetDocuments(ctx context.Context, params GetDocsParams) ([]domain.Document, error) {
	query := `
		SELECT 
			d.id, d.owner_id, d.name, d.mime, d.has_file, 
			d.is_public, d.file_path, d.json_data, d.created_at
		FROM documents d
		JOIN users u ON d.owner_id = u.id
		WHERE 1=1`

	args := []any{}
	argID := 1

	if params.OwnerID != nil {
		query += fmt.Sprintf(" AND d.owner_id = $%d", argID)
		args = append(args, *params.OwnerID)
		argID++
	} else if params.TargetLogin != "" {
		query += fmt.Sprintf(" AND u.login = $%d", argID)
		args = append(args, params.TargetLogin)
		argID++
	}

	if params.FilterKey != "" && params.FilterVal != "" {
		switch params.FilterKey {
		case "name":
			query += fmt.Sprintf(" AND d.name = $%d", argID)
			args = append(args, params.FilterVal)
			argID++
		case "public":
			query += fmt.Sprintf(" AND d.is_public = $%d", argID)
			args = append(args, params.FilterVal == "true")
			argID++
		case "mime":
			query += fmt.Sprintf(" AND d.mime = $%d", argID)
			args = append(args, params.FilterVal)
			argID++
		}
	}

	query += fmt.Sprintf(" ORDER BY d.name ASC, d.created_at DESC LIMIT $%d", argID)
	args = append(args, params.Limit)

	fmt.Println(args...)
	rows, err := r.dbPool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []domain.Document
	for rows.Next() {
		var d domain.Document
		err := rows.Scan(
			&d.ID,
			&d.OwnerID,
			&d.Name,
			&d.Mime,
			&d.HasFile,
			&d.IsPublic,
			&d.FilePath,
			&d.JSONData,
			&d.Created,
		)
		if err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}

	return docs, rows.Err()
}
