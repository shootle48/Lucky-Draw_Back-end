package room

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

func (s *Service) Create(ctx context.Context, req request.CreateRoom) (*model.Room, bool, error) {

	m := &model.Room{
		Name: req.Name,
	}

	_, err := s.db.NewInsert().Model(m).Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("room already exists")
		}
	}

	return m, false, err
}

func (s *Service) Update(ctx context.Context, req request.UpdateRoom, id request.GetByIDRoom) (*model.Room, bool, error) {
	ex, err := s.db.NewSelect().Table("rooms").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, false, err
	}

	m := &model.Room{
		ID:   id.ID,
		Name: req.Name,
	}
	logger.Info(m)
	m.SetUpdateNow()
	_, err = s.db.NewUpdate().Model(m).
		Set("name = ?name").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("room already exists")
		}
	}
	return m, false, err
}

func (s *Service) List(ctx context.Context, req request.ListRoom) ([]response.ListRoom, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []response.ListRoom{}
	query := s.db.NewSelect().
		TableExpr("rooms AS r").
		Column("r.id", "r.name").
		Where("deleted_at IS NULL")

	if req.Search != "" {
		search := fmt.Sprintf("%" + strings.ToLower(req.Search) + "%")
		if req.SearchBy != "" {
			search := strings.ToLower(req.Search)
			query.Where(fmt.Sprintf("LOWER(r.%s) LIKE ?", req.SearchBy), search)
		} else {
			query.Where("LOWER(name) LIKE ?", search)
		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("r.%s %s", req.SortBy, req.OrderBy)

	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}
	return m, count, err
}

func (s *Service) Get(ctx context.Context, id request.GetByIDRoom) (*response.ListRoom, error) {
	m := response.ListRoom{}
	err := s.db.NewSelect().
		TableExpr("rooms AS u").
		Column("u.id", "u.name").
		Where("id = ?", id.ID).Where("deleted_at IS NULL").Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id request.GetByIDRoom) error {
	ex, err := s.db.NewSelect().Table("rooms").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
	if err != nil {
		return err
	}

	if !ex {
		return errors.New("room not found")
	}

	// data, err := s.db.NewDelete().Table("room").Where("id = ?", id.ID).Exec(ctx)
	_, err = s.db.NewDelete().Model((*model.Room)(nil)).Where("id = ?", id.ID).Exec(ctx)
	return err
}

// new function
func (s *Service) ListAll(ctx context.Context, id request.GetByIDRoom) (*response.ListAllRoomResponse, error) {
	room := &model.Room{}
	err := s.db.NewSelect().
		Model(room).
		Relation("Players").
		Relation("Prizes").
		Relation("Prizes.Winners").
		Relation("Prizes.DrawConditions").
		Relation("DrawConditions").
		Relation("Winners").
		Where("room.id = ?", id.ID).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return &response.ListAllRoomResponse{
		Players:        room.Players,
		Prizes:         room.Prizes,
		DrawConditions: room.DrawConditions,
		Winners:        room.Winners,
	}, nil
}
