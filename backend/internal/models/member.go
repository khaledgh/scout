package models

import (
	"time"

	"gorm.io/gorm"
)

type Section string

const (
	SectionAshbal  Section = "ashbal"
	SectionKashaf  Section = "kashaf"
	SectionJawala  Section = "jawala"
	SectionMukashe Section = "mukashe"
)

type MemberStatus string

const (
	MemberStatusActive   MemberStatus = "active"
	MemberStatusInactive MemberStatus = "inactive"
)

type Member struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	UserID      *uint          `gorm:"index" json:"user_id,omitempty"`
	User        *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	FullName    string         `gorm:"size:255;not null" json:"full_name"`
	BirthDate   time.Time      `json:"birth_date"`
	Gender      string         `gorm:"type:enum('male','female');not null" json:"gender"`
	Section     Section        `gorm:"type:enum('ashbal','kashaf','jawala','mukashe');not null" json:"section"`
	RankStage   string         `gorm:"size:100" json:"rank_stage"`
	JoinDate    time.Time      `json:"join_date"`
	PhotoURL    *string        `gorm:"size:500" json:"photo_url,omitempty"`
	ParentName  string         `gorm:"size:255" json:"parent_name"`
	ParentPhone string         `gorm:"size:20" json:"parent_phone"`
	SecondPhone string         `gorm:"size:20" json:"secondary_phone"`
	Address     string         `gorm:"size:500" json:"address"`
	XPTotal     int            `gorm:"default:0" json:"xp_total"`
	Level       int            `gorm:"default:1" json:"level"`
	Status      MemberStatus   `gorm:"type:enum('active','inactive');default:'active'" json:"status"`

	Medical    *MemberMedical  `gorm:"foreignKey:MemberID" json:"medical,omitempty"`
	Badges     []MemberBadge   `gorm:"foreignKey:MemberID" json:"badges,omitempty"`
	Skills     []MemberSkill   `gorm:"foreignKey:MemberID" json:"skills,omitempty"`
	Units      []UnitMember    `gorm:"foreignKey:MemberID" json:"units,omitempty"`
}

type MemberMedical struct {
	ID                uint      `gorm:"primarykey" json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	MemberID          uint      `gorm:"not null;uniqueIndex" json:"member_id"`
	BloodType         string    `gorm:"size:10" json:"blood_type"`
	Allergies         string    `gorm:"type:text" json:"allergies"`
	ChronicConditions string    `gorm:"type:text" json:"chronic_conditions"`
	Medications       string    `gorm:"type:text" json:"medications"`
	EmergencyNotes    string    `gorm:"type:text" json:"emergency_notes"`
}
