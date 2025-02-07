package models

type CoursePrimaryKey struct {
	Id string `json:"id"`
}

type CreateCourse struct {
	CourseTitle string `json:"course_title"`
	SemesterID  string `json:"semester_id"`
}

type Course struct {
	Id          string `json:"id"`
	CourseTitle string `json:"course_title"`
	SemesterID  string `json:"semester_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UpdateCourse struct {
	Id          string `json:"id"`
	CourseTitle string `json:"course_title"`
	SemesterID  string `json:"semester_id"`
}

type CourseGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type CourseGetListResponse struct {
	Count   int       `json:"count"`
	Courses []*Course `json:"courses"`
}
