package services

import (
	"context"
	"kashfi/internal/config"
	"kashfi/internal/models"
	"kashfi/internal/utils"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type ActivityService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewActivityService(db *gorm.DB, cfg *config.Config) *ActivityService {
	return &ActivityService{db: db, cfg: cfg}
}

type ActivityFilter struct {
	Type   *string
	From   *time.Time
	To     *time.Time
	UnitID *uint
	Status *string
}

func (s *ActivityService) List(ctx context.Context, f ActivityFilter, page, pageSize int) ([]*models.Activity, int64, error) {
	q := s.db.WithContext(ctx).Model(&models.Activity{})
	if f.Type != nil { q = q.Where("type = ?", *f.Type) }
	if f.From != nil { q = q.Where("starts_at >= ?", *f.From) }
	if f.To != nil { q = q.Where("starts_at <= ?", *f.To) }
	if f.UnitID != nil { q = q.Where("unit_id = ? OR unit_id IS NULL", *f.UnitID) }
	if f.Status != nil { q = q.Where("status = ?", *f.Status) }

	var total int64
	q.Count(&total)
	var activities []*models.Activity
	err := q.Preload("ResponsibleUser").Preload("Unit").
		Order("starts_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&activities).Error
	return activities, total, err
}

func (s *ActivityService) Get(ctx context.Context, id uint) (*models.Activity, error) {
	var a models.Activity
	err := s.db.WithContext(ctx).
		Preload("ResponsibleUser").
		Preload("Unit").
		Preload("Media").
		First(&a, id).Error
	if err != nil {
		return nil, ErrNotFound
	}
	return &a, nil
}

type CreateActivityInput struct {
	Title             string
	Description       string
	Type              models.ActivityType
	Location          string
	LocationLat       *float64
	LocationLng       *float64
	StartsAt          time.Time
	EndsAt            time.Time
	ResponsibleUserID uint
	UnitID            *uint
}

func (s *ActivityService) Create(ctx context.Context, in CreateActivityInput) (*models.Activity, error) {
	a := &models.Activity{
		Title:             in.Title,
		Description:       in.Description,
		Type:              in.Type,
		Location:          in.Location,
		LocationLat:       in.LocationLat,
		LocationLng:       in.LocationLng,
		StartsAt:          in.StartsAt,
		EndsAt:            in.EndsAt,
		ResponsibleUserID: in.ResponsibleUserID,
		UnitID:            in.UnitID,
		Status:            models.ActivityStatusPlanned,
	}
	if err := s.db.WithContext(ctx).Create(a).Error; err != nil {
		return nil, err
	}
	return a, nil
}

type UpdateActivityInput struct {
	Title       *string
	Description *string
	Type        *models.ActivityType
	Location    *string
	LocationLat *float64
	LocationLng *float64
	StartsAt    *time.Time
	EndsAt      *time.Time
	Status      *models.ActivityStatus
	UnitID      *uint
}

func (s *ActivityService) Update(ctx context.Context, id uint, in UpdateActivityInput) (*models.Activity, error) {
	var a models.Activity
	if err := s.db.WithContext(ctx).First(&a, id).Error; err != nil {
		return nil, ErrNotFound
	}
	updates := map[string]interface{}{}
	if in.Title != nil { updates["title"] = *in.Title }
	if in.Description != nil { updates["description"] = *in.Description }
	if in.Type != nil { updates["type"] = *in.Type }
	if in.Location != nil { updates["location"] = *in.Location }
	if in.LocationLat != nil { updates["location_lat"] = *in.LocationLat }
	if in.LocationLng != nil { updates["location_lng"] = *in.LocationLng }
	if in.StartsAt != nil { updates["starts_at"] = *in.StartsAt }
	if in.EndsAt != nil { updates["ends_at"] = *in.EndsAt }
	if in.Status != nil { updates["status"] = *in.Status }
	if in.UnitID != nil { updates["unit_id"] = *in.UnitID }

	if err := s.db.WithContext(ctx).Model(&a).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *ActivityService) Delete(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&models.Activity{}, id)
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}

func (s *ActivityService) AddMedia(ctx context.Context, activityID, uploaderID uint, filename string, mediaType models.MediaType) (*models.ActivityMedia, error) {
	url := filepath.Join(s.cfg.Upload.PublicPath, filename)
	media := &models.ActivityMedia{
		ActivityID: activityID,
		URL:        url,
		MediaType:  mediaType,
		UploadedBy: uploaderID,
	}
	if err := s.db.WithContext(ctx).Create(media).Error; err != nil {
		return nil, err
	}
	return media, nil
}

func (s *ActivityService) GetFeedback(ctx context.Context, activityID uint) ([]models.ActivityFeedback, error) {
	var fb []models.ActivityFeedback
	err := s.db.WithContext(ctx).Where("activity_id = ?", activityID).
		Preload("Member").Find(&fb).Error
	return fb, err
}

type FeedbackSummary struct {
	AverageRating float64 `json:"average_rating"`
	Count         int     `json:"count"`
	Strengths     []string `json:"strengths"`
	Improvements  []string `json:"improvements"`
}

func (s *ActivityService) FeedbackSummary(ctx context.Context, activityID uint) (*FeedbackSummary, error) {
	var fb []models.ActivityFeedback
	s.db.WithContext(ctx).Where("activity_id = ?", activityID).Find(&fb)

	if len(fb) == 0 {
		return &FeedbackSummary{}, nil
	}
	total := 0
	strengths := []string{}
	improvements := []string{}
	for _, f := range fb {
		total += f.Rating
		if f.WhatWentWell != "" { strengths = append(strengths, f.WhatWentWell) }
		if f.WhatToImprove != "" { improvements = append(improvements, f.WhatToImprove) }
	}
	return &FeedbackSummary{
		AverageRating: float64(total) / float64(len(fb)),
		Count:         len(fb),
		Strengths:     strengths,
		Improvements:  improvements,
	}, nil
}

func (s *ActivityService) CreateFeedback(ctx context.Context, activityID, memberID uint, rating int, well, improve, comment string) (*models.ActivityFeedback, error) {
	fb := &models.ActivityFeedback{
		ActivityID:    activityID,
		MemberID:      memberID,
		Rating:        rating,
		WhatWentWell:  well,
		WhatToImprove: improve,
		Comment:       comment,
	}
	if err := s.db.WithContext(ctx).Create(fb).Error; err != nil {
		return nil, err
	}
	return fb, nil
}

func (s *ActivityService) UploadPath() string { return s.cfg.Upload.Dir }
func (s *ActivityService) UploadMaxMB() int64  { return s.cfg.Upload.MaxSizeMB }
func (s *ActivityService) UploadTypes() string  { return s.cfg.Upload.AllowedTypes }

var _ = utils.SaveUpload // referenced via handler
