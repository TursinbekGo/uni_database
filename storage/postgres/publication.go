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

type publicationRepo struct {
	db *pgxpool.Pool
}

func NewPublicationRepo(db *pgxpool.Pool) *publicationRepo {
	return &publicationRepo{
		db: db,
	}
}

func (r *publicationRepo) Create(ctx context.Context, req *models.CreatePublication) (string, error) {
	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO publications(id, course_id, title, description,tags, image_id, file_id, contributor_id, status, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9, NOW())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.CourseId,
		req.Title,
		req.Description,
		req.Tags,
		req.ImageID,
		req.FileID,
		req.ContributorID,
		req.Status,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *publicationRepo) GetByID(ctx context.Context, req *models.PublicationPrimaryKey) (*models.Publication, error) {

	var (
		query string

		id          sql.NullString
		courseId    sql.NullString
		title       sql.NullString
		description sql.NullString
		tags        sql.NullString
		imageID     sql.NullString
		fileID      sql.NullString
		contributor sql.NullString
		status      sql.NullString
		createdAt   sql.NullString
		updatedAt   sql.NullString
	)

	query = `
		SELECT
			id,
			course_id,
			title,
			description,
			tags,
			image_id,
			file_id,
			contributor_id,
			status,
			created_at,
			updated_at
		FROM publications
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&courseId,
		&title,
		&description,
		&tags,
		&imageID,
		&fileID,
		&contributor,
		&status,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Publication{
		Id:            id.String,
		CourseId:      courseId.String,
		Title:         title.String,
		Description:   description.String,
		Tags:          tags.String,
		ImageID:       imageID.String,
		FileID:        fileID.String,
		ContributorID: contributor.String,
		Status:        status.String,
		CreatedAt:     createdAt.String,
		UpdatedAt:     updatedAt.String,
	}, nil
}

func (r *publicationRepo) GetList(ctx context.Context, req *models.PublicationGetListRequest) (*models.PublicationGetListResponse, error) {

	var (
		resp   = &models.PublicationGetListResponse{}
		query  string
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			course_id,
			title,
			description,
			tags,
			image_id,
			file_id,
			contributor_id,
			status,
			created_at,
			updated_at
		FROM publications
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.Search != "" {
		where += ` AND (title ILIKE '%' || '` + req.Search + `' || '%' OR tags ILIKE '%' || '` + req.Search + `' || '%' OR status ILIKE '%' || '` + req.Search + `' || '%')`
	}

	query += where + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id          sql.NullString
			courseId    sql.NullString
			title       sql.NullString
			description sql.NullString
			tags        sql.NullString
			imageID     sql.NullString
			fileID      sql.NullString
			contributor sql.NullString
			status      sql.NullString
			createdAt   sql.NullString
			updatedAt   sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&courseId,
			&title,
			&description,
			&tags,
			&imageID,
			&fileID,
			&contributor,
			&status,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Publications = append(resp.Publications, &models.Publication{
			Id:            id.String,
			CourseId:      courseId.String,
			Title:         title.String,
			Description:   description.String,
			Tags:          tags.String,
			ImageID:       imageID.String,
			FileID:        fileID.String,
			ContributorID: contributor.String,
			Status:        status.String,
			CreatedAt:     createdAt.String,
			UpdatedAt:     updatedAt.String,
		})
	}

	return resp, nil
}

func (r *publicationRepo) Update(ctx context.Context, req *models.UpdatePublication) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			publications
		SET
			course_id = :course_id,
			title = :title,
			description = :description,
			tags = :tags,
			image_id = :image_id,
			file_id = :file_id,
			contributor_id = :contributor_id,
			status = :status,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":             req.Id,
		"course_id":      req.CourseId,
		"title":          req.Title,
		"description":    req.Description,
		"tags":           req.Tags,
		"image_id":       req.ImageID,
		"file_id":        req.FileID,
		"contributor_id": req.ContributorID,
		"status":         req.Status,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *publicationRepo) Delete(ctx context.Context, req *models.PublicationPrimaryKey) error {

	_, err := r.db.Exec(ctx, "DELETE FROM publications WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return nil
}
func (r *publicationRepo) GetPublicationStats(ctx context.Context, publicationID string) (*models.PublicationStats, error) {
	var (
		likeCount     sql.NullFloat64
		downloadCount sql.NullFloat64
	)

	query := `
        SELECT
            (SELECT count FROM likes WHERE publication_id = $1) AS like_count,
            (SELECT count FROM downloads WHERE publication_id = $1) AS download_count;
    `

	err := r.db.QueryRow(ctx, query, publicationID).Scan(
		&likeCount,
		&downloadCount,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get publication stats: %v", err)
	}

	return &models.PublicationStats{
		LikeCount:     (likeCount.Float64),
		DownloadCount: (downloadCount.Float64),
	}, nil
}
func (r *publicationRepo) GetPublicationsByTag(ctx context.Context, tag string) ([]*models.Publication, error) {
	var publications []*models.Publication

	query := `
		SELECT id, course_id, title, description, tags, image_id, file_id, contributor_id, status
		FROM publications
		WHERE tags ILIKE '%' || $1 || '%';
	`

	rows, err := r.db.Query(ctx, query, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to get publications by tag: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var publication models.Publication
		err := rows.Scan(
			&publication.Id,
			&publication.CourseId,
			&publication.Title,
			&publication.Description,
			&publication.Tags,
			&publication.ImageID,
			&publication.FileID,
			&publication.ContributorID,
			&publication.Status,
			// &publication.CreatedAt,
			// &publication.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		publications = append(publications, &publication)
	}

	return publications, nil
}
