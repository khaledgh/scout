package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleSuperAdmin Role = "super_admin"
	RoleLeader     Role = "leader"
	RoleAssistant  Role = "assistant"
	RoleMember     Role = "member"
	RoleParent     Role = "parent"
)

type User struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	FullName    string         `gorm:"size:255;not null" json:"full_name"`
	Email       *string        `gorm:"size:191;uniqueIndex" json:"email,omitempty"`
	Phone       string         `gorm:"size:20;not null;uniqueIndex" json:"phone"`
	PasswordHash string        `gorm:"size:255;not null" json:"-"`
	Role        Role           `gorm:"type:enum('super_admin','leader','assistant','member','parent');not null;default:'member'" json:"role"`
	AvatarURL   *string        `gorm:"size:500" json:"avatar_url,omitempty"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	LastLoginAt *time.Time     `json:"last_login_at,omitempty"`

	// Relations
	UnitLeaderships []UnitLeader `gorm:"foreignKey:UserID" json:"-"`
}
