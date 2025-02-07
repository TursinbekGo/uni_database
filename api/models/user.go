package models

type UserPrimaryKey struct {
	Id        string `json:"id"`
	StudentId string `json:"student_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Status    string `json:"status"`
}

type CreateUser struct {
	StudentId    string `json:"student_id"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Email        string `json:"email"`
	Grade        string `json:"grade"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ProfileImage string `json:"profile_image"`
	Status       bool   `json:"status"`
}

type User struct {
	Id           string `json:"id"`
	StudentId    string `json:"student_id"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Email        string `json:"email"`
	Grade        string `json:"grade"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ProfileImage string `json:"profile_image"`
	Status       bool   `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type UpdateUser struct {
	Id           string `json:"id"`
	StudentId    string `json:"student_id"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Email        string `json:"email"`
	Grade        string `json:"grade"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ProfileImage string `json:"profile_image"`
	Status       bool   `json:"status"`
}

type UserGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type UserGetListResponse struct {
	Count int     `json:"count"`
	Users []*User `json:"users"`
}

type UserActivityCounts struct {
	PublicationCount int     `json:"publication_count"`
	DownloadCount    float64 `json:"download_count"`
	LikeCount        float64 `json:"like_count"`
}
type UserpublicationCounts struct {
	UserID           string `json:"user_id"`
	Name             string `json:"name"`
	PublicationCount int    `json:"publication_count"`
}

type UserScore struct {
	UserID string `json:"user_id"`
	Score  int    `json:"score"`
}
type UserRank struct {
	UserCount int `json:"user_count"`
	UserRank  int `json:"user_rank"`
}
type UserStatistics struct {
	UserID            string `json:"user_id"`
	PublicationsCount int    `json:"publications_count"`
	LikesCount        int    `json:"likes_count"`
	DownloadsCount    int    `json:"downloads_count"`
	Rank              int    `json:"rank"`
}
