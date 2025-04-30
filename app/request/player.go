package request

type CreatePlayer struct {
	Prefix    string `json:"prefix"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	MemberID  string `json:"member_id"`
	Position  string `json:"position"`
	RoomID    string `json:"room_id"`
	IsActive  bool   `json:"is_active"`
	Status    string `json:"status"`
}

type UpdatePlayer struct {
	CreatePlayer
}

type ListPlayer struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
}

type GetByIDPlayer struct {
	ID string `uri:"id" binding:"required"`
}

type CreatePlayerIM struct {
	Prefix    string `json:"prefix" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	MemberID  string `json:"member_id" binding:"required"`
	Position  string `json:"position" binding:"required"`
	RoomID    string `json:"room_id" binding:"required"`
	IsActive  bool   `json:"is_active"`
	Status    string `json:"status"`
}
