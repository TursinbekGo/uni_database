package postgres

import (
	"context"
	"database/sql"
	"fmt"

	uuid "github.com/google/uuid"

	// "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"app/api/models"
	"app/pkg/helper"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(ctx context.Context, req *models.CreateUser) (string, error) {
	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO users(id, student_id, name, surname, email, grade, username, password,profile_image,status, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9,$10, NOW())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.StudentId,
		req.Name,
		req.Surname,
		req.Email,
		req.Grade,
		req.Username,
		req.Password,
		req.ProfileImage,
		req.Status,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *userRepo) GetByID(ctx context.Context, req *models.UserPrimaryKey) (*models.User, error) {

	var whereField = "id"
	if len(req.Email) > 0 {
		whereField = "email"
		req.Id = req.Email
	}

	var (
		query string

		id            sql.NullString
		studentId     sql.NullString
		name          sql.NullString
		surname       sql.NullString
		email         sql.NullString
		grade         sql.NullString
		username      sql.NullString
		password      sql.NullString
		profile_image sql.NullString
		status        sql.NullBool
		createdAt     sql.NullString
		updatedAt     sql.NullString
	)

	query = `
		SELECT
			id,
			student_id,
			name,
			surname,
			email,
			grade,
			username,
			password,
			profile_image,
			status,
			created_at,
			updated_at
		FROM users
	WHERE ` + whereField + ` = $1

	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&studentId,
		&name,
		&surname,
		&email,
		&grade,
		&username,
		&password,
		&profile_image,
		&status,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.User{
		Id:           id.String,
		StudentId:    studentId.String,
		Name:         name.String,
		Surname:      surname.String,
		Email:        email.String,
		Grade:        grade.String,
		Username:     username.String,
		Password:     password.String,
		ProfileImage: profile_image.String,
		Status:       status.Bool,
		CreatedAt:    createdAt.String,
		UpdatedAt:    updatedAt.String,
	}, nil
}

func (r *userRepo) GetList(ctx context.Context, req *models.UserGetListRequest) (*models.UserGetListResponse, error) {
	var (
		resp   = &models.UserGetListResponse{}
		query  string
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			student_id,
			name,
			surname,
			email,
			grade,
			username,
			password,
		    profile_image,
			status,
			created_at,
			updated_at
		FROM users
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.Search != "" {
		where += ` AND (name ILIKE '%' || '` + req.Search + `' || '%' OR surname ILIKE '%' || '` + req.Search + `' || '%' OR email ILIKE '%' || '` + req.Search + `' || '%')`
	}

	query += where + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id            sql.NullString
			studentId     sql.NullString
			name          sql.NullString
			surname       sql.NullString
			email         sql.NullString
			grade         sql.NullString
			username      sql.NullString
			password      sql.NullString
			profile_image sql.NullString
			status        sql.NullBool
			createdAt     sql.NullString
			updatedAt     sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&studentId,
			&name,
			&surname,
			&email,
			&grade,
			&username,
			&password,
			&profile_image,
			&status,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Users = append(resp.Users, &models.User{
			Id:           id.String,
			StudentId:    studentId.String,
			Name:         name.String,
			Surname:      surname.String,
			Email:        email.String,
			Grade:        grade.String,
			Username:     username.String,
			Password:     password.String,
			ProfileImage: profile_image.String,
			Status:       status.Bool,
			CreatedAt:    createdAt.String,
			UpdatedAt:    updatedAt.String,
		})
	}

	return resp, nil
}

func (r *userRepo) Update(ctx context.Context, req *models.UpdateUser) (int64, error) {
	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			users
		SET
			student_id = :student_id,
			name = :name,
			surname = :surname,
			email = :email,
			grade = :grade,
			username = :username,
			password = :password,
			profile_image = :profile_image,
			status = :status,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":            req.Id,
		"student_id":    req.StudentId,
		"name":          req.Name,
		"surname":       req.Surname,
		"email":         req.Email,
		"grade":         req.Grade,
		"username":      req.Username,
		"password":      req.Password,
		"profile_image": req.ProfileImage,
		"status":        req.Status,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *userRepo) Delete(ctx context.Context, req *models.UserPrimaryKey) error {
	_, err := r.db.Exec(ctx, "DELETE FROM users WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) GetUserActivityCounts(ctx context.Context, userID string) (*models.UserActivityCounts, error) {
	var (
		publicationCount sql.NullInt64
		downloadCount    sql.NullFloat64
		likeCount        sql.NullFloat64
	)

	query := `
        SELECT
            (SELECT COUNT(*) FROM publications WHERE contributor_id = $1) AS publication_count,
    (SELECT count FROM downloads WHERE contributor_id = $1) AS download_count,
    (SELECT count FROM likes WHERE contributor_id = $1) AS like_count;
    `

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&publicationCount,
		&downloadCount,
		&likeCount,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user activity counts: %v", err)
	}

	return &models.UserActivityCounts{
		PublicationCount: int(publicationCount.Int64),
		DownloadCount:    (downloadCount.Float64),
		LikeCount:        (likeCount.Float64),
	}, nil
}

func (r *userRepo) GetTopContributors(ctx context.Context) ([]*models.UserpublicationCounts, error) {
	var (
		query = `
			SELECT
				u.id AS user_id,
				u.name,
				COUNT(p.id) AS publication_count
			FROM publications p
			JOIN users u ON p.contributor_id = u.id
			GROUP BY u.id
			ORDER BY publication_count DESC
		`
		// rows *pgx.Rows
	)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contributors []*models.UserpublicationCounts
	for rows.Next() {
		var userID, name string
		var publicationCount int

		err := rows.Scan(&userID, &name, &publicationCount)
		if err != nil {
			return nil, err
		}

		contributors = append(contributors, &models.UserpublicationCounts{
			UserID:           userID,
			Name:             name,
			PublicationCount: publicationCount,
		})
	}

	return contributors, nil
}

func (r *userRepo) GetUserScores(ctx context.Context) ([]*models.UserScore, error) {
	var scores []*models.UserScore

	query := `
     SELECT 
    u.id AS user_id,
    ROUND((COALESCE(SUM(l.count), 0) + COALESCE(SUM(d.count), 0)) / 
          CASE 
              WHEN COUNT(p.id) = 0 THEN 1 
              ELSE COUNT(p.id) 
          END) AS score  -- Round to nearest integer
FROM users u
LEFT JOIN publications p ON p.contributor_id = u.id
LEFT JOIN likes l ON l.publication_id = p.id
LEFT JOIN downloads d ON d.publication_id = p.id
GROUP BY u.id
ORDER BY score DESC
    `

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get user scores: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var score models.UserScore
		err := rows.Scan(
			&score.UserID,
			&score.Score,
		)
		if err != nil {
			return nil, err
		}
		scores = append(scores, &score)
	}

	return scores, nil
}

func (r *userRepo) GetUserRank(ctx context.Context, userID string) (*models.UserRank, error) {
	var rankInfo models.UserRank

	query := `
			SELECT 
    (SELECT COUNT(*) FROM users) AS user_count,  -- Total user count (calculated separately)
    (
        SELECT COUNT(*) + 1  -- Calculate rank
        FROM (
            SELECT 
                u.id AS user_id,
                (COALESCE(SUM(l.count), 0) + COALESCE(SUM(d.count), 0)) / 
                CASE 
                    WHEN COUNT(p.id) = 0 THEN 1 
                    ELSE COUNT(p.id) 
                END AS score
            FROM users u
            LEFT JOIN publications p ON p.contributor_id = u.id
            LEFT JOIN likes l ON l.publication_id = p.id
            LEFT JOIN downloads d ON d.publication_id = p.id
            GROUP BY u.id
            ORDER BY score DESC
        ) AS ranked_users
        WHERE score > (
            SELECT 
                (COALESCE(SUM(l.count), 0) + COALESCE(SUM(d.count), 0)) / 
                CASE 
                    WHEN COUNT(p.id) = 0 THEN 1 
                    ELSE COUNT(p.id) 
                END AS score
            FROM users u
            LEFT JOIN publications p ON p.contributor_id = u.id
            LEFT JOIN likes l ON l.publication_id = p.id
            LEFT JOIN downloads d ON d.publication_id = p.id
            WHERE u.id = $1  -- Score of the specific user
        )
    ) AS user_rank;
	`

	err := r.db.QueryRow(ctx, query, userID).Scan(&rankInfo.UserCount, &rankInfo.UserRank)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err) // Handle user not found
		}
		return nil, fmt.Errorf("failed to get user rank: %w", err)
	}

	return &rankInfo, nil
}
func (r *userRepo) GetUserStatistics(ctx context.Context, userID string) (*models.UserStatistics, error) {
	// Query to get count of publications, likes, and downloads for a specific user
	query := `
        SELECT 
            u.id AS user_id,
            COUNT(p.id) AS publications_count,
            COALESCE(SUM(l.count), 0) AS likes_count,
            COALESCE(SUM(d.count), 0) AS downloads_count
        FROM users u
        LEFT JOIN publications p ON p.contributor_id = u.id
        LEFT JOIN likes l ON l.publication_id = p.id
        LEFT JOIN downloads d ON d.publication_id = p.id
        WHERE u.id = $1
        GROUP BY u.id;
    `

	// Execute the query
	var stats models.UserStatistics
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&stats.UserID,
		&stats.PublicationsCount,
		&stats.LikesCount,
		&stats.DownloadsCount,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user statistics: %v", err)
	}

	// Now, calculate the rank of the user
	rank, err := r.GetUserRank(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate user rank: %v", err)
	}
	stats.Rank = rank.UserRank

	return &stats, nil
}
