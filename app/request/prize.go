package request

type CreatePrize struct {
	Name     string `json:"name" binding:"required"`
	ImageURL string `json:"image_url" binding:"required"`
	Quantity int64  `json:"quantity" binding:"required,gt=0"`
	RoomID   string `json:"room_id" binding:"required"`
}
type UpdatePrize struct {
	CreatePrize
}

type ListPrize struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
}

type GetByIDPrize struct {
	ID string `uri:"id" binding:"required"`
}
