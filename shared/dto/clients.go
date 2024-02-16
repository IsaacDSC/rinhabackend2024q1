package dto

type Client struct {
	ID      string `json:"id"`
	Limit   int64  `json:"limit"`
	Balance int64  `json:"balance"`
}
