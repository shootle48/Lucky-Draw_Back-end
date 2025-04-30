package winner

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"app/internal/logger"
	"context"
	"errors"
	"fmt"
	"strings"
)

func contains(strs []string, target string) bool {
	for _, s := range strs {
		if s == target {
			return true
		}
	}
	return false
}

func (s *Service) Create(ctx context.Context, req request.CreateWinner) (*response.ListWinnerDetail, bool, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, false, err
	}
	defer tx.Rollback()

	drawCondition := model.DrawCondition{}
	err = tx.NewSelect().
		Model(&drawCondition).
		Where("id = ?", req.DrawConditionID).
		Scan(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get draw condition: %w", err)
	}

	prize := model.Prize{}
	err = tx.NewSelect().
		Model(&prize).
		Where("id = ?", req.PrizeID).
		Where("quantity >= ?", drawCondition.Quantity).
		Scan(ctx)
	if err != nil {
		return nil, true, errors.New("not enough prize quantity")
	}

	m := &model.Winner{
		RoomID:          req.RoomID,
		PlayerID:        req.PlayerID,
		PrizeID:         req.PrizeID,
		DrawConditionID: req.DrawConditionID,
	}

	_, err = tx.NewInsert().Model(m).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("winner already exists")
		}
		return nil, false, err
	}

	if req.PlayerStatus == "not_received" {
		_, err = tx.NewUpdate().
			Model((*model.Player)(nil)).
			Set("status = ?", "received").
			Where("id = ?", req.PlayerID).
			Exec(ctx)
		if err != nil {
			return nil, false, fmt.Errorf("failed to update player status: %w", err)
		}

		_, err = tx.NewUpdate().
			Model((*model.Prize)(nil)).
			// Set("quantity = quantity - ?", drawCondition.Quantity).
			Set("quantity = quantity - 1").
			Where("id = ?", req.PrizeID).
			Where("quantity >= ?", drawCondition.Quantity).
			Exec(ctx)
		if err != nil {
			return nil, false, fmt.Errorf("failed to update prize quantity: %w", err)
		}
	} else if req.PlayerStatus == "no_show" || req.PlayerStatus == "waive" {
		_, err = tx.NewUpdate().
			Model((*model.Player)(nil)).
			Set("status = ?", req.PlayerStatus).
			Where("id = ?", req.PlayerID).
			Exec(ctx)
		if err != nil {
			return nil, false, fmt.Errorf("failed to update player status: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, false, err
	}

	result := &response.ListWinnerDetail{}
	err = s.db.NewSelect().
		TableExpr("winners AS w").
		ColumnExpr("w.id::uuid").
		ColumnExpr("w.room_id::uuid").
		ColumnExpr("r.name AS room_name").
		ColumnExpr("w.player_id::uuid").
		ColumnExpr("p.prefix").
		ColumnExpr("p.first_name").
		ColumnExpr("p.last_name").
		ColumnExpr("p.position").
		ColumnExpr("p.member_id").
		ColumnExpr("p.is_active").
		ColumnExpr("p.status").
		ColumnExpr("w.prize_id::uuid").
		ColumnExpr("pr.name AS prize_name").
		ColumnExpr("pr.image_url").
		ColumnExpr("w.draw_condition_id::uuid").
		ColumnExpr("dc.filter_status").
		ColumnExpr("dc.filter_position").
		ColumnExpr("dc.filter_is_active").
		ColumnExpr("dc.quantity").
		ColumnExpr("w.created_at").
		Join("JOIN rooms r ON r.id = w.room_id::uuid").
		Join("JOIN players p ON p.id = w.player_id::uuid").
		Join("JOIN prizes pr ON pr.id = w.prize_id::uuid").
		Join("JOIN draw_conditions dc ON dc.id = w.draw_condition_id::uuid").
		Where("w.id = ?", m.ID).
		Scan(ctx, result)

	if err != nil {
		return nil, false, err
	}

	return result, false, nil
}

func (s *Service) Update(ctx context.Context, req request.UpdateWinner, id request.GetByIDWinner) (*model.Winner, bool, error) {
	ex, err := s.db.NewSelect().Table("winners").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, false, err
	}

	m := &model.Winner{
		ID:              id.ID,
		RoomID:          req.RoomID,
		PlayerID:        req.PlayerID,
		PrizeID:         req.PrizeID,
		DrawConditionID: req.DrawConditionID,
	}
	logger.Info(m)
	m.SetUpdateNow()
	_, err = s.db.NewUpdate().Model(m).
		Set("room_id = ?room_id, player_id = ?player_id, prize_id = ?prize_id, draw_condition_id = ?draw_condition_id").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("winner already exists")
		}
	}
	return m, false, err
}

