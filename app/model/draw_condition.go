package model

import (
	"github.com/uptrace/bun"
)

type DrawCondition struct {
	bun.BaseModel `bun:"table:draw_conditions"`

	ID             string `json:"id" bun:",pk,type:uuid,default:gen_random_uuid()"`
	RoomID         string `bun:"room_id,notnull"`
	PrizeID        string `bun:"prize_id,notnull"`
	FilterStatus   string `bun:"filter_status,notnull"`
	FilterPosition string `bun:"filter_position,notnull"`
	Quantity       int64  `bun:"quantity,notnull"`

	Room  *Room  `bun:"rel:belongs-to,join:room_id=id"`
	Prize *Prize `bun:"rel:belongs-to,join:prize_id=id"`

	Winners []Winner `bun:"rel:has-many,join:id=draw_condition_id"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
