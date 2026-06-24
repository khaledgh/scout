package models

import (
	"time"

	"gorm.io/gorm"
)

type Announcement struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Body        string         `gorm:"type:text;not null" json:"body"`
	Audience    string         `gorm:"type:enum('all','unit','leaders');default:'all'" json:"audience"`
	UnitID      *uint          `gorm:"index" json:"unit_id,omitempty"`
	AuthorID    uint           `gorm:"not null;index" json:"author_id"`
	Pinned      bool           `gorm:"default:false" json:"pinned"`
	PublishedAt *time.Time     `json:"published_at,omitempty"`
	Author      User           `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Unit        *Unit          `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
}

type ChannelType string

const (
	ChannelUnit   ChannelType = "unit"
	ChannelGroup  ChannelType = "group"
	ChannelDirect ChannelType = "direct"
)

type Channel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Type      ChannelType    `gorm:"type:enum('unit','group','direct');not null" json:"type"`
	UnitID    *uint          `gorm:"index" json:"unit_id,omitempty"`
	Name      string         `gorm:"size:255" json:"name"`
	Unit      *Unit          `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	Messages  []Message      `gorm:"foreignKey:ChannelID" json:"messages,omitempty"`
}

type Message struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time `gorm:"index" json:"created_at"`
	ChannelID     uint      `gorm:"not null;index" json:"channel_id"`
	SenderID      uint      `gorm:"not null;index" json:"sender_id"`
	Body          string    `gorm:"type:text" json:"body"`
	AttachmentURL *string   `gorm:"size:500" json:"attachment_url,omitempty"`
	Channel       Channel   `gorm:"foreignKey:ChannelID" json:"channel,omitempty"`
	Sender        User      `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}

type Notification struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	CreatedAt time.Time  `gorm:"index" json:"created_at"`
	UserID    uint       `gorm:"not null;index" json:"user_id"`
	Title     string     `gorm:"size:255;not null" json:"title"`
	Body      string     `gorm:"type:text" json:"body"`
	Type      string     `gorm:"size:100" json:"type"`
	DataJSON  string     `gorm:"type:json" json:"data_json"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	User      User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
