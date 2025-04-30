package response

type ListDrawCondition struct {
	ID             string `bun:"id" json:"id"`
	RoomID         string `bun:"room_id" json:"room_id"`
	PrizeID        string `bun:"prize_id" json:"prize_id"`
	FilterStatus   string `bun:"filter_status" json:"filter_status"`
	FilterPosition string `bun:"filter_position" json:"filter_position"`
	FilterIsActive bool   `bun:"filter_is_active" json:"filter_is_active"`
	Quantity       int64  `bun:"quantity" json:"quantity"`
}

type PreviewPlayer struct {
	ID        string `bun:"id" json:"id"`
	Prefix    string `bun:"prefix" json:"prefix"`
	FirstName string `bun:"first_name" json:"first_name"`
	LastName  string `bun:"last_name" json:"last_name"`
	MemberID  string `bun:"member_id" json:"member_id"`
	Position  string `bun:"position" json:"position"`
	IsActive  bool   `bun:"is_active" json:"is_active"`
	Status    string `bun:"status" json:"status"`
}

type DrawConditionPreview struct {
	ID             string          `bun:"id" json:"id"`
	RoomID         string          `bun:"room_id" json:"room_id"`
	PrizeID        string          `bun:"prize_id" json:"prize_id"`
	FilterStatus   []string        `bun:"filter_status" json:"filter_status"`
	FilterPosition []string        `bun:"filter_position" json:"filter_position"`
	FilterIsActive bool            `bun:"filter_is_active" json:"filter_is_active"`
	Quantity       int64           `bun:"quantity" json:"quantity"`
	Players        []PreviewPlayer `bun:"players" json:"players"`
}
