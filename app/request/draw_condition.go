package request

type CreateDrawCondition struct {
	RoomID         string `json:"room_id"`
	PrizeID        string `json:"prize_id"`
	FilterStatus   string `json:"filter_status"`
	FilterPosition string `json:"filter_position"`
	Quantity       int64  `json:"quantity"`
}

type UpdateDrawCondition struct {
	CreateDrawCondition
}

type ListDrawCondition struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
}

type GetByIDDrawCondition struct {
	ID string `uri:"id" binding:"required"`
}

type PreviewPlayers struct {
	RoomID         string   `json:"room_id" binding:"required"`
	FilterStatus   string   `json:"filter_status"`
	FilterPosition []string `json:"filter_position"`
}
