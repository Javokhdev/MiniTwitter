package entity

type User struct {
	ID          string `json:"id"`
	FullName    string `json:"full_name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	UserType    string `json:"user_type"`
	UserRole    string `json:"user_role"`
	Status      string `json:"status"`
	AccessToken string `json:"access_token"`
	AvatarId    string `json:"avatar_id"`
	Gender      string `json:"gender"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UserSingleRequest struct {
	ID       string `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

type UserList struct {
	Items []User `json:"users"`
	Count int    `json:"count"`
}
