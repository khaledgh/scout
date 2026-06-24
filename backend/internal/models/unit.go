package models

import (
	"time"

	"gorm.io/gorm"
)

type Unit struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Name       string         `gorm:"size:255;not null" json:"name"`
	Section    Section        `gorm:"type:enum('ashbal','kashaf','jawala','mukashe');not null" json:"section"`
	Motto      string         `gorm:"size:500" json:"motto"`
	EmblemURL  *string        `gorm:"size:500" json:"emblem_url,omitempty"`
	Level      int            `gorm:"default:1" json:"level"`
	ScoreTotal int            `gorm:"default:0" json:"score_total"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`

	Leaders []UnitLeader `gorm:"foreignKey:UnitID" json:"leaders,omitempty"`
	Members []UnitMember `gorm:"foreignKey:UnitID" json:"members,omitempty"`
}

type UnitLeaderRole string

const (
	UnitLeaderRoleLeader    UnitLeaderRole = "leader"
	UnitLeaderRoleAssistant UnitLeaderRole = "assistant"
)

type UnitLeader struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UnitID     uint           `gorm:"not null;index" json:"unit_id"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	RoleInUnit UnitLeaderRole `gorm:"type:enum('leader','assistant');default:'assistant'" json:"role_in_unit"`
	Unit       Unit           `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	User       User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type UnitMember struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UnitID    uint           `gorm:"not null;index" json:"unit_id"`
	MemberID  uint           `gorm:"not null;index" json:"member_id"`
	IsPrimary bool           `gorm:"default:true" json:"is_primary"`
	JoinedAt  time.Time      `json:"joined_at"`
	LeftAt    *time.Time     `json:"left_at,omitempty"`
	Unit      Unit           `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	Member    Member         `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
