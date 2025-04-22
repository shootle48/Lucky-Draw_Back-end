package model

import (
	"github.com/uptrace/bun"
)

type Player struct {
	bun.BaseModel `bun:"table:players"`

	ID        string `json:"id" bun:",pk,type:uuid,default:gen_random_uuid()"`
	Prefix    string `bun:"prefix,notnull"`
	FirstName string `bun:"first_name,notnull"`
	LastName  string `bun:"last_name,notnull"`
	MemberID  string `bun:"member_id,unique:member_room,notnull"`
	Position  string `bun:"position,notnull"`
	RoomID    string `bun:"room_id,unique:member_room,notnull"`
	IsActive  bool   `bun:"is_active,type:boolean,default:false,notnull"`

	Room *Room `bun:"rel:belongs-to,join:room_id=id"`

	Winners []Winner `bun:"rel:has-many,join:id=player_id"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
