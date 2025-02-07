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

type notificationRepo struct {
	db *pgxpool.Pool
}

func NewNotificationRepo(db *pgxpool.Pool) *notificationRepo {
	return &notificationRepo{
		db: db,
	}
}

func (r *notificationRepo) Create(ctx context.Context, req *models.CreateNotification) (string, error) {
	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO notifications(id, user_image, message, file, username, message_type, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6,NOW())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.UserImage,
		req.Message,
		req.File,
		req.UserName,
		req.MessageType,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *notificationRepo) GetByID(ctx context.Context, req *models.NotificationPrimaryKey) (*models.Notification, error) {

	var (
		query string

		id          sql.NullString
		userImage   sql.NullString
		message     sql.NullString
		file        sql.NullString
		userName    sql.NullString
		messageType sql.NullString
		createdAt   sql.NullString
		updatedAt   sql.NullString
	)

	query = `
		SELECT
			id,
			user_image,
			message,
			file,
			username,
			message_type,
			created_at,
			updated_at
		FROM notifications
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&userImage,
		&message,
		&file,
		&userName,
		&messageType,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Notification{
		Id:          id.String,
		UserImage:   userImage.String,
		Message:     message.String,
		File:        file.String,
		UserName:    userName.String,
		MessageType: messageType.String,
		CreatedAt:   createdAt.String,
		UpdatedAt:   updatedAt.String,
	}, nil
}

func (r *notificationRepo) GetList(ctx context.Context, req *models.NotificationGetListRequest) (*models.NotificationGetListResponse, error) {

	var (
		resp   = &models.NotificationGetListResponse{}
		query  string
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			user_image,
			message,
			file,
			username,
			message_type,
			created_at,
			updated_at
		FROM notifications
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.Search != "" {
		where += ` AND message ILIKE '%' || '` + req.Search + `' || '%'`
	}

	query += where + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id          sql.NullString
			userImage   sql.NullString
			message     sql.NullString
			file        sql.NullString
			userName    sql.NullString
			messageType sql.NullString
			createdAt   sql.NullString
			updatedAt   sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&userImage,
			&message,
			&file,
			&userName,
			&messageType,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Notifications = append(resp.Notifications, &models.Notification{
			Id:          id.String,
			UserImage:   userImage.String,
			Message:     message.String,
			File:        file.String,
			UserName:    userName.String,
			MessageType: messageType.String,
			CreatedAt:   createdAt.String,
			UpdatedAt:   updatedAt.String,
		})
	}

	return resp, nil
}

func (r *notificationRepo) Update(ctx context.Context, req *models.UpdateNotification) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			notifications
		SET
			user_image = :user_image,
			message = :message,
			file = :file,
			username = :username,
			message_type = :message_type,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":           req.Id,
		"user_image":   req.UserImage,
		"message":      req.Message,
		"file":         req.File,
		"username":     req.UserName,
		"message_type": req.MessageType,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *notificationRepo) Delete(ctx context.Context, req *models.NotificationPrimaryKey) error {

	_, err := r.db.Exec(ctx, "DELETE FROM notifications WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return nil
}