func (s *Service) List(ctx context.Context, req request.ListWinner) ([]response.ListWinner, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []response.ListWinner{}
	query := s.db.NewSelect().
		TableExpr("winners AS w").
		Column("w.id", "w.room_id", "w.player_id", "w.prize_id", "w.draw_condition_id").
		Where("w.deleted_at IS NULL")

	if req.Search != "" {
		if req.SearchBy != "" {
			search := strings.ToLower(req.Search)
			query.Where(fmt.Sprintf("LOWER(w.%s) LIKE ?", req.SearchBy), "%"+search+"%")
		} else {
			query.Where("w.room_id::uuid = ?", req.Search)
		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("w.%s %s", req.SortBy, req.OrderBy)

	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}
	return m, count, err
}

func (s *Service) Get(ctx context.Context, id request.GetByIDWinner) (*response.ListWinner, error) {
	m := response.ListWinner{}
	err := s.db.NewSelect().
		TableExpr("winners AS w").
		Column("w.id", "w.room_id", "w.player_id", "w.prize_id", "w.draw_condition_id").
		Where("w.id = ?", id.ID).
		Where("w.deleted_at IS NULL").
		Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id request.GetByIDWinner) error {
	ex, err := s.db.NewSelect().Table("winners").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
	if err != nil {
		return err
	}

	if !ex {
		return errors.New("winner not found")
	}

	// data, err := s.db.NewDelete().Table("room").Where("id = ?", id.ID).Exec(ctx)
	_, err = s.db.NewDelete().Model((*model.Winner)(nil)).Where("id = ?", id.ID).Exec(ctx)
	return err
}

// new function

func (s *Service) DashboardByRoomID(ctx context.Context, roomID string) (*response.WinnerDashboard, error) {
	var winners []response.ListWinnerDetail
	err := s.db.NewSelect().
		TableExpr(`
			(
				SELECT DISTINCT ON (w.player_id) 
                w.id::uuid,
                w.room_id::uuid,
                r.name AS room_name,
                w.player_id::uuid,
                p.prefix,
                p.first_name,
                p.last_name,
                p.position,
                p.member_id,
                p.is_active,
                p.status,
                w.prize_id::uuid,
                pr.name AS prize_name,
                pr.image_url,
                pr.quantity,
                w.draw_condition_id::uuid,
                dc.filter_status,
                dc.filter_position,
                dc.filter_is_active,
                dc.quantity,
                w.created_at
            FROM winners w
            JOIN rooms r ON r.id = w.room_id::uuid
            JOIN players p ON p.id = w.player_id::uuid
            JOIN prizes pr ON pr.id = w.prize_id::uuid
            JOIN draw_conditions dc ON dc.id = w.draw_condition_id::uuid
            WHERE w.room_id = ?
            ORDER BY w.player_id, w.created_at DESC
        ) AS latest_winners
		`, roomID).
		Scan(ctx, &winners)
	if err != nil {
		return nil, fmt.Errorf("failed to get winners: %w", err)
	}

	var prizes []response.PrizeDashboard
	err = s.db.NewSelect().
		TableExpr("prizes AS pr").
		ColumnExpr("pr.id::uuid").
		ColumnExpr("pr.room_id::uuid").
		ColumnExpr("pr.name").
		ColumnExpr("pr.image_url").
		ColumnExpr("pr.quantity").
		Where("pr.room_id = ?", roomID).
		Scan(ctx, &prizes)
	if err != nil {
		return nil, fmt.Errorf("failed to get prizes: %w", err)
	}

	return &response.WinnerDashboard{
		Winners: winners,
		Prizes:  prizes,
	}, nil
}
