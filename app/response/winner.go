package response

type ListWinner struct {
	ID              string `bun:"id" json:"id"`
	RoomID          string `bun:"room_id" json:"room_id"`
	PlayerID        string `bun:"player_id" json:"player_id"`
	PrizeID         string `bun:"prize_id" json:"prize_id"`
	DrawConditionID string `bun:"draw_condition_id" json:"draw_condition_id"`
}

type ListWinnerDetail struct {
	ID       string `bun:"id" json:"id"`
	RoomID   string `bun:"room_id" json:"room_id"`
	RoomName string `bun:"room_name" json:"room_name"`

	PlayerID  string `bun:"player_id" json:"player_id"`
	Prefix    string `bun:"prefix" json:"prefix"`
	FirstName string `bun:"first_name" json:"first_name"`
	LastName  string `bun:"last_name" json:"last_name"`
	MemberID  string `bun:"member_id" json:"member_id"`
	Position  string `bun:"position" json:"position"`
	IsActive  bool   `bun:"is_active" json:"is_active"`
	Status    string `bun:"status" json:"status"`

	PrizeID   string `bun:"prize_id" json:"prize_id"`
	PrizeName string `bun:"prize_name" json:"prize_name"`
	ImageURL  string `bun:"image_url" json:"image_url"`

	DrawConditionID string `bun:"draw_condition_id" json:"draw_condition_id"`
	FilterStatus    string `bun:"filter_status" json:"filter_status"`
	FilterPosition  string `bun:"filter_position" json:"filter_position"`
	FilterIsActive  bool   `bun:"filter_is_active" json:"filter_is_active"`
	Quantity        int64  `bun:"quantity" json:"quantity"`
	CreatedAt       string `bun:"created_at" json:"created_at"`
}

type WinnerDashboard struct {
	Winners []ListWinnerDetail `json:"winners"`
	Prizes  []PrizeDashboard   `json:"prizes"`
}

type PrizeDashboard struct {
	ID       string `json:"id"`
	RoomID   string `json:"room_id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	Quantity int    `json:"quantity"`
}
