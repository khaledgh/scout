package services

import (
	"context"
	"encoding/json"
	"kashfi/internal/config"
	"kashfi/internal/models"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type TrainingService struct {
	db     *gorm.DB
	gamify *GamificationService
	cfg    *config.Config
}

func NewTrainingService(db *gorm.DB, gamify *GamificationService, cfg *config.Config) *TrainingService {
	return &TrainingService{db: db, gamify: gamify, cfg: cfg}
}

func (s *TrainingService) UploadPath() string  { return s.cfg.Upload.Dir }
func (s *TrainingService) UploadMaxMB() int64  { return s.cfg.Upload.MaxSizeMB }
func (s *TrainingService) UploadTypes() string { return s.cfg.Upload.AllowedTypes }

func (s *TrainingService) SetCover(ctx context.Context, id uint, filename string) error {
	return s.db.WithContext(ctx).Model(&models.TrainingLesson{}).Where("id = ?", id).
		Update("cover_url", filename).Error
}

func (s *TrainingService) ListLessons(ctx context.Context) ([]models.TrainingLesson, error) {
	var lessons []models.TrainingLesson
	err := s.db.WithContext(ctx).Where("is_published = true").
		Order("order_index").Find(&lessons).Error
	return lessons, err
}

func (s *TrainingService) GetLesson(ctx context.Context, id uint) (*models.TrainingLesson, error) {
	var lesson models.TrainingLesson
	err := s.db.WithContext(ctx).Preload("Quizzes").Preload("Media").First(&lesson, id).Error
	if err != nil {
		return nil, ErrNotFound
	}
	return &lesson, nil
}

func (s *TrainingService) AddMedia(ctx context.Context, lessonID, uploaderID uint, filename string, mediaType models.MediaType) (*models.TrainingLessonMedia, error) {
	url := filepath.ToSlash(filepath.Join(s.cfg.Upload.PublicPath, filename))
	media := &models.TrainingLessonMedia{
		LessonID:   lessonID,
		URL:        url,
		MediaType:  mediaType,
		UploadedBy: uploaderID,
	}
	if err := s.db.WithContext(ctx).Create(media).Error; err != nil {
		return nil, err
	}
	return media, nil
}

func (s *TrainingService) DeleteMedia(ctx context.Context, mediaID uint) error {
	result := s.db.WithContext(ctx).Delete(&models.TrainingLessonMedia{}, mediaID)
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}

type CreateLessonInput struct {
	Title       string
	Category    string
	Content     string
	OrderIndex  int
	IsPublished bool
}

func (s *TrainingService) CreateLesson(ctx context.Context, in CreateLessonInput) (*models.TrainingLesson, error) {
	lesson := &models.TrainingLesson{
		Title:       in.Title,
		Category:    in.Category,
		Content:     in.Content,
		OrderIndex:  in.OrderIndex,
		IsPublished: in.IsPublished,
	}
	if err := s.db.WithContext(ctx).Create(lesson).Error; err != nil {
		return nil, err
	}
	return lesson, nil
}

type UpdateLessonInput struct {
	Title       *string
	Category    *string
	Content     *string
	OrderIndex  *int
	IsPublished *bool
}

func (s *TrainingService) UpdateLesson(ctx context.Context, id uint, in UpdateLessonInput) (*models.TrainingLesson, error) {
	var lesson models.TrainingLesson
	if err := s.db.WithContext(ctx).First(&lesson, id).Error; err != nil {
		return nil, ErrNotFound
	}
	updates := map[string]interface{}{}
	if in.Title != nil { updates["title"] = *in.Title }
	if in.Category != nil { updates["category"] = *in.Category }
	if in.Content != nil { updates["content"] = *in.Content }
	if in.OrderIndex != nil { updates["order_index"] = *in.OrderIndex }
	if in.IsPublished != nil { updates["is_published"] = *in.IsPublished }
	s.db.WithContext(ctx).Model(&lesson).Updates(updates)
	return &lesson, nil
}

func (s *TrainingService) DeleteLesson(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&models.TrainingLesson{}, id)
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}

