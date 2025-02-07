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

type adminRepo struct {
	db *pgxpool.Pool
}

func NewAdminRepo(db *pgxpool.Pool) *adminRepo {
	return &adminRepo{
		db: db,
	}
}

func (r *adminRepo) Create(ctx context.Context, req *models.CreateAdmin) (string, error) {

	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO admins(id, email,password, updated_at)
		VALUES ($1, $2, $3,NOW())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.Email,
		req.Password,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *adminRepo) GetByID(ctx context.Context, req *models.AdminPrimaryKey) (*models.Admin, error) {

	var whereField = "id"
	if len(req.Email) > 0 {
		whereField = "email"
		req.Id = req.Email
	}
	var (
		query string

		id        sql.NullString
		email     sql.NullString
		password  sql.NullString
		createdAt sql.NullString
		updatedAt sql.NullString
	)

	query = `
		SELECT
			id,
			email,
			password,
			created_at,
			updated_at
		FROM admins
		WHERE ` + whereField + ` = $1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&email,
		&password,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Admin{
		Id:        id.String,
		Email:     email.String,
		Password:  password.String,
		CreatedAt: createdAt.String,
		UpdatedAt: updatedAt.String,
	}, nil
}

func (r *adminRepo) GetList(ctx context.Context, req *models.AdminGetListRequest) (*models.AdminGetListResponse, error) {

	var (
		resp   = &models.AdminGetListResponse{}
		query  string
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			email,
			password,
			created_at,
			updated_at
		FROM admins
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.Search != "" {
		where += ` AND email ILIKE '%' || '` + req.Search + `' || '%'`
	}

	query += where + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id        sql.NullString
			email     sql.NullString
			password  sql.NullString
			createdAt sql.NullString
			updatedAt sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&email,
			&password,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Admins = append(resp.Admins, &models.Admin{
			Id:        id.String,
			Email:     email.String,
			Password:  password.String,
			CreatedAt: createdAt.String,
			UpdatedAt: updatedAt.String,
		})
	}

	return resp, nil
}

func (r *adminRepo) Update(ctx context.Context, req *models.UpdateAdmin) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			admins
		SET
		id = :id,
		email = :email,
		password = :password,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":       req.Id,
		"email":    req.Email,
		"password": req.Password,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *adminRepo) Delete(ctx context.Context, req *models.AdminPrimaryKey) error {

	_, err := r.db.Exec(ctx, "DELETE FROM admins WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return nil
}
