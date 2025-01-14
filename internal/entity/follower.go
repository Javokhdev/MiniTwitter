package entity

type Follower struct {
	FollowingId string `json:"follwing_id"`
	FollowerId  string `json:"follower_id"`
	UnFollowed    bool   `json:"followed"`
}
