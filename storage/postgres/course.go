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

type courseRepo struct {
	db *pgxpool.Pool
}

func NewCourseRepo(db *pgxpool.Pool) *courseRepo {
	return &courseRepo{
		db: db,
	}
}

func (r *courseRepo) Create(ctx context.Context, req *models.CreateCourse) (string, error) {
	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO courses(id, course_title, semester_id, updated_at)
		VALUES ($1, $2, $3,NOW())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.CourseTitle,
		req.SemesterID,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *courseRepo) GetByID(ctx context.Context, req *models.CoursePrimaryKey) (*models.Course, error) {

	var (
		query string

		id          sql.NullString
		courseTitle sql.NullString
		semesterID  sql.NullString
		createdAt   sql.NullString
		updatedAt   sql.NullString
	)

	query = `
		SELECT
			id,
			course_title,
			semester_id,
			created_at,
			updated_at
		FROM courses
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&courseTitle,
		&semesterID,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Course{
		Id:          id.String,
		CourseTitle: courseTitle.String,
		SemesterID:  semesterID.String,
		CreatedAt:   createdAt.String,
		UpdatedAt:   updatedAt.String,
	}, nil
}

func (r *courseRepo) GetList(ctx context.Context, req *models.CourseGetListRequest) (*models.CourseGetListResponse, error) {

	var (
		resp   = &models.CourseGetListResponse{}
		query  string
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			course_title,
			semester_id,
			created_at,
			updated_at
		FROM courses
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.Search != "" {
		where += ` AND course_title ILIKE '%' || '` + req.Search + `' || '%'`
	}

	query += where + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id          sql.NullString
			courseTitle sql.NullString
			semesterID  sql.NullString
			createdAt   sql.NullString
			updatedAt   sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&courseTitle,
			&semesterID,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Courses = append(resp.Courses, &models.Course{
			Id:          id.String,
			CourseTitle: courseTitle.String,
			SemesterID:  semesterID.String,
			CreatedAt:   createdAt.String,
			UpdatedAt:   updatedAt.String,
		})
	}

	return resp, nil
}

func (r *courseRepo) Update(ctx context.Context, req *models.UpdateCourse) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			courses
		SET
			course_title = :course_title,
			semester_id = :semester_id,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":           req.Id,
		"course_title": req.CourseTitle,
		"semester_id":  req.SemesterID,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *courseRepo) Delete(ctx context.Context, req *models.CoursePrimaryKey) error {

	_, err := r.db.Exec(ctx, "DELETE FROM courses WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return nil
}
