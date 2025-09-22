package utils

type PaginationRequest struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	OrderBy string `json:"order_by"`
	SortBy  string `json:"sort_by"`
	Search  string `json:"search"`
}

type PaginationResponse struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	OrderBy   string `json:"order_by"`
	SortBy    string `json:"sort_by"`
	Search    string `json:"search"`
	TotalData int64  `json:"total_data"`
	TotalPage int    `json:"total_page"`
}
