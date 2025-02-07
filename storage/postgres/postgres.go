package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"app/config"
	"app/storage"
)

type store struct {
	db           *pgxpool.Pool
	admin        *adminRepo
	user         *userRepo
	course       *courseRepo
	semester     *semesterRepo
	like         *likeRepo
	download     *downloadRepo
	publication  *publicationRepo
	notification *notificationRepo
}

func NewConnectionPostgres(cfg *config.Config) (storage.StorageI, error) {

	connect, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%d sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresDatabase,
		cfg.PostgresPassword,
		cfg.PostgresPort,
	))

	if err != nil {
		return nil, err
	}
	connect.MaxConns = cfg.PostgresMaxConnection

	pgxpool, err := pgxpool.ConnectConfig(context.Background(), connect)
	if err != nil {
		return nil, err
	}

	return &store{
		db: pgxpool,
	}, nil
}

func (s *store) Close() {
	s.db.Close()
}

func (s *store) Admin() storage.AdminRepoI {

	if s.admin == nil {
		s.admin = NewAdminRepo(s.db)
	}

	return s.admin
}

func (s *store) User() storage.UserRepoI {

	if s.user == nil {
		s.user = NewUserRepo(s.db)
	}

	return s.user
}
func (s *store) Course() storage.CourseRepoI {

	if s.course == nil {
		s.course = NewCourseRepo(s.db)
	}

	return s.course
}

func (s *store) Semester() storage.SemesterRepoI {

	if s.semester == nil {
		s.semester = NewSemesterRepo(s.db)
	}

	return s.semester
}
func (s *store) Like() storage.LikeRepoI {

	if s.like == nil {
		s.like = NewLikeRepo(s.db)
	}

	return s.like
}
func (s *store) Download() storage.DownloadRepoI {

	if s.download == nil {
		s.download = NewDownloadRepo(s.db)
	}

	return s.download
}
func (s *store) Publication() storage.PublicationRepoI {

	if s.publication == nil {
		s.publication = NewPublicationRepo(s.db)
	}

	return s.publication
}

func (s *store) Notification() storage.NotificationRepoI {

	if s.notification == nil {
		s.notification = NewNotificationRepo(s.db)
	}

	return s.notification
}
