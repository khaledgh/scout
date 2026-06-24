package main

import (
	"kashfi/internal/config"
	"kashfi/internal/handlers"
	appMiddleware "kashfi/internal/middleware"
	"kashfi/internal/services"
	"kashfi/internal/ws"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func registerRoutes(api *echo.Group, cfg *config.Config, database *gorm.DB, hub *ws.Hub) {
	// ── Services (order matters due to dependencies) ──────────────────────────
	gamifySvc    := services.NewGamificationService(database, cfg)
	authSvc      := services.NewAuthService(database, cfg)
	memberSvc    := services.NewMemberService(database, gamifySvc, cfg)
	unitSvc      := services.NewUnitService(database)
	activitySvc  := services.NewActivityService(database, cfg)
	attendSvc    := services.NewAttendanceService(database, cfg, gamifySvc, hub)
	badgeSvc     := services.NewBadgeService(database, gamifySvc)
	trainingSvc  := services.NewTrainingService(database, gamifySvc, cfg)
	notifSvc     := services.NewNotificationService(database)
	chatSvc      := services.NewChatService(database, hub)
	announceSvc  := services.NewAnnouncementService(database, notifSvc)
	reportSvc    := services.NewReportService(database)
	equipmentSvc := services.NewEquipmentService(database)

	// ── Handlers ─────────────────────────────────────────────────────────────
	authH     := handlers.NewAuthHandler(authSvc)
	memberH   := handlers.NewMemberHandler(memberSvc, cfg.Geo.QRSigningSecret)
	unitH     := handlers.NewUnitHandler(unitSvc)
	activityH := handlers.NewActivityHandler(activitySvc, attendSvc)
	badgeH    := handlers.NewBadgeHandler(badgeSvc, gamifySvc)
	trainingH := handlers.NewTrainingHandler(trainingSvc)
	announceH := handlers.NewAnnouncementHandler(announceSvc, notifSvc)
	notifH    := handlers.NewNotificationHandler(notifSvc)
	chatH     := handlers.NewChatHandler(chatSvc, hub, cfg)
	reportH   := handlers.NewReportHandler(reportSvc)
	equipmentH := handlers.NewEquipmentHandler(equipmentSvc)

	authRL := appMiddleware.RateLimit(10, time.Minute)
	jwtMW  := appMiddleware.JWT(cfg.JWT.Secret)

	leaderRoles := []string{"super_admin", "leader", "assistant"}
	adminRoles  := []string{"super_admin"}

	// ── Auth (public) ─────────────────────────────────────────────────────────
	auth := api.Group("/auth", authRL)
	auth.POST("/login", authH.Login)
	auth.POST("/refresh", authH.Refresh)
	auth.POST("/logout", authH.Logout)
	auth.GET("/me", authH.Me, jwtMW)
	auth.POST("/change-password", authH.ChangePassword, jwtMW)

	// ── Authenticated ─────────────────────────────────────────────────────────
	r := api.Group("", jwtMW)

	// Members
	r.GET("/members", memberH.List)
	r.POST("/members", memberH.Create, appMiddleware.RequireRole(leaderRoles...))
	r.GET("/members/:id", memberH.Get)
	r.PUT("/members/:id", memberH.Update, appMiddleware.RequireRole(leaderRoles...))
	r.DELETE("/members/:id", memberH.Delete, appMiddleware.RequireRole(leaderRoles...))
	r.POST("/members/:id/photo", memberH.UploadPhoto, appMiddleware.RequireRole(leaderRoles...))
	r.GET("/members/:id/medical", memberH.GetMedical, appMiddleware.RequireRole(leaderRoles...))
	r.PUT("/members/:id/medical", memberH.UpsertMedical, appMiddleware.RequireRole(leaderRoles...))
	r.GET("/members/:id/timeline", memberH.Timeline)
	r.POST("/members/:id/evaluate", memberH.CreateEvaluation, appMiddleware.RequireRole(leaderRoles...))
	r.GET("/members/:id/qr", memberH.GetQRToken)
	r.GET("/members/:id/badges", badgeH.MemberBadges)
	r.POST("/members/:id/badges", badgeH.AwardBadge, appMiddleware.RequireRole(leaderRoles...))
	r.DELETE("/members/:id/badges/:badgeId", badgeH.RevokeBadge, appMiddleware.RequireRole(leaderRoles...))
	r.POST("/members/:id/skills", badgeH.AssessSkill, appMiddleware.RequireRole(leaderRoles...))
	r.GET("/members/:id/xp", badgeH.XPHistory)

	// Units
	r.GET("/units", unitH.List)
	r.POST("/units", unitH.Create, appMiddleware.RequireRole(adminRoles...))
	r.GET("/units/leaderboard", unitH.Leaderboard)
	r.GET("/units/:id", unitH.Get)
	r.PUT("/units/:id", unitH.Update, appMiddleware.RequireRole(leaderRoles...))
	r.DELETE("/units/:id", unitH.Delete, appMiddleware.RequireRole(adminRoles...))
	r.POST("/units/:id/members", unitH.AddMembers, appMiddleware.RequireRole(leaderRoles...))
	r.DELETE("/units/:id/members/:mid", unitH.RemoveMember, appMiddleware.RequireRole(leaderRoles...))
	r.POST("/units/:id/leaders", unitH.AssignLeader, appMiddleware.RequireRole(adminRoles...))

	// Activities
	r.GET("/activities", activityH.List)
	r.POST("/activities", activityH.Create, appMiddleware.RequireRole(leaderRoles...))
	r.GET("/activities/:id", activityH.Get)
	r.PUT("/activities/:id", activityH.Update, appMiddleware.RequireRole(leaderRoles...))
	r.DELETE("/activities/:id", activityH.Delete, appMiddleware.RequireRole(leaderRoles...))
	r.POST("/activities/:id/media", activityH.UploadMedia, appMiddleware.RequireRole(leaderRoles...))
	r.GET("/activities/:id/attendance", activityH.GetAttendance)
	r.POST("/activities/:id/attendance", activityH.RecordAttendance, appMiddleware.RequireRole(leaderRoles...))
	r.POST("/activities/:id/checkin", activityH.CheckIn, appMiddleware.RateLimit(30, time.Minute))
	r.POST("/activities/:id/feedback", activityH.CreateFeedback)
	r.GET("/activities/:id/feedback/summary", activityH.FeedbackSummary, appMiddleware.RequireRole(leaderRoles...))

	// Badges & Skills catalog
	r.GET("/badges", badgeH.ListBadges)
	r.POST("/badges", badgeH.CreateBadge, appMiddleware.RequireRole(leaderRoles...))
	r.GET("/badges/:id", badgeH.GetBadge)
	r.PUT("/badges/:id", badgeH.UpdateBadge, appMiddleware.RequireRole(adminRoles...))
	r.DELETE("/badges/:id", badgeH.DeleteBadge, appMiddleware.RequireRole(adminRoles...))
	r.GET("/skills", badgeH.ListSkills)
	r.POST("/skills", badgeH.CreateSkill, appMiddleware.RequireRole(adminRoles...))
	r.PUT("/skills/:id", badgeH.UpdateSkill, appMiddleware.RequireRole(adminRoles...))
	r.DELETE("/skills/:id", badgeH.DeleteSkill, appMiddleware.RequireRole(adminRoles...))

	// Training
	r.GET("/training/lessons", trainingH.ListLessons)
	r.POST("/training/lessons", trainingH.CreateLesson, appMiddleware.RequireRole(leaderRoles...))
	r.GET("/training/lessons/:id", trainingH.GetLesson)
	r.PUT("/training/lessons/:id", trainingH.UpdateLesson, appMiddleware.RequireRole(leaderRoles...))
	r.DELETE("/training/lessons/:id", trainingH.DeleteLesson, appMiddleware.RequireRole(leaderRoles...))
	r.GET("/training/lessons/:id/quiz", trainingH.GetQuiz)
	r.POST("/training/lessons/:id/quiz", trainingH.CreateQuiz, appMiddleware.RequireRole(leaderRoles...))
	r.POST("/training/quizzes/:id/attempt", trainingH.SubmitAttempt)
	r.POST("/training/lessons/:id/cover", trainingH.UploadLessonCover, appMiddleware.RequireRole(leaderRoles...))
	r.POST("/training/lessons/:id/media", trainingH.UploadLessonMedia, appMiddleware.RequireRole(leaderRoles...))
	r.DELETE("/training/lessons/:id/media/:mid", trainingH.DeleteLessonMedia, appMiddleware.RequireRole(leaderRoles...))
	r.GET("/training/me/progress", trainingH.MyProgress)

	// Gamification
	r.GET("/leaderboard/members", badgeH.MemberLeaderboard)
	r.GET("/leaderboard/units", badgeH.UnitLeaderboard)

	// Communication
	r.GET("/announcements", announceH.List)
	r.POST("/announcements", announceH.Create, appMiddleware.RequireRole(leaderRoles...))
	r.GET("/announcements/:id", announceH.Get)
	r.PUT("/announcements/:id", announceH.Update, appMiddleware.RequireRole(leaderRoles...))
	r.DELETE("/announcements/:id", announceH.Delete, appMiddleware.RequireRole(leaderRoles...))
	r.GET("/channels", chatH.Channels)
	r.GET("/channels/:id/messages", chatH.Messages)
	r.GET("/notifications", notifH.List)
	r.PUT("/notifications/:id/read", notifH.MarkRead)

	// Dashboard & Reports
	r.GET("/dashboard", reportH.Dashboard, appMiddleware.RequireRole(leaderRoles...))
	r.GET("/reports/member/:id", reportH.MemberReport, appMiddleware.RequireRole(leaderRoles...))
	r.GET("/reports/unit/:id", reportH.UnitReport, appMiddleware.RequireRole(leaderRoles...))
	r.GET("/reports/monthly", reportH.Monthly, appMiddleware.RequireRole(leaderRoles...))
	r.GET("/reports/export", reportH.Export, appMiddleware.RequireRole(leaderRoles...))

	// Equipment
	r.GET("/equipment", equipmentH.List)
	r.POST("/equipment", equipmentH.Create, appMiddleware.RequireRole(leaderRoles...))
	r.GET("/equipment/:id", equipmentH.Get)
	r.PUT("/equipment/:id", equipmentH.Update, appMiddleware.RequireRole(leaderRoles...))
	r.DELETE("/equipment/:id", equipmentH.Delete, appMiddleware.RequireRole(leaderRoles...))
	r.GET("/equipment/:id/loans", equipmentH.Loans)
	r.POST("/equipment/:id/loan", equipmentH.Loan, appMiddleware.RequireRole(leaderRoles...))
	r.POST("/equipment/loans/:loanId/return", equipmentH.ReturnLoan, appMiddleware.RequireRole(leaderRoles...))

	// WebSocket
	api.GET("/ws/chat", chatH.WebSocket, jwtMW)
}
