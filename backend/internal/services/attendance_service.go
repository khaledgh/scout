package services

import (
	"context"
	"encoding/json"
	"kashfi/internal/config"
	"kashfi/internal/models"
	"kashfi/internal/utils"
	"kashfi/internal/ws"
	"time"

	"gorm.io/gorm"
)

type AttendanceService struct {
	db     *gorm.DB
	cfg    *config.Config
	gamify *GamificationService
	hub    *ws.Hub
}

func NewAttendanceService(db *gorm.DB, cfg *config.Config, gamify *GamificationService, hub *ws.Hub) *AttendanceService {
	return &AttendanceService{db: db, cfg: cfg, gamify: gamify, hub: hub}
}

func (s *AttendanceService) GetForActivity(ctx context.Context, activityID uint) ([]models.ActivityAttendance, error) {
	var records []models.ActivityAttendance
	err := s.db.WithContext(ctx).Where("activity_id = ?", activityID).
		Preload("Member").Find(&records).Error
	return records, err
}

type AttendanceRecord struct {
	MemberID uint
	Status   models.AttendanceStatus
}

func (s *AttendanceService) BulkRecord(ctx context.Context, activityID, recordedBy uint, records []AttendanceRecord) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, r := range records {
			method := models.CheckInManual
			now := time.Now()
			att := models.ActivityAttendance{
				ActivityID:    activityID,
				MemberID:      r.MemberID,
				Status:        r.Status,
				CheckInAt:     &now,
				CheckInMethod: &method,
				RecordedBy:    recordedBy,
			}
			if err := tx.Where("activity_id = ? AND member_id = ?", activityID, r.MemberID).
				Assign(att).FirstOrCreate(&att).Error; err != nil {
				return err
			}
			if r.Status == models.AttendancePresent {
				s.awardAttendanceXP(ctx, tx, r.MemberID, activityID)
			}
		}
		s.broadcastCount(ctx, activityID)
		return nil
	})
}

func (s *AttendanceService) CheckInQR(ctx context.Context, activityID uint, token string, recordedBy uint) (*models.ActivityAttendance, error) {
	memberID, ok := utils.VerifyQRToken(token, s.cfg.Geo.QRSigningSecret, 24*time.Hour)
	if !ok {
		return nil, ErrBadRequest
	}
	return s.recordPresent(ctx, activityID, memberID, models.CheckInQR, nil, nil, recordedBy)
}

func (s *AttendanceService) CheckInGPS(ctx context.Context, activityID uint, memberID uint, lat, lng float64, recordedBy uint) (*models.ActivityAttendance, error) {
	var activity models.Activity
	if err := s.db.WithContext(ctx).First(&activity, activityID).Error; err != nil {
		return nil, ErrNotFound
	}
	if activity.LocationLat == nil || activity.LocationLng == nil {
		return nil, ErrBadRequest
	}
	dist := utils.HaversineMeters(lat, lng, *activity.LocationLat, *activity.LocationLng)
	if dist > s.cfg.Geo.GeofenceRadiusMeters {
		return nil, ErrForbidden
	}
	return s.recordPresent(ctx, activityID, memberID, models.CheckInGPS, &lat, &lng, recordedBy)
}

func (s *AttendanceService) recordPresent(ctx context.Context, activityID, memberID uint, method models.CheckInMethod, lat, lng *float64, recordedBy uint) (*models.ActivityAttendance, error) {
	now := time.Now()
	att := models.ActivityAttendance{
		ActivityID:    activityID,
		MemberID:      memberID,
		Status:        models.AttendancePresent,
		CheckInAt:     &now,
		CheckInMethod: &method,
		Lat:           lat,
		Lng:           lng,
		RecordedBy:    recordedBy,
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("activity_id = ? AND member_id = ?", activityID, memberID).
			Assign(att).FirstOrCreate(&att).Error; err != nil {
			return err
		}
		s.awardAttendanceXP(ctx, tx, memberID, activityID)
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.broadcastCount(ctx, activityID)
	return &att, nil
}

func (s *AttendanceService) awardAttendanceXP(ctx context.Context, tx *gorm.DB, memberID, activityID uint) {
	refID := activityID
	event := models.XPEvent{
		MemberID: memberID,
		Source:   models.XPSourceAttendance,
		Points:   s.cfg.Gamify.XPPerAttendance,
		RefID:    &refID,
		Note:     "activity attendance",
	}
	var existing models.XPEvent
	if tx.Where("member_id = ? AND ref_id = ? AND source = 'attendance'", memberID, activityID).First(&existing).Error == nil {
		return // already awarded
	}
	tx.Create(&event)
	tx.Model(&models.Member{}).Where("id = ?", memberID).
		UpdateColumn("xp_total", gorm.Expr("xp_total + ?", event.Points))
}

type countMsg struct {
	Type       string `json:"type"`
	ActivityID uint   `json:"activity_id"`
	Present    int64  `json:"present"`
}

func (s *AttendanceService) broadcastCount(ctx context.Context, activityID uint) {
	var count int64
	s.db.WithContext(ctx).Model(&models.ActivityAttendance{}).
		Where("activity_id = ? AND status = 'present'", activityID).Count(&count)
	msg := countMsg{Type: "attendance_count", ActivityID: activityID, Present: count}
	b, _ := json.Marshal(msg)
	s.hub.Broadcast(0, b) // broadcast to all connected clients
}
