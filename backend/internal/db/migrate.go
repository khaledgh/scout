package db

import (
	"kashfi/internal/models"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Member{},
		&models.MemberMedical{},
		&models.Unit{},
		&models.UnitLeader{},
		&models.UnitMember{},
		&models.Activity{},
		&models.ActivityAttendance{},
		&models.ActivityFeedback{},
		&models.ActivityMedia{},
		&models.Badge{},
		&models.MemberBadge{},
		&models.Skill{},
		&models.MemberSkill{},
		&models.Evaluation{},
		&models.TrainingLesson{},
		&models.TrainingLessonMedia{},
		&models.Quiz{},
		&models.QuizQuestion{},
		&models.QuizAttempt{},
		&models.Announcement{},
		&models.Channel{},
		&models.Message{},
		&models.Notification{},
		&models.Equipment{},
		&models.EquipmentLoan{},
		&models.XPEvent{},
	)
}
