package models

type Pagination struct {
	Limit        int         `json:"limit"`
	Page         int         `json:"page"`
	Sort         string      `json:"sort"`
	TotalRows    int64       `json:"total_rows"`
	TotalPages   int         `json:"total_pages"`
	PreviousPage string      `json:"previous_page"`
	NextPage     string      `json:"next_page"`
	FirstPage    string      `json:"first_page"`
	LastPage     string      `json:"last_page"`
	FromRow      int         `json:"from_row"`
	ToRow        int         `json:"to_row"`
	Rows         interface{} `json:"rows"`
}
