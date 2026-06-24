package models

import (
	"time"

	"gorm.io/gorm"
)

type Equipment struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	Name              string         `gorm:"size:255;not null" json:"name"`
	Category          string         `gorm:"size:100" json:"category"`
	QuantityTotal     int            `gorm:"default:0" json:"quantity_total"`
	QuantityAvailable int            `gorm:"default:0" json:"quantity_available"`
	Condition         string         `gorm:"size:100" json:"condition"`
	Notes             string         `gorm:"type:text" json:"notes"`
}

type EquipmentLoan struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	EquipmentID uint       `gorm:"not null;index" json:"equipment_id"`
	BorrowedBy  uint       `gorm:"not null;index" json:"borrowed_by"`
	ActivityID  *uint      `gorm:"index" json:"activity_id,omitempty"`
	Quantity    int        `gorm:"default:1" json:"quantity"`
	DueDate     time.Time  `json:"due_date"`
	ReturnedAt  *time.Time `json:"returned_at,omitempty"`
	Equipment   Equipment  `gorm:"foreignKey:EquipmentID" json:"equipment,omitempty"`
	Borrower    User       `gorm:"foreignKey:BorrowedBy" json:"borrower,omitempty"`
	Activity    *Activity  `gorm:"foreignKey:ActivityID" json:"activity,omitempty"`
}

type XPSource string

const (
	XPSourceAttendance XPSource = "attendance"
	XPSourceBadge      XPSource = "badge"
	XPSourceQuiz       XPSource = "quiz"
	XPSourceLeadership XPSource = "leadership"
	XPSourceManual     XPSource = "manual"
)

type XPEvent struct {
	ID       uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
	MemberID uint      `gorm:"not null;index" json:"member_id"`
	UnitID   *uint     `gorm:"index" json:"unit_id,omitempty"`
	Source   XPSource  `gorm:"type:enum('attendance','badge','quiz','leadership','manual');not null" json:"source"`
	Points   int       `gorm:"not null" json:"points"`
	RefID    *uint     `json:"ref_id,omitempty"`
	Note     string    `gorm:"size:500" json:"note"`
	Member   Member    `gorm:"foreignKey:MemberID" json:"member,omitempty"`
}