func (s *TrainingService) GetQuiz(ctx context.Context, lessonID uint) (*models.Quiz, error) {
	var quiz models.Quiz
	err := s.db.WithContext(ctx).Where("lesson_id = ?", lessonID).
		Preload("Questions").First(&quiz).Error
	if err != nil {
		return nil, ErrNotFound
	}
	return &quiz, nil
}

type QuizInput struct {
	LessonID  uint
	Title     string
	PassScore int
	XPReward  int
	Questions []QuestionInput
}

type QuestionInput struct {
	Text         string
	Options      []string
	CorrectIndex int
	Points       int
}

func (s *TrainingService) CreateQuiz(ctx context.Context, in QuizInput) (*models.Quiz, error) {
	var result *models.Quiz
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		quiz := &models.Quiz{
			LessonID:  in.LessonID,
			Title:     in.Title,
			PassScore: in.PassScore,
			XPReward:  in.XPReward,
		}
		if err := tx.Create(quiz).Error; err != nil {
			return err
		}
		for _, q := range in.Questions {
			opts, _ := json.Marshal(q.Options)
			qq := models.QuizQuestion{
				QuizID:       quiz.ID,
				Text:         q.Text,
				OptionsJSON:  string(opts),
				CorrectIndex: q.CorrectIndex,
				Points:       q.Points,
			}
			if err := tx.Create(&qq).Error; err != nil {
				return err
			}
		}
		result = quiz
		return nil
	})
	return result, err
}

func (s *TrainingService) SubmitAttempt(ctx context.Context, quizID, memberID uint, answers []int) (*models.QuizAttempt, error) {
	var quiz models.Quiz
	if err := s.db.WithContext(ctx).Preload("Questions").First(&quiz, quizID).Error; err != nil {
		return nil, ErrNotFound
	}

	score := 0
	totalPoints := 0
	for i, q := range quiz.Questions {
		totalPoints += q.Points
		if i < len(answers) && answers[i] == q.CorrectIndex {
			score += q.Points
		}
	}
	pct := 0
	if totalPoints > 0 {
		pct = score * 100 / totalPoints
	}
	passed := pct >= quiz.PassScore

	answersJSON, _ := json.Marshal(answers)
	attempt := &models.QuizAttempt{
		QuizID:      quizID,
		MemberID:    memberID,
		Score:       pct,
		Passed:      passed,
		AnswersJSON: string(answersJSON),
		AttemptedAt: time.Now(),
	}
	if err := s.db.WithContext(ctx).Create(attempt).Error; err != nil {
		return nil, err
	}

	if passed && quiz.XPReward > 0 {
		refID := quizID
		s.gamify.AwardXP(ctx, memberID, nil, models.XPSourceQuiz, quiz.XPReward, &refID, "quiz: "+quiz.Title)
	}

	return attempt, nil
}

type LessonProgress struct {
	LessonID    uint   `json:"lesson_id"`
	LessonTitle string `json:"lesson_title"`
	Passed      bool   `json:"passed"`
	BestScore   int    `json:"best_score"`
}

func (s *TrainingService) MyProgress(ctx context.Context, memberID uint) ([]LessonProgress, error) {
	var lessons []models.TrainingLesson
	s.db.WithContext(ctx).Where("is_published = true").Order("order_index").Find(&lessons)

	result := []LessonProgress{}
	for _, lesson := range lessons {
		var quiz models.Quiz
		if s.db.WithContext(ctx).Where("lesson_id = ?", lesson.ID).First(&quiz).Error != nil {
			result = append(result, LessonProgress{LessonID: lesson.ID, LessonTitle: lesson.Title})
			continue
		}
		var best models.QuizAttempt
		s.db.WithContext(ctx).Where("quiz_id = ? AND member_id = ?", quiz.ID, memberID).
			Order("score DESC").First(&best)
		result = append(result, LessonProgress{
			LessonID:    lesson.ID,
			LessonTitle: lesson.Title,
			Passed:      best.Passed,
			BestScore:   best.Score,
		})
	}
	return result, nil
}
