package response

import "app/app/model"

type ListPlayer struct {
	ID        string `bun:"id" json:"id"`
	Prefix    string `bun:"prefix" json:"prefix"`
	FirstName string `bun:"first_name" json:"first_name"`
	LastName  string `bun:"last_name" json:"last_name"`
	MemberID  string `bun:"member_id" json:"member_id"`
	Position  string `bun:"position" json:"position"`
	RoomID    string `bun:"room_id" json:"room_id"`
	RoomName  string `bun:"room_name" json:"room_name"`
	IsActive  bool   `bun:"is_active" json:"is_active"`
}

type ListAllRoomResponse struct {
	Players        []model.Player        `bun:"players" json:"players"`
	Prizes         []model.Prize         `bun:"prizes" json:"prizes"`
	DrawConditions []model.DrawCondition `bun:"draw_conditions" json:"draw_conditions"`
	Winners        []model.Winner        `bun:"winners" json:"winners"`
}
