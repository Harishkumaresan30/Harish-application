package models

type Metric struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Time  int64   `json:"time"` // Unix timestamp
}
