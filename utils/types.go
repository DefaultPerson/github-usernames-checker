package utils

type User struct {
	Username string
	Filename string
}

type GitHubAccount struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	AuthToken       string `json:"auth_token"`
	Cookie          string `json:"cookie"`
	Timestamp       string `json:"timestamp"`
	TimestampSecret string `json:"timestamp_secret"`
}
