package services

import (
	"context"
	"kashfi/internal/models"

	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

func (s *NotificationService) Create(ctx context.Context, userID uint, title, body, notifType, dataJSON string) error {
	n := &models.Notification{
		UserID:   userID,
		Title:    title,
		Body:     body,
		Type:     notifType,
		DataJSON: dataJSON,
	}
	return s.db.WithContext(ctx).Create(n).Error
}

func (s *NotificationService) ForUser(ctx context.Context, userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	err := s.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").Limit(50).Find(&notifications).Error
	return notifications, err
}

func (s *NotificationService) MarkRead(ctx context.Context, notifID, userID uint) error {
	result := s.db.WithContext(ctx).Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notifID, userID).
		Update("read_at", gorm.Expr("NOW()"))
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}

func (s *NotificationService) NotifyAll(ctx context.Context, title, body, notifType string) error {
	var userIDs []uint
	s.db.WithContext(ctx).Model(&models.User{}).Where("is_active = true").Pluck("id", &userIDs)
	for _, uid := range userIDs {
		s.Create(ctx, uid, title, body, notifType, "{}")
	}
	return nil
}

func (s *NotificationService) NotifyUnit(ctx context.Context, unitID uint, title, body, notifType string) error {
	var userIDs []uint
	s.db.WithContext(ctx).Model(&models.UnitMember{}).
		Joins("JOIN members m ON m.id = unit_members.member_id").
		Where("unit_members.unit_id = ? AND m.user_id IS NOT NULL", unitID).
		Pluck("m.user_id", &userIDs)
	for _, uid := range userIDs {
		s.Create(ctx, uid, title, body, notifType, "{}")
	}
	return nil
}
