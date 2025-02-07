package models

type SemesterPrimaryKey struct {
	Id string `json:"id"`
}

type CreateSemester struct {
	SemesterNumber string `json:"semester_number"`
}

type Semester struct {
	Id             string `json:"id"`
	SemesterNumber string `json:"semester_number"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type UpdateSemester struct {
	Id             string `json:"id"`
	SemesterNumber string `json:"semester_number"`
}

type SemesterGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type SemesterGetListResponse struct {
	Count     int         `json:"count"`
	Semesters []*Semester `json:"semesters"`
}
