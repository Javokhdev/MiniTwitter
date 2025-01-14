package entity

type Attachment struct {
	Id          string `json:"id"`
	TweetId     string `json:"-"`
	FilePath    string `json:"filepath"`
	ContentType string `json:"content_type"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type AttachmentList struct {
	Items []Attachment `json:"items"`
	Count int64        `json:"count"`
}

type AttachmentMultipleInsertRequest struct {
	TweetId     string       `json:"tweet_id"`
	Attachments []Attachment `json:"attachments"`
}

type Tweet struct {
	Id          string              `json:"id"`
	Owner       User                `json:"owner"`
	Content     string              `json:"content"`
	Tags        map[string][]string `json:"tags"`
	Attachments []Attachment        `json:"attachments"`
	Status      string              `json:"status"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
}

type TweetList struct {
	Items []Tweet `json:"items"`
	Count int64   `json:"count"`
}
