package dto

type CreateLessonRequest struct {
	Title       string `json:"title" validate:"required"`
	Category    string `json:"category" validate:"required"`
	Content     string `json:"content"`
	OrderIndex  int    `json:"order_index"`
	IsPublished bool   `json:"is_published"`
}

type UpdateLessonRequest struct {
	Title       *string `json:"title"`
	Category    *string `json:"category"`
	Content     *string `json:"content"`
	OrderIndex  *int    `json:"order_index"`
	IsPublished *bool   `json:"is_published"`
}

type CreateQuizRequest struct {
	LessonID  uint            `json:"lesson_id" validate:"required"`
	Title     string          `json:"title" validate:"required"`
	PassScore int             `json:"pass_score" validate:"min=0,max=100"`
	XPReward  int             `json:"xp_reward" validate:"min=0"`
	Questions []QuestionInput `json:"questions" validate:"required,min=1"`
}

type QuestionInput struct {
	Text         string   `json:"text" validate:"required"`
	Options      []string `json:"options" validate:"required,min=2"`
	CorrectIndex int      `json:"correct_index" validate:"min=0"`
	Points       int      `json:"points" validate:"min=1"`
}

type QuizAttemptRequest struct {
	Answers []int `json:"answers" validate:"required"`
}
