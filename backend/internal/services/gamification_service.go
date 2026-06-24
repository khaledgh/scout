package services

import (
	"context"
	"kashfi/internal/config"
	"kashfi/internal/models"
	"math"

	"gorm.io/gorm"
)

type GamificationService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewGamificationService(db *gorm.DB, cfg *config.Config) *GamificationService {
	return &GamificationService{db: db, cfg: cfg}
}

// AwardXP adds XP to a member and propagates to their primary unit.
func (s *GamificationService) AwardXP(ctx context.Context, memberID uint, unitID *uint, source models.XPSource, points int, refID *uint, note string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		event := models.XPEvent{
			MemberID: memberID,
			UnitID:   unitID,
			Source:   source,
			Points:   points,
			RefID:    refID,
			Note:     note,
		}
		if err := tx.Create(&event).Error; err != nil {
			return err
		}

		// Update member XP + level
		var member models.Member
		if err := tx.First(&member, memberID).Error; err != nil {
			return err
		}
		newXP := member.XPTotal + points
		newLevel := s.computeLevel(newXP)
		if err := tx.Model(&member).Updates(map[string]interface{}{
			"xp_total": newXP,
			"level":    newLevel,
		}).Error; err != nil {
			return err
		}

		// Propagate to unit score
		if unitID != nil {
			if err := tx.Model(&models.Unit{}).Where("id = ?", *unitID).
				UpdateColumn("score_total", gorm.Expr("score_total + ?", points)).Error; err != nil {
				return err
			}
		} else {
			// Find primary unit
			var um models.UnitMember
			if err := tx.Where("member_id = ? AND is_primary = true", memberID).First(&um).Error; err == nil {
				tx.Model(&models.Unit{}).Where("id = ?", um.UnitID).
					UpdateColumn("score_total", gorm.Expr("score_total + ?", points))
			}
		}
		return nil
	})
}

func (s *GamificationService) computeLevel(xp int) int {
	base := s.cfg.Gamify.LevelBaseXP
	if base <= 0 {
		base = 100
	}
	// Level = 1 + floor(sqrt(xp / base))
	level := 1 + int(math.Floor(math.Sqrt(float64(xp)/float64(base))))
	if level < 1 {
		level = 1
	}
	return level
}

type LeaderboardEntry struct {
	MemberID uint   `json:"member_id"`
	FullName string  `json:"full_name"`
	Section  string  `json:"section"`
	XPTotal  int     `json:"xp_total"`
	Level    int     `json:"level"`
	PhotoURL *string `json:"photo_url,omitempty"`
}

func (s *GamificationService) MemberLeaderboard(ctx context.Context, section *string, limit int) ([]LeaderboardEntry, error) {
	q := s.db.WithContext(ctx).Model(&models.Member{}).
		Select("id AS member_id, full_name, section, xp_total, level, photo_url").
		Where("status = 'active'")
	if section != nil {
		q = q.Where("section = ?", *section)
	}
	var result []LeaderboardEntry
	err := q.Order("xp_total DESC").Limit(limit).Scan(&result).Error
	return result, err
}

type UnitLeaderboardEntry struct {
	UnitID     uint   `json:"unit_id"`
	Name       string `json:"name"`
	Section    string `json:"section"`
	ScoreTotal int    `json:"score_total"`
	Level      int    `json:"level"`
	EmblemURL  *string `json:"emblem_url,omitempty"`
}

func (s *GamificationService) UnitLeaderboard(ctx context.Context) ([]UnitLeaderboardEntry, error) {
	var result []UnitLeaderboardEntry
	err := s.db.WithContext(ctx).Model(&models.Unit{}).
		Select("id AS unit_id, name, section, score_total, level, emblem_url").
		Where("is_active = true").
		Order("score_total DESC").Scan(&result).Error
	return result, err
}

func (s *GamificationService) GetXPHistory(ctx context.Context, memberID uint) ([]models.XPEvent, error) {
	var events []models.XPEvent
	err := s.db.WithContext(ctx).Where("member_id = ?", memberID).
		Order("created_at DESC").Limit(100).Find(&events).Error
	return events, err
}
