package storage

import (
	"app/api/models"
	"context"
)

type StorageI interface {
	Close()
	Admin() AdminRepoI
	User() UserRepoI
	Course() CourseRepoI
	Like() LikeRepoI
	Semester() SemesterRepoI
	Download() DownloadRepoI
	Publication() PublicationRepoI
	Notification() NotificationRepoI
}

type AdminRepoI interface {
	Create(context.Context, *models.CreateAdmin) (string, error)
	GetByID(context.Context, *models.AdminPrimaryKey) (*models.Admin, error)
	GetList(context.Context, *models.AdminGetListRequest) (*models.AdminGetListResponse, error)
	Update(context.Context, *models.UpdateAdmin) (int64, error)
	Delete(context.Context, *models.AdminPrimaryKey) error
}

type UserRepoI interface {
	Create(context.Context, *models.CreateUser) (string, error)
	GetByID(context.Context, *models.UserPrimaryKey) (*models.User, error)
	GetList(context.Context, *models.UserGetListRequest) (*models.UserGetListResponse, error)
	Update(context.Context, *models.UpdateUser) (int64, error)
	Delete(context.Context, *models.UserPrimaryKey) error
	GetUserActivityCounts(ctx context.Context, userID string) (*models.UserActivityCounts, error)
	GetTopContributors(ctx context.Context) ([]*models.UserpublicationCounts, error)
	GetUserScores(ctx context.Context) ([]*models.UserScore, error)
	GetUserRank(ctx context.Context, userID string) (*models.UserRank, error)
	GetUserStatistics(ctx context.Context, userID string) (*models.UserStatistics, error)
}

type CourseRepoI interface {
	Create(context.Context, *models.CreateCourse) (string, error)
	GetByID(context.Context, *models.CoursePrimaryKey) (*models.Course, error)
	GetList(context.Context, *models.CourseGetListRequest) (*models.CourseGetListResponse, error)
	Update(context.Context, *models.UpdateCourse) (int64, error)
	Delete(context.Context, *models.CoursePrimaryKey) error
}

type SemesterRepoI interface {
	Create(context.Context, *models.CreateSemester) (string, error)
	GetByID(context.Context, *models.SemesterPrimaryKey) (*models.Semester, error)
	GetList(context.Context, *models.SemesterGetListRequest) (*models.SemesterGetListResponse, error)
	Update(context.Context, *models.UpdateSemester) (int64, error)
	Delete(context.Context, *models.SemesterPrimaryKey) error
}

type LikeRepoI interface {
	Create(context.Context, *models.CreateLike) (string, error)
	GetByID(context.Context, *models.LikePrimaryKey) (*models.Like, error)
	GetList(context.Context, *models.LikeGetListRequest) (*models.LikeGetListResponse, error)
	Update(context.Context, *models.UpdateLike) (int64, error)
	Delete(context.Context, *models.LikePrimaryKey) error
}
type DownloadRepoI interface {
	Create(context.Context, *models.CreateDownload) (string, error)
	GetByID(context.Context, *models.DownloadPrimaryKey) (*models.Download, error)
	GetList(context.Context, *models.DownloadGetListRequest) (*models.DownloadGetListResponse, error)
	Update(context.Context, *models.UpdateDownload) (int64, error)
	Delete(context.Context, *models.DownloadPrimaryKey) error
}
type PublicationRepoI interface {
	Create(context.Context, *models.CreatePublication) (string, error)
	GetByID(context.Context, *models.PublicationPrimaryKey) (*models.Publication, error)
	GetList(context.Context, *models.PublicationGetListRequest) (*models.PublicationGetListResponse, error)
	Update(context.Context, *models.UpdatePublication) (int64, error)
	Delete(context.Context, *models.PublicationPrimaryKey) error
	GetPublicationStats(ctx context.Context, publicationID string) (*models.PublicationStats, error)
	GetPublicationsByTag(ctx context.Context, tag string) ([]*models.Publication, error)
}
type NotificationRepoI interface {
	Create(context.Context, *models.CreateNotification) (string, error)
	GetByID(context.Context, *models.NotificationPrimaryKey) (*models.Notification, error)
	GetList(context.Context, *models.NotificationGetListRequest) (*models.NotificationGetListResponse, error)
	Update(context.Context, *models.UpdateNotification) (int64, error)
	Delete(context.Context, *models.NotificationPrimaryKey) error
}
