package services

import (
	"context"
	"kashfi/internal/models"
	"time"

	"gorm.io/gorm"
)

type UnitService struct {
	db *gorm.DB
}

func NewUnitService(db *gorm.DB) *UnitService {
	return &UnitService{db: db}
}

func (s *UnitService) List(ctx context.Context) ([]*models.Unit, error) {
	var units []*models.Unit
	err := s.db.WithContext(ctx).Where("is_active = true").
		Preload("Leaders.User").Preload("Members").Find(&units).Error
	return units, err
}

func (s *UnitService) Get(ctx context.Context, id uint) (*models.Unit, error) {
	var unit models.Unit
	err := s.db.WithContext(ctx).
		Preload("Leaders.User").
		Preload("Members.Member").
		First(&unit, id).Error
	if err != nil {
		return nil, ErrNotFound
	}
	return &unit, nil
}

type CreateUnitInput struct {
	Name    string
	Section models.Section
	Motto   string
}

func (s *UnitService) Create(ctx context.Context, in CreateUnitInput) (*models.Unit, error) {
	unit := &models.Unit{
		Name:     in.Name,
		Section:  in.Section,
		Motto:    in.Motto,
		IsActive: true,
		Level:    1,
	}
	if err := s.db.WithContext(ctx).Create(unit).Error; err != nil {
		return nil, err
	}
	return unit, nil
}

type UpdateUnitInput struct {
	Name     *string
	Section  *models.Section
	Motto    *string
	IsActive *bool
}

func (s *UnitService) Update(ctx context.Context, id uint, in UpdateUnitInput) (*models.Unit, error) {
	var unit models.Unit
	if err := s.db.WithContext(ctx).First(&unit, id).Error; err != nil {
		return nil, ErrNotFound
	}
	updates := map[string]interface{}{}
	if in.Name != nil { updates["name"] = *in.Name }
	if in.Section != nil { updates["section"] = *in.Section }
	if in.Motto != nil { updates["motto"] = *in.Motto }
	if in.IsActive != nil { updates["is_active"] = *in.IsActive }

	if err := s.db.WithContext(ctx).Model(&unit).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &unit, nil
}

func (s *UnitService) Delete(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&models.Unit{}, id)
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}

func (s *UnitService) AddMembers(ctx context.Context, unitID uint, memberIDs []uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, mid := range memberIDs {
			um := models.UnitMember{
				UnitID:    unitID,
				MemberID:  mid,
				IsPrimary: true,
				JoinedAt:  time.Now(),
			}
			if err := tx.Where("unit_id = ? AND member_id = ?", unitID, mid).
				FirstOrCreate(&um).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *UnitService) RemoveMember(ctx context.Context, unitID, memberID uint) error {
	result := s.db.WithContext(ctx).
		Where("unit_id = ? AND member_id = ?", unitID, memberID).
		Delete(&models.UnitMember{})
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}

func (s *UnitService) AssignLeader(ctx context.Context, unitID, userID uint, role models.UnitLeaderRole) error {
	ul := models.UnitLeader{
		UnitID:     unitID,
		UserID:     userID,
		RoleInUnit: role,
	}
	return s.db.WithContext(ctx).
		Where("unit_id = ? AND user_id = ?", unitID, userID).
		Assign(ul).FirstOrCreate(&ul).Error
}

func (s *UnitService) Leaderboard(ctx context.Context) ([]*models.Unit, error) {
	var units []*models.Unit
	err := s.db.WithContext(ctx).Where("is_active = true").
		Order("score_total DESC").Find(&units).Error
	return units, err
}

func (s *UnitService) IsUserLeaderOf(ctx context.Context, userID, unitID uint) bool {
	var count int64
	s.db.WithContext(ctx).Model(&models.UnitLeader{}).
		Where("user_id = ? AND unit_id = ?", userID, unitID).Count(&count)
	return count > 0
}
