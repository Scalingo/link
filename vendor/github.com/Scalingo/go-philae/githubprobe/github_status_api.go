package githubprobe

import "time"

type GithubStatusResponse struct {
	Status        string    `json:"status"`
	LastUpdatedAt time.Time `json:"last_updated"`
}
