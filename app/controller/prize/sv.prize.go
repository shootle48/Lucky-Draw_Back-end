package prize

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"app/internal/logger"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

func (s *Service) Create(ctx context.Context, req request.CreatePrize) (*response.Prize, bool, error) {

	roomExists, err := s.db.NewSelect().Model((*model.Room)(nil)).Where("id = ?", req.RoomID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}
	if !roomExists {
		return nil, true, errors.New("room not found")
	}

	m := &model.Prize{
		Name:     req.Name,
		ImageURL: req.ImageURL,
		Quantity: req.Quantity,
		RoomID:   req.RoomID,
	}

	_, err = s.db.NewInsert().Model(m).Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("prize already exists")
		}
		return nil, false, err
	}

	resp := &response.Prize{
		ID:        m.ID,
		Name:      m.Name,
		ImageURL:  m.ImageURL,
		Quantity:  m.Quantity,
		RoomID:    m.RoomID,
		CreatedAt: time.Unix(m.CreatedAt, 0).Format("2006-01-02 15:04:05"),
	}

	return resp, false, nil
}

func (s *Service) Update(ctx context.Context, req request.UpdatePrize, id request.GetByIDPrize) (*model.Prize, bool, error) {
	ex, err := s.db.NewSelect().Table("prizes").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, false, err
	}

	m := &model.Prize{
		ID:       id.ID,
		Name:     req.Name,
		ImageURL: req.ImageURL,
		Quantity: req.Quantity,
		RoomID:   req.RoomID,
	}
	logger.Info(m)
	m.SetUpdateNow()
	_, err = s.db.NewUpdate().Model(m).
		Set("name = ?name, image_url = ?image_url, quantity = ?quantity, room_id = ?room_id").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("prize already exists")
		}
	}
	return m, false, err
}

func (s *Service) List(ctx context.Context, req request.ListPrize) ([]response.ListPrize, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []response.ListPrize{}
	query := s.db.NewSelect().
		TableExpr("prizes AS p").
		Column("p.id", "p.name", "p.image_url", "p.quantity", "p.room_id").
		Where("deleted_at IS NULL")

	if req.Search != "" {
		if req.SearchBy != "" {
			search := strings.ToLower(req.Search)
			query.Where(fmt.Sprintf("LOWER(p.%s) LIKE ?", req.SearchBy), "%"+search+"%")
		} else {
			query.Where("p.room_id = ?", req.Search)
		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("p.%s %s", req.SortBy, req.OrderBy)

	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}
	return m, count, err
}

func (s *Service) Get(ctx context.Context, id request.GetByIDPrize) (*response.ListPrize, error) {
	m := response.ListPrize{}
	err := s.db.NewSelect().
		TableExpr("prizes AS p").
		Column("p.id", "p.name", "p.image_url", "p.quantity", "p.room_id").
		Where("id = ?", id.ID).
		Where("deleted_at IS NULL").
		Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id request.GetByIDPrize) error {
	ex, err := s.db.NewSelect().Table("prizes").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
	if err != nil {
		return err
	}

	if !ex {
		return errors.New("prize not found")
	}

	// data, err := s.db.NewDelete().Table("room").Where("id = ?", id.ID).Exec(ctx)
	_, err = s.db.NewDelete().Model((*model.Prize)(nil)).Where("id = ?", id.ID).Exec(ctx)
	return err
}
