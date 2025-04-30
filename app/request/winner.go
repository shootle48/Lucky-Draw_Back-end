package request

type CreateWinner struct {
	RoomID          string `json:"room_id"`
	PlayerID        string `json:"player_id"`
	PrizeID         string `json:"prize_id"`
	DrawConditionID string `json:"draw_condition_id"`
	PlayerStatus    string `json:"player_status"`
}

type UpdateWinner struct {
	CreateWinner
}

type ListWinner struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
}

type GetByIDWinner struct {
	ID string `uri:"id" binding:"required"`
}
