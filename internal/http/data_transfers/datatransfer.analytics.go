package data_transfers

import "time"

type DayWiseAnalyticsResponse struct {
	Date         time.Time           `json:"date"`
	Sets         int                 `json:"sets"`
	Reps         int                 `json:"reps"`
	Exercises    int                 `json:"exercises"`
	SessionsTime string              `json:"session_time"`
	Details      map[string][]string `json:"details"`
	Sessions     []SessionResponse   `json:"sessions"`
}
