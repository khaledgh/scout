package services

import (
	"context"
	"kashfi/internal/models"
	"time"

	"gorm.io/gorm"
)

type AnnouncementService struct {
	db    *gorm.DB
	notif *NotificationService
}

func NewAnnouncementService(db *gorm.DB, notif *NotificationService) *AnnouncementService {
	return &AnnouncementService{db: db, notif: notif}
}

func (s *AnnouncementService) List(ctx context.Context, userID uint, role string) ([]models.Announcement, error) {
	q := s.db.WithContext(ctx).Preload("Author").Order("pinned DESC, published_at DESC")
	if role == "member" {
		q = q.Where("audience IN ('all') AND published_at IS NOT NULL AND published_at <= ?", time.Now())
	}
	var announcements []models.Announcement
	err := q.Limit(50).Find(&announcements).Error
	return announcements, err
}

func (s *AnnouncementService) Get(ctx context.Context, id uint) (*models.Announcement, error) {
	var a models.Announcement
	if err := s.db.WithContext(ctx).Preload("Author").First(&a, id).Error; err != nil {
		return nil, ErrNotFound
	}
	return &a, nil
}

type UpdateAnnouncementInput struct {
	Title    *string
	Body     *string
	Audience *string
	UnitID   *uint
	Pinned   *bool
}

// Update modifies an announcement. The caller must be the author or super_admin.
func (s *AnnouncementService) Update(ctx context.Context, id, userID uint, role string, in UpdateAnnouncementInput) (*models.Announcement, error) {
	var a models.Announcement
	if err := s.db.WithContext(ctx).First(&a, id).Error; err != nil {
		return nil, ErrNotFound
	}
	if role != "super_admin" && a.AuthorID != userID {
		return nil, ErrForbidden
	}
	updates := map[string]interface{}{}
	if in.Title != nil { updates["title"] = *in.Title }
	if in.Body != nil { updates["body"] = *in.Body }
	if in.Audience != nil { updates["audience"] = *in.Audience }
	if in.UnitID != nil { updates["unit_id"] = *in.UnitID }
	if in.Pinned != nil { updates["pinned"] = *in.Pinned }

	if err := s.db.WithContext(ctx).Model(&a).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *AnnouncementService) Delete(ctx context.Context, id, userID uint, role string) error {
	var a models.Announcement
	if err := s.db.WithContext(ctx).First(&a, id).Error; err != nil {
		return ErrNotFound
	}
	if role != "super_admin" && a.AuthorID != userID {
		return ErrForbidden
	}
	return s.db.WithContext(ctx).Delete(&a).Error
}

type CreateAnnouncementInput struct {
	Title    string
	Body     string
	Audience string
	UnitID   *uint
	Pinned   bool
	AuthorID uint
}

func (s *AnnouncementService) Create(ctx context.Context, in CreateAnnouncementInput) (*models.Announcement, error) {
	now := time.Now()
	a := &models.Announcement{
		Title:       in.Title,
		Body:        in.Body,
		Audience:    in.Audience,
		UnitID:      in.UnitID,
		AuthorID:    in.AuthorID,
		Pinned:      in.Pinned,
		PublishedAt: &now,
	}
	if err := s.db.WithContext(ctx).Create(a).Error; err != nil {
		return nil, err
	}

	// Send notifications
	go func() {
		switch in.Audience {
		case "all":
			s.notif.NotifyAll(ctx, in.Title, in.Body, "announcement")
		case "unit":
			if in.UnitID != nil {
				s.notif.NotifyUnit(ctx, *in.UnitID, in.Title, in.Body, "announcement")
			}
		}
	}()

	return a, nil
}
