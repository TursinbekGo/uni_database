package models

type PublicationPrimaryKey struct {
	Id string `json:"id"`
}

type CreatePublication struct {
	CourseId      string `json:"course_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Tags          string `json:"tags"`
	ImageID       string `json:"image_id"`
	FileID        string `json:"file_id"`
	ContributorID string `json:"contributor_id"`
	Status        string `json:"status"`
}

type Publication struct {
	Id            string `json:"id"`
	CourseId      string `json:"course_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Tags          string `json:"tags"`
	ImageID       string `json:"image_id"`
	FileID        string `json:"file_id"`
	ContributorID string `json:"contributor_id"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type UpdatePublication struct {
	Id            string `json:"id"`
	CourseId      string `json:"course_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Tags          string `json:"tags"`
	ImageID       string `json:"image_id"`
	FileID        string `json:"file_id"`
	ContributorID string `json:"contributor_id"`
	Status        string `json:"status"`
}

type PublicationGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type PublicationGetListResponse struct {
	Count        int            `json:"count"`
	Publications []*Publication `json:"publications"`
}

type PublicationStats struct {
	LikeCount     float64 `json:"like_count"`
	DownloadCount float64 `json:"download_count"`
}
