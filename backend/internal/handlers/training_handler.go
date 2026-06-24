package handlers

import (
	"encoding/json"
	"kashfi/internal/dto"
	appMiddleware "kashfi/internal/middleware"
	"kashfi/internal/services"
	"kashfi/internal/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TrainingHandler struct {
	svc *services.TrainingService
}

func NewTrainingHandler(svc *services.TrainingService) *TrainingHandler {
	return &TrainingHandler{svc: svc}
}

func (h *TrainingHandler) ListLessons(c echo.Context) error {
	lessons, err := h.svc.ListLessons(c.Request().Context())
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, lessons)
}

func (h *TrainingHandler) GetLesson(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	lesson, err := h.svc.GetLesson(c.Request().Context(), id)
	if err != nil {
		return utils.NotFound(c, "lesson")
	}
	return utils.OK(c, lesson)
}

func (h *TrainingHandler) CreateLesson(c echo.Context) error {
	var req dto.CreateLessonRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	if err := utils.Validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	lesson, err := h.svc.CreateLesson(c.Request().Context(), services.CreateLessonInput{
		Title: req.Title, Category: req.Category, Content: req.Content,
		OrderIndex: req.OrderIndex, IsPublished: req.IsPublished,
	})
	if err != nil {
		return utils.Internal(c)
	}
	return utils.Created(c, lesson)
}

func (h *TrainingHandler) UpdateLesson(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	var req dto.UpdateLessonRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	lesson, err := h.svc.UpdateLesson(c.Request().Context(), id, services.UpdateLessonInput{
		Title: req.Title, Category: req.Category, Content: req.Content,
		OrderIndex: req.OrderIndex, IsPublished: req.IsPublished,
	})
	if err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "lesson") }
		return utils.Internal(c)
	}
	return utils.OK(c, lesson)
}

func (h *TrainingHandler) DeleteLesson(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	if err := h.svc.DeleteLesson(c.Request().Context(), id); err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "lesson") }
		return utils.Internal(c)
	}
	return utils.NoContent(c)
}

func (h *TrainingHandler) CreateQuiz(c echo.Context) error {
	lessonID, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid lesson id")
	}
	var req dto.CreateQuizRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	if err := utils.Validate.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	questions := make([]services.QuestionInput, len(req.Questions))
	for i, q := range req.Questions {
		questions[i] = services.QuestionInput{
			Text: q.Text, Options: q.Options, CorrectIndex: q.CorrectIndex, Points: q.Points,
		}
	}
	quiz, err := h.svc.CreateQuiz(c.Request().Context(), services.QuizInput{
		LessonID: lessonID, Title: req.Title, PassScore: req.PassScore,
		XPReward: req.XPReward, Questions: questions,
	})
	if err != nil {
		return utils.Internal(c)
	}
	return utils.Created(c, quiz)
}

func (h *TrainingHandler) GetQuiz(c echo.Context) error {
	lessonID, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid lesson id")
	}
	quiz, err := h.svc.GetQuiz(c.Request().Context(), lessonID)
	if err != nil {
		return utils.NotFound(c, "quiz")
	}
	return utils.OK(c, quiz)
}

func (h *TrainingHandler) SubmitAttempt(c echo.Context) error {
	quizID, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid quiz id")
	}
	var req dto.QuizAttemptRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid body")
	}
	memberID := appMiddleware.GetUserID(c)
	attempt, err := h.svc.SubmitAttempt(c.Request().Context(), quizID, memberID, req.Answers)
	if err != nil {
		if err == services.ErrNotFound { return utils.NotFound(c, "quiz") }
		return utils.Internal(c)
	}
	return utils.Created(c, attempt)
}

func (h *TrainingHandler) MyProgress(c echo.Context) error {
	memberID := appMiddleware.GetUserID(c)
	progress, err := h.svc.MyProgress(c.Request().Context(), memberID)
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, progress)
}

func (h *TrainingHandler) UploadLessonMedia(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	file, err := c.FormFile("file")
	if err != nil {
		return utils.BadRequest(c, "file required")
	}
	filename, err := utils.SaveUpload(file, h.svc.UploadPath(), h.svc.UploadTypes(), h.svc.UploadMaxMB())
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}
	mediaType := detectMediaType(file)
	media, err := h.svc.AddMedia(c.Request().Context(), id, appMiddleware.GetUserID(c), filename, mediaType)
	if err != nil {
		return utils.Internal(c)
	}
	return utils.Created(c, media)
}

func (h *TrainingHandler) DeleteLessonMedia(c echo.Context) error {
	mid, err := strconv.ParseUint(c.Param("mid"), 10, 64)
	if err != nil {
		return utils.BadRequest(c, "invalid media id")
	}
	if err := h.svc.DeleteMedia(c.Request().Context(), uint(mid)); err != nil {
		if err == services.ErrNotFound {
			return utils.NotFound(c, "media")
		}
		return utils.Internal(c)
	}
	return utils.NoContent(c)
}

func (h *TrainingHandler) UploadLessonCover(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return utils.BadRequest(c, "invalid id")
	}
	file, err := c.FormFile("file")
	if err != nil {
		return utils.BadRequest(c, "file required")
	}
	filename, err := utils.SaveUpload(file, h.svc.UploadPath(), "image/jpeg,image/png,image/webp", h.svc.UploadMaxMB())
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}
	if err := h.svc.SetCover(c.Request().Context(), id, filename); err != nil {
		if err == services.ErrNotFound {
			return utils.NotFound(c, "lesson")
		}
		return utils.Internal(c)
	}
	return utils.OK(c, map[string]string{"filename": filename})
}

var _ = json.Marshal // used in service
