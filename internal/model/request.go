package model

type Paginated struct {
	Page       int64 `json:"page"`
	Size       int64 `json:"size"`
	TotalPages int64 `json:"totalPages"`
	Total      int64 `json:"total"`
	HasNext    bool  `json:"hasNext"`
	HasPrev    bool  `json:"hasPrev"`
}
