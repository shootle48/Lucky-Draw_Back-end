package draw_condition

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"app/internal/logger"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func (s *Service) Create(ctx context.Context, req request.CreateDrawCondition) (*model.DrawCondition, bool, error) {
	prize := model.Prize{}
	err := s.db.NewSelect().
		Model(&prize).
		Where("id = ?", req.PrizeID).
		Where("quantity >= ?", req.Quantity).
		Scan(ctx)
	if err != nil {
		return nil, true, errors.New("not enough prize quantity")
	}

	m := &model.DrawCondition{
		RoomID:         req.RoomID,
		PrizeID:        req.PrizeID,
		FilterStatus:   req.FilterStatus,
		FilterPosition: req.FilterPosition,
		FilterIsActive: req.FilterIsActive,
		Quantity:       int64(req.Quantity),
	}

	_, err = s.db.NewInsert().Model(m).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("draw_condition already exists")
		}
		return nil, false, err
	}

	return m, false, nil
}

func (s *Service) Update(ctx context.Context, req request.UpdateDrawCondition, id request.GetByIDDrawCondition) (*model.DrawCondition, bool, error) {
	ex, err := s.db.NewSelect().Table("draw_conditions").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, false, err
	}

	m := &model.DrawCondition{
		ID:             id.ID,
		RoomID:         req.RoomID,
		PrizeID:        req.PrizeID,
		FilterStatus:   req.FilterStatus,
		FilterPosition: req.FilterPosition,
		FilterIsActive: req.FilterIsActive,
		Quantity:       int64(req.Quantity),
	}
	logger.Info(m)
	m.SetUpdateNow()
	_, err = s.db.NewUpdate().Model(m).
		Set("room_id = ?room_id, prize_id = ?prize_id, filter_status = ?filter_status, filter_position = ?filter_position, filter_is_active = ?filter_is_active, quantity = ?quantity").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("draw_conditions already exists")
		}
	}
	return m, false, err
}

func (s *Service) List(ctx context.Context, req request.ListDrawCondition) ([]response.ListDrawCondition, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []response.ListDrawCondition{}
	query := s.db.NewSelect().
		TableExpr("draw_conditions AS d").
		Column("d.id", "d.room_id", "d.prize_id", "d.filter_status", "d.filter_position", "d.filter_is_active", "d.quantity").
		Where("deleted_at IS NULL")

	if req.Search != "" {
		if req.SearchBy != "" {
			search := strings.ToLower(req.Search)
			query.Where(fmt.Sprintf("LOWER(d.%s) LIKE ?", req.SearchBy), "%"+search+"%")
		} else {
			query.Where("d.room_id::uuid = ?", req.Search)
		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("d.%s %s", req.SortBy, req.OrderBy)

	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}
	return m, count, err
}

func (s *Service) Get(ctx context.Context, id request.GetByIDDrawCondition) (*response.ListDrawCondition, error) {
	m := response.ListDrawCondition{}
	err := s.db.NewSelect().
		TableExpr("draw_conditions AS d").
		Column("d.id", "d.room_id", "d.prize_id", "d.filter_status", "d.filter_position", "d.filter_is_active", "d.quantity").
		Where("id = ?", id.ID).
		Where("deleted_at IS NULL").
		Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id request.GetByIDDrawCondition) error {
	ex, err := s.db.NewSelect().Table("draw_conditions").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
	if err != nil {
		return err
	}

	if !ex {
		return errors.New("draw_condition not found")
	}

	// data, err := s.db.NewDelete().Table("room").Where("id = ?", id.ID).Exec(ctx)
	_, err = s.db.NewDelete().Model((*model.DrawCondition)(nil)).Where("id = ?", id.ID).Exec(ctx)
	return err
}

// new function
func (s *Service) PreviewPlayer(ctx context.Context, req request.PreviewPlayers) ([]response.PreviewPlayer, error) {
	query := s.db.NewSelect().
		TableExpr("players AS p").
		Column("p.id", "p.prefix", "p.first_name", "p.last_name", "p.member_id", "p.position", "p.is_active", "p.status").
		Where("p.room_id = ?", req.RoomID).
		Where("p.deleted_at IS NULL")

	// Filter by Position
	if len(req.FilterPosition) > 0 {
		query = query.Where("p.position IN (?)", bun.In(req.FilterPosition))
	}

	// Filter by Status — ถ้าเป็น "not_received" หรือไม่มีค่าเลย ให้แสดงทุกคน
	if len(req.FilterStatus) > 0 {
		query = query.Where("p.status IN (?)", bun.In(req.FilterStatus))
	}

	// Filter by IsActive
	if req.FilterIsActive {
		query = query.Where("p.is_active = TRUE")
	}

	var players []response.PreviewPlayer
	err := query.Scan(ctx, &players)
	return players, err
}

func (s *Service) GetDrawConditionPreview(ctx context.Context, id string) (*response.DrawConditionPreview, error) {
	// Convert the string id to UUID
	drawConditionID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID format for id: %v", err)
	}

	// Fetch the draw condition
	dc := model.DrawCondition{}
	err = s.db.NewSelect().Model(&dc).Where("id = ?", drawConditionID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	// Fetch the associated players
	query := s.db.NewSelect().
		TableExpr("players AS p").
		Column("p.id", "p.prefix", "p.first_name", "p.last_name", "p.member_id", "p.position", "p.is_active", "p.status").
		Where("p.room_id = ?", dc.RoomID).
		Where("p.deleted_at IS NULL")

	// Apply filters
	if len(dc.FilterPosition) > 0 {
		query = query.Where("p.position IN (?)", bun.In(dc.FilterPosition))
	}

	if len(dc.FilterStatus) > 0 {
		query = query.Where("p.status IN (?)", bun.In(dc.FilterStatus))
	}

	if dc.FilterIsActive {
		query = query.Where("p.is_active = TRUE")
	}

	/// function check player_id & draw_condition_id in winners

	// subQuery := s.db.NewSelect().
	// 	Table("winners").
	// 	ColumnExpr("1").
	// 	Where("winners.player_id = p.id::text").
	// 	Where("winners.draw_condition_id = ?", dc.ID)

	// query = query.Where("NOT EXISTS (?)", subQuery)

	var players []response.PreviewPlayer
	err = query.Scan(ctx, &players)
	if err != nil {
		return nil, err
	}

	// Prepare the final response
	result := &response.DrawConditionPreview{
		ID:             dc.ID,
		RoomID:         dc.RoomID,
		PrizeID:        dc.PrizeID,
		FilterStatus:   dc.FilterStatus,
		FilterPosition: dc.FilterPosition,
		FilterIsActive: dc.FilterIsActive,
		Quantity:       dc.Quantity,
		Players:        players,
	}

	return result, nil
}
