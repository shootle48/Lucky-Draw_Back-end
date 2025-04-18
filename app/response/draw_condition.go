package response

type ListDrawCondition struct {
	ID             string `bun:"id" json:"id"`
	RoomID         string `bun:"room_id" json:"room_id"`
	PrizeID        string `bun:"prize_id" json:"prize_id"`
	FilterStatus   string `bun:"filter_status" json:"filter_status"`
	FilterPosition string `bun:"filter_position" json:"filter_position"`
	Quantity       int64  `bun:"quantity" json:"quantity"`
}

type PreviewPlayer struct {
	ID        string `bun:"id" json:"id"`
	Prefix    string `bun:"prefix" json:"prefix"`
	FirstName string `bun:"first_name" json:"first_name"`
	LastName  string `bun:"last_name" json:"last_name"`
	Position  string `bun:"position" json:"position"`
}
