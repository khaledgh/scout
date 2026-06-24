package models

import (
	"time"

	"gorm.io/gorm"
)

type ActivityType string

const (
	ActivityTypeCamp     ActivityType = "camp"
	ActivityTypeHike     ActivityType = "hike"
	ActivityTypeTraining ActivityType = "training"
	ActivityTypeMeeting  ActivityType = "meeting"
	ActivityTypeService  ActivityType = "service"
)

type ActivityStatus string

const (
	ActivityStatusPlanned   ActivityStatus = "planned"
	ActivityStatusOngoing   ActivityStatus = "ongoing"
	ActivityStatusCompleted ActivityStatus = "completed"
	ActivityStatusCancelled ActivityStatus = "cancelled"
)

type Activity struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	Title             string         `gorm:"size:255;not null" json:"title"`
	Description       string         `gorm:"type:text" json:"description"`
	Type              ActivityType   `gorm:"type:enum('camp','hike','training','meeting','service');not null" json:"type"`
	Location          string         `gorm:"size:500" json:"location"`
	LocationLat       *float64       `json:"location_lat,omitempty"`
	LocationLng       *float64       `json:"location_lng,omitempty"`
	StartsAt          time.Time      `gorm:"index" json:"starts_at"`
	EndsAt            time.Time      `json:"ends_at"`
	ResponsibleUserID uint           `gorm:"index" json:"responsible_user_id"`
	UnitID            *uint          `gorm:"index" json:"unit_id,omitempty"`
	Status            ActivityStatus `gorm:"type:enum('planned','ongoing','completed','cancelled');default:'planned'" json:"status"`
	CoverImageURL     *string        `gorm:"size:500" json:"cover_image_url,omitempty"`

	ResponsibleUser User                  `gorm:"foreignKey:ResponsibleUserID" json:"responsible_user,omitempty"`
	Unit            *Unit                 `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	Attendances     []ActivityAttendance  `gorm:"foreignKey:ActivityID" json:"attendances,omitempty"`
	Feedback        []ActivityFeedback    `gorm:"foreignKey:ActivityID" json:"feedback,omitempty"`
	Media           []ActivityMedia       `gorm:"foreignKey:ActivityID" json:"media,omitempty"`
}

type AttendanceStatus string

const (
	AttendancePresent AttendanceStatus = "present"
	AttendanceAbsent  AttendanceStatus = "absent"
	AttendanceExcused AttendanceStatus = "excused"
	AttendanceLate    AttendanceStatus = "late"
)

type CheckInMethod string

const (
	CheckInQR     CheckInMethod = "qr"
	CheckInGPS    CheckInMethod = "gps"
	CheckInManual CheckInMethod = "manual"
)

type ActivityAttendance struct {
	ID             uint             `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	ActivityID     uint             `gorm:"not null;index" json:"activity_id"`
	MemberID       uint             `gorm:"not null;index" json:"member_id"`
	Status         AttendanceStatus `gorm:"type:enum('present','absent','excused','late');not null" json:"status"`
	CheckInAt      *time.Time       `json:"check_in_at,omitempty"`
	CheckInMethod  *CheckInMethod   `gorm:"type:enum('qr','gps','manual')" json:"check_in_method,omitempty"`
	Lat            *float64         `json:"lat,omitempty"`
	Lng            *float64         `json:"lng,omitempty"`
	RecordedBy     uint             `gorm:"index" json:"recorded_by"`
	Activity       Activity         `gorm:"foreignKey:ActivityID" json:"activity,omitempty"`
	Member         Member           `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	RecordedByUser User             `gorm:"foreignKey:RecordedBy" json:"recorded_by_user,omitempty"`
}

type ActivityFeedback struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	ActivityID    uint      `gorm:"not null;index" json:"activity_id"`
	MemberID      uint      `gorm:"not null;index" json:"member_id"`
	Rating        int       `gorm:"not null" json:"rating"`
	WhatWentWell  string    `gorm:"type:text" json:"what_went_well"`
	WhatToImprove string    `gorm:"type:text" json:"what_to_improve"`
	Comment       string    `gorm:"type:text" json:"comment"`
	Member        Member    `gorm:"foreignKey:MemberID" json:"member,omitempty"`
}

type MediaType string

const (
	MediaImage MediaType = "image"
	MediaVideo MediaType = "video"
)

type ActivityMedia struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	ActivityID uint      `gorm:"not null;index" json:"activity_id"`
	URL        string    `gorm:"size:500;not null" json:"url"`
	MediaType  MediaType `gorm:"type:enum('image','video');not null" json:"media_type"`
	UploadedBy uint      `gorm:"index" json:"uploaded_by"`
	Caption    string    `gorm:"size:500" json:"caption"`
	Uploader   User      `gorm:"foreignKey:UploadedBy" json:"uploader,omitempty"`
}
