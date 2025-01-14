package entity

type Tag struct {
	Id        string `json:"id"`
	Slug      string `json:"slug"`
	Level     int    `json:"level"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TagList struct {
	Items []Tag `json:"items"`
	Count int   `json:"count"`
}

type UserTag struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	Tag       Tag    `json:"tag"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserTagList struct {
	Items []UserTag `json:"items"`
	Count int       `json:"count"`
}
