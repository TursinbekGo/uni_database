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

type semesterRepo struct {
	db *pgxpool.Pool
}

func NewSemesterRepo(db *pgxpool.Pool) *semesterRepo {
	return &semesterRepo{
		db: db,
	}
}

func (r *semesterRepo) Create(ctx context.Context, req *models.CreateSemester) (string, error) {
	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO semesters(id, semester_number, updated_at)
		VALUES ($1, $2, NOW())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.SemesterNumber,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *semesterRepo) GetByID(ctx context.Context, req *models.SemesterPrimaryKey) (*models.Semester, error) {

	var (
		query string

		id             sql.NullString
		semesterNumber sql.NullString
		createdAt      sql.NullString
		updatedAt      sql.NullString
	)

	query = `
		SELECT
			id,
			semester_number,
			created_at,
			updated_at
		FROM semesters
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&semesterNumber,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Semester{
		Id:             id.String,
		SemesterNumber: semesterNumber.String,
		CreatedAt:      createdAt.String,
		UpdatedAt:      updatedAt.String,
	}, nil
}

func (r *semesterRepo) GetList(ctx context.Context, req *models.SemesterGetListRequest) (*models.SemesterGetListResponse, error) {

	var (
		resp   = &models.SemesterGetListResponse{}
		query  string
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			semester_number,
			created_at,
			updated_at
		FROM semesters
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.Search != "" {
		where += ` AND semester_number ILIKE '%' || '` + req.Search + `' || '%'`
	}

	query += where + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id             sql.NullString
			semesterNumber sql.NullString
			createdAt      sql.NullString
			updatedAt      sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&semesterNumber,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Semesters = append(resp.Semesters, &models.Semester{
			Id:             id.String,
			SemesterNumber: semesterNumber.String,
			CreatedAt:      createdAt.String,
			UpdatedAt:      updatedAt.String,
		})
	}

	return resp, nil
}

func (r *semesterRepo) Update(ctx context.Context, req *models.UpdateSemester) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			semesters
		SET
			semester_number = :semester_number,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":              req.Id,
		"semester_number": req.SemesterNumber,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *semesterRepo) Delete(ctx context.Context, req *models.SemesterPrimaryKey) error {

	_, err := r.db.Exec(ctx, "DELETE FROM semesters WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return nil
}
