package services

import (
	"context"
	"kashfi/internal/models"
	"time"

	"gorm.io/gorm"
)

type BadgeService struct {
	db     *gorm.DB
	gamify *GamificationService
}

func NewBadgeService(db *gorm.DB, gamify *GamificationService) *BadgeService {
	return &BadgeService{db: db, gamify: gamify}
}

func (s *BadgeService) ListCatalog(ctx context.Context) ([]models.Badge, error) {
	var badges []models.Badge
	err := s.db.WithContext(ctx).Where("is_active = true").Find(&badges).Error
	return badges, err
}

type CreateBadgeInput struct {
	Name        string
	Description string
	Category    string
	XPReward    int
}

func (s *BadgeService) Create(ctx context.Context, in CreateBadgeInput) (*models.Badge, error) {
	b := &models.Badge{
		Name:        in.Name,
		Description: in.Description,
		Category:    in.Category,
		XPReward:    in.XPReward,
		IsActive:    true,
	}
	if err := s.db.WithContext(ctx).Create(b).Error; err != nil {
		return nil, err
	}
	return b, nil
}

func (s *BadgeService) Get(ctx context.Context, id uint) (*models.Badge, error) {
	var b models.Badge
	if err := s.db.WithContext(ctx).First(&b, id).Error; err != nil {
		return nil, ErrNotFound
	}
	return &b, nil
}

type UpdateBadgeInput struct {
	Name        *string
	Description *string
	Category    *string
	XPReward    *int
	IsActive    *bool
}

func (s *BadgeService) Update(ctx context.Context, id uint, in UpdateBadgeInput) (*models.Badge, error) {
	var b models.Badge
	if err := s.db.WithContext(ctx).First(&b, id).Error; err != nil {
		return nil, ErrNotFound
	}
	updates := map[string]interface{}{}
	if in.Name != nil { updates["name"] = *in.Name }
	if in.Description != nil { updates["description"] = *in.Description }
	if in.Category != nil { updates["category"] = *in.Category }
	if in.XPReward != nil { updates["xp_reward"] = *in.XPReward }
	if in.IsActive != nil { updates["is_active"] = *in.IsActive }

	if err := s.db.WithContext(ctx).Model(&b).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (s *BadgeService) Delete(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&models.Badge{}, id)
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}

// RevokeBadge removes a badge award from a member.
func (s *BadgeService) RevokeBadge(ctx context.Context, memberID, badgeID uint) error {
	result := s.db.WithContext(ctx).
		Where("member_id = ? AND badge_id = ?", memberID, badgeID).
		Delete(&models.MemberBadge{})
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}

func (s *BadgeService) AwardBadge(ctx context.Context, memberID, badgeID, awardedByID uint) (*models.MemberBadge, error) {
	var badge models.Badge
	if err := s.db.WithContext(ctx).First(&badge, badgeID).Error; err != nil {
		return nil, ErrNotFound
	}

	var existing models.MemberBadge
	if err := s.db.WithContext(ctx).Where("member_id = ? AND badge_id = ?", memberID, badgeID).
		First(&existing).Error; err == nil {
		return nil, ErrConflict
	}

	mb := &models.MemberBadge{
		MemberID:  memberID,
		BadgeID:   badgeID,
		AwardedAt: time.Now(),
		AwardedBy: &awardedByID,
		Progress:  100,
	}
	if err := s.db.WithContext(ctx).Create(mb).Error; err != nil {
		return nil, err
	}

	if badge.XPReward > 0 {
		refID := badgeID
		s.gamify.AwardXP(ctx, memberID, nil, models.XPSourceBadge, badge.XPReward, &refID, "badge: "+badge.Name)
	}

	return mb, nil
}

func (s *BadgeService) MemberBadges(ctx context.Context, memberID uint) ([]models.MemberBadge, error) {
	var badges []models.MemberBadge
	err := s.db.WithContext(ctx).Where("member_id = ?", memberID).
		Preload("Badge").Find(&badges).Error
	return badges, err
}

func (s *BadgeService) ListSkills(ctx context.Context) ([]models.Skill, error) {
	var skills []models.Skill
	err := s.db.WithContext(ctx).Find(&skills).Error
	return skills, err
}

func (s *BadgeService) CreateSkill(ctx context.Context, name, category, description string, maxLevel int) (*models.Skill, error) {
	sk := &models.Skill{Name: name, Category: category, Description: description, MaxLevel: maxLevel}
	if err := s.db.WithContext(ctx).Create(sk).Error; err != nil {
		return nil, err
	}
	return sk, nil
}

func (s *BadgeService) GetSkill(ctx context.Context, id uint) (*models.Skill, error) {
	var sk models.Skill
	if err := s.db.WithContext(ctx).First(&sk, id).Error; err != nil {
		return nil, ErrNotFound
	}
	return &sk, nil
}

type UpdateSkillInput struct {
	Name        *string
	Category    *string
	Description *string
	MaxLevel    *int
}

func (s *BadgeService) UpdateSkill(ctx context.Context, id uint, in UpdateSkillInput) (*models.Skill, error) {
	var sk models.Skill
	if err := s.db.WithContext(ctx).First(&sk, id).Error; err != nil {
		return nil, ErrNotFound
	}
	updates := map[string]interface{}{}
	if in.Name != nil { updates["name"] = *in.Name }
	if in.Category != nil { updates["category"] = *in.Category }
	if in.Description != nil { updates["description"] = *in.Description }
	if in.MaxLevel != nil { updates["max_level"] = *in.MaxLevel }

	if err := s.db.WithContext(ctx).Model(&sk).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &sk, nil
}

func (s *BadgeService) DeleteSkill(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&models.Skill{}, id)
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}

func (s *BadgeService) AssessSkill(ctx context.Context, memberID, skillID, assessorID uint, level int) (*models.MemberSkill, error) {
	now := time.Now()
	ms := models.MemberSkill{
		MemberID:   memberID,
		SkillID:    skillID,
		Level:      level,
		AssessedBy: &assessorID,
		AssessedAt: &now,
	}
	err := s.db.WithContext(ctx).
		Where("member_id = ? AND skill_id = ?", memberID, skillID).
		Assign(ms).FirstOrCreate(&ms).Error
	if err != nil {
		return nil, err
	}
	return &ms, nil
}
