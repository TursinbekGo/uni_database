package models

type NotificationPrimaryKey struct {
	Id string `json:"id"`
}

type CreateNotification struct {
	UserImage   string `json:"user_image"`
	Message     string `json:"message"`
	File        string `json:"file"`
	UserName    string `json:"username"`
	MessageType string `json:"message_type"`
}

type Notification struct {
	Id          string `json:"id"`
	UserImage   string `json:"user_image"`
	Message     string `json:"message"`
	File        string `json:"file"`
	UserName    string `json:"username"`
	MessageType string `json:"message_type"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UpdateNotification struct {
	Id          string `json:"id"`
	UserImage   string `json:"user_image"`
	Message     string `json:"message"`
	File        string `json:"file"`
	UserName    string `json:"username"`
	MessageType string `json:"message_type"`
}

type NotificationGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type NotificationGetListResponse struct {
	Count         int             `json:"count"`
	Notifications []*Notification `json:"notifications"`
}
