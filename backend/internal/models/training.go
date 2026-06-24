package models

import (
	"time"

	"gorm.io/gorm"
)

type TrainingLesson struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Category    string         `gorm:"size:100;not null" json:"category"`
	Content     string         `gorm:"type:longtext" json:"content"`
	CoverURL    *string        `gorm:"size:500" json:"cover_url,omitempty"`
	MediaJSON   *string        `gorm:"type:text" json:"media_json,omitempty"`
	OrderIndex  int            `gorm:"default:0" json:"order_index"`
	IsPublished bool           `gorm:"default:false" json:"is_published"`

	Quizzes []Quiz                 `gorm:"foreignKey:LessonID" json:"quizzes,omitempty"`
	Media   []TrainingLessonMedia  `gorm:"foreignKey:LessonID" json:"media,omitempty"`
}

type TrainingLessonMedia struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	LessonID   uint      `gorm:"not null;index" json:"lesson_id"`
	URL        string    `gorm:"size:500;not null" json:"url"`
	MediaType  MediaType `gorm:"type:enum('image','video');not null" json:"media_type"`
	UploadedBy uint      `gorm:"index" json:"uploaded_by"`
	Caption    string    `gorm:"size:500" json:"caption"`
}

type Quiz struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	LessonID  uint           `gorm:"not null;index" json:"lesson_id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	PassScore int            `gorm:"default:70" json:"pass_score"`
	XPReward  int            `gorm:"default:0" json:"xp_reward"`

	Lesson    TrainingLesson `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
	Questions []QuizQuestion `gorm:"foreignKey:QuizID" json:"questions,omitempty"`
}

type QuizQuestion struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	QuizID       uint      `gorm:"not null;index" json:"quiz_id"`
	Text         string    `gorm:"type:text;not null" json:"text"`
	OptionsJSON  string    `gorm:"type:json;not null" json:"options_json"`
	CorrectIndex int       `gorm:"not null" json:"correct_index"`
	Points       int       `gorm:"default:1" json:"points"`
}

type QuizAttempt struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	QuizID      uint      `gorm:"not null;index" json:"quiz_id"`
	MemberID    uint      `gorm:"not null;index" json:"member_id"`
	Score       int       `json:"score"`
	Passed      bool      `json:"passed"`
	AnswersJSON string    `gorm:"type:json" json:"answers_json"`
	AttemptedAt time.Time `json:"attempted_at"`
	Quiz        Quiz      `gorm:"foreignKey:QuizID" json:"quiz,omitempty"`
	Member      Member    `gorm:"foreignKey:MemberID" json:"member,omitempty"`
}
