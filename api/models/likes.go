package models

type LikePrimaryKey struct {
	Id string `json:"id"`
}

type CreateLike struct {
	Count         float64 `json:"count"`
	PublicationID string  `json:"publication_id"`
	ContributorID string  `json:"contributor_id"`
}

type Like struct {
	Id            string  `json:"id"`
	Count         float64 `json:"count"`
	PublicationID string  `json:"publication_id"`
	ContributorID string  `json:"contributor_id"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type UpdateLike struct {
	Id            string  `json:"id"`
	Count         float64 `json:"count"`
	PublicationID string  `json:"publication_id"`
	ContributorID string  `json:"contributor_id"`
}

type LikeGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type LikeGetListResponse struct {
	Count int     `json:"count"`
	Likes []*Like `json:"likes"`
}
