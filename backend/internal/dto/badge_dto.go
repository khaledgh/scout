package dto

type CreateBadgeRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
	XPReward    int    `json:"xp_reward" validate:"min=0"`
}

type AwardBadgeRequest struct {
	BadgeID uint `json:"badge_id" validate:"required"`
}

type AssessSkillRequest struct {
	SkillID uint `json:"skill_id" validate:"required"`
	Level   int  `json:"level" validate:"required,min=0,max=10"`
}

type CreateSkillRequest struct {
	Name        string `json:"name" validate:"required"`
	Category    string `json:"category"`
	Description string `json:"description"`
	MaxLevel    int    `json:"max_level" validate:"min=1,max=10"`
}

type UpdateBadgeRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	XPReward    *int    `json:"xp_reward"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateSkillRequest struct {
	Name        *string `json:"name"`
	Category    *string `json:"category"`
	Description *string `json:"description"`
	MaxLevel    *int    `json:"max_level"`
}
