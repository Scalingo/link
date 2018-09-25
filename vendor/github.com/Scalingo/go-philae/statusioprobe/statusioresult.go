package statusioprobe

import "time"

type StatusIOResponse struct {
	Result *StatusIOResult `json:"result"`
}

type StatusIOResult struct {
	Overall *StatusIOOverallResult  `json:"status_overall"`
	Status  *[]StatusIOStatusResult `json:"status"`
}

type StatusIOOverallResult struct {
	Updated    time.Time `json:"updated"`
	Status     string    `json:"status"`
	StatusCode int       `json:"status_code"`
}

type StatusIOStatusResult struct {
	ID         string                     `json:"id"`
	Name       string                     `json:"name"`
	Updated    time.Time                  `json:"updated"`
	Status     string                     `json:"status"`
	StatusCode int                        `json:"status_code"`
	Containers []*StatusIOContainerResult `json:"containers"`
}

type StatusIOContainerResult struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Updated    time.Time `json:"updated"`
	Status     string    `json:"status"`
	StatusCode int       `json:"status_code"`
}
