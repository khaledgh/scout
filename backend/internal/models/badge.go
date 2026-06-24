package models

import (
	"time"

	"gorm.io/gorm"
)

type Badge struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Name         string         `gorm:"size:255;not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	Category     string         `gorm:"size:100" json:"category"`
	IconURL      *string        `gorm:"size:500" json:"icon_url,omitempty"`
	XPReward     int            `gorm:"default:0" json:"xp_reward"`
	CriteriaJSON *string        `gorm:"type:text" json:"criteria_json,omitempty"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
}

type MemberBadge struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	MemberID  uint       `gorm:"not null;index" json:"member_id"`
	BadgeID   uint       `gorm:"not null;index" json:"badge_id"`
	AwardedAt time.Time  `json:"awarded_at"`
	AwardedBy *uint      `gorm:"index" json:"awarded_by,omitempty"`
	Progress  int        `gorm:"default:0" json:"progress"`
	Badge     Badge      `gorm:"foreignKey:BadgeID" json:"badge,omitempty"`
	Member    Member     `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	Awarder   *User      `gorm:"foreignKey:AwardedBy" json:"awarder,omitempty"`
}

type Skill struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Category    string         `gorm:"size:100" json:"category"`
	Description string         `gorm:"type:text" json:"description"`
	MaxLevel    int            `gorm:"default:5" json:"max_level"`
}

type MemberSkill struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	MemberID   uint       `gorm:"not null;index" json:"member_id"`
	SkillID    uint       `gorm:"not null;index" json:"skill_id"`
	Level      int        `gorm:"default:0" json:"level"`
	AssessedBy *uint      `gorm:"index" json:"assessed_by,omitempty"`
	AssessedAt *time.Time `json:"assessed_at,omitempty"`
	Skill      Skill      `gorm:"foreignKey:SkillID" json:"skill,omitempty"`
	Assessor   *User      `gorm:"foreignKey:AssessedBy" json:"assessor,omitempty"`
}

type Evaluation struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	MemberID      uint      `gorm:"not null;index" json:"member_id"`
	EvaluatorID   uint      `gorm:"not null;index" json:"evaluator_id"`
	Period        string    `gorm:"size:20;not null" json:"period"`
	Discipline    int       `json:"discipline"`
	Participation int       `json:"participation"`
	Leadership    int       `json:"leadership"`
	Skill         int       `json:"skill"`
	Overall       int       `json:"overall"`
	Notes         string    `gorm:"type:text" json:"notes"`
	Member        Member    `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	Evaluator     User      `gorm:"foreignKey:EvaluatorID" json:"evaluator,omitempty"`
}
