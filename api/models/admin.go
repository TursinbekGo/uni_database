package models

type AdminPrimaryKey struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type CreateAdmin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Admin struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UpdateAdmin struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type AdminGetListResponse struct {
	Count  int      `json:"count"`
	Admins []*Admin `json:"admins"`
}
