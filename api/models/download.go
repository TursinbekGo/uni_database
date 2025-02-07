package models

type DownloadPrimaryKey struct {
	Id string `json:"id"`
}

type CreateDownload struct {
	Count         float64 `json:"count"`
	PublicationID string  `json:"publication_id"`
	ContributorID string  `json:"contributor_id"`
}

type Download struct {
	Id            string  `json:"id"`
	Count         float64 `json:"count"`
	PublicationID string  `json:"publication_id"`
	ContributorID string  `json:"contributor_id"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type UpdateDownload struct {
	Id            string  `json:"id"`
	Count         float64 `json:"count"`
	PublicationID string  `json:"publication_id"`
	ContributorID string  `json:"contributor_id"`
}

type DownloadGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type DownloadGetListResponse struct {
	Count     int         `json:"count"`
	Downloads []*Download `json:"downloads"`
}
