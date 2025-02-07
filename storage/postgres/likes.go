package postgres

import (
	"context"
	"database/sql"
	"fmt"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"app/api/models"
	"app/pkg/helper"
)

type likeRepo struct {
	db *pgxpool.Pool
}

func NewLikeRepo(db *pgxpool.Pool) *likeRepo {
	return &likeRepo{
		db: db,
	}
}

func (r *likeRepo) Create(ctx context.Context, req *models.CreateLike) (string, error) {
	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO likes(id, count, publication_id, contributor_id, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.Count,
		req.PublicationID,
		req.ContributorID,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *likeRepo) GetByID(ctx context.Context, req *models.LikePrimaryKey) (*models.Like, error) {

	var (
		query string

		id            sql.NullString
		count         sql.NullFloat64
		publicationID sql.NullString
		contributorID sql.NullString
		createdAt     sql.NullString
		updatedAt     sql.NullString
	)

	query = `
		SELECT
			id,
			count,
			publication_id,
			contributor_id,
			created_at,
			updated_at
		FROM likes
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&count,
		&publicationID,
		&contributorID,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Like{
		Id:            id.String,
		Count:         count.Float64,
		PublicationID: publicationID.String,
		ContributorID: contributorID.String,
		CreatedAt:     createdAt.String,
		UpdatedAt:     updatedAt.String,
	}, nil
}

func (r *likeRepo) GetList(ctx context.Context, req *models.LikeGetListRequest) (*models.LikeGetListResponse, error) {

	var (
		resp   = &models.LikeGetListResponse{}
		query  string
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			count,
			publication_id,
			contributor_id,
			created_at,
			updated_at
		FROM likes
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.Search != "" {
		where += ` AND publication_id ILIKE '%' || '` + req.Search + `' || '%'`
	}

	query += where + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id            sql.NullString
			count         sql.NullFloat64
			publicationID sql.NullString
			contributorID sql.NullString
			createdAt     sql.NullString
			updatedAt     sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&count,
			&publicationID,
			&contributorID,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Likes = append(resp.Likes, &models.Like{
			Id:            id.String,
			Count:         count.Float64,
			PublicationID: publicationID.String,
			ContributorID: contributorID.String,
			CreatedAt:     createdAt.String,
			UpdatedAt:     updatedAt.String,
		})
	}

	return resp, nil
}

func (r *likeRepo) Update(ctx context.Context, req *models.UpdateLike) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			likes
		SET
			count = :count,
			publication_id = :publication_id,
			contributor_id = :contributor_id,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":             req.Id,
		"count":          req.Count,
		"publication_id": req.PublicationID,
		"contributor_id": req.ContributorID,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *likeRepo) Delete(ctx context.Context, req *models.LikePrimaryKey) error {

	_, err := r.db.Exec(ctx, "DELETE FROM likes WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return nil
}
