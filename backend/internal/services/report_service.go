package services

import (
	"context"
	"kashfi/internal/models"
	"time"

	"gorm.io/gorm"
)

type ReportService struct {
	db *gorm.DB
}

func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{db: db}
}

type SectionCount struct {
	Section string `json:"section"`
	Count   int64  `json:"count"`
}

type DashboardData struct {
	MemberCount         int64              `json:"member_count"`
	ActiveMembers       int64              `json:"active_members"`
	AttendanceRate      float64            `json:"attendance_rate"`
	TopUnit             *models.Unit       `json:"top_unit,omitempty"`
	UpcomingActivities  []models.Activity  `json:"upcoming_activities"`
	AtRiskMembers       []models.Member    `json:"at_risk_members"`
	RecentActivities    []models.Activity  `json:"recent_activities"`
	MembersBySection    []SectionCount     `json:"members_by_section"`
	TopMembers          []models.Member    `json:"top_members"`
	RecentBadges        []models.MemberBadge `json:"recent_badges"`
}

func (s *ReportService) Dashboard(ctx context.Context) (*DashboardData, error) {
	d := &DashboardData{}

	s.db.WithContext(ctx).Model(&models.Member{}).Count(&d.MemberCount)
	s.db.WithContext(ctx).Model(&models.Member{}).Where("status = 'active'").Count(&d.ActiveMembers)

	// Attendance rate (last 30 days)
	var total, present int64
	since := time.Now().AddDate(0, 0, -30)
	s.db.WithContext(ctx).Model(&models.ActivityAttendance{}).
		Joins("JOIN activities ON activities.id = activity_attendance.activity_id").
		Where("activities.starts_at >= ?", since).Count(&total)
	s.db.WithContext(ctx).Model(&models.ActivityAttendance{}).
		Joins("JOIN activities ON activities.id = activity_attendance.activity_id").
		Where("activities.starts_at >= ? AND activity_attendance.status = 'present'", since).Count(&present)
	if total > 0 {
		d.AttendanceRate = float64(present) / float64(total) * 100
	}

	// Top unit
	s.db.WithContext(ctx).Order("score_total DESC").Where("is_active = true").First(&d.TopUnit)

	// Upcoming activities (next 7 days)
	s.db.WithContext(ctx).
		Where("starts_at BETWEEN ? AND ? AND status IN ('planned','ongoing')", time.Now(), time.Now().AddDate(0, 0, 7)).
		Preload("ResponsibleUser").Order("starts_at").Limit(5).Find(&d.UpcomingActivities)

	// At-risk: attendance rate < 50% in last 30 days with at least 4 activities
	var atRiskIDs []uint
	s.db.WithContext(ctx).Raw(`
		SELECT aa.member_id
		FROM activity_attendance aa
		JOIN activities a ON a.id = aa.activity_id
		WHERE a.starts_at >= ?
		GROUP BY aa.member_id
		HAVING COUNT(*) >= 4 AND SUM(aa.status = 'present') / COUNT(*) < 0.5
	`, since).Scan(&atRiskIDs)
	if len(atRiskIDs) > 0 {
		s.db.WithContext(ctx).Where("id IN ?", atRiskIDs).Limit(10).Find(&d.AtRiskMembers)
	}

	// Recent activities
	s.db.WithContext(ctx).
		Where("starts_at <= ? AND status = 'completed'", time.Now()).
		Order("starts_at DESC").Limit(5).Find(&d.RecentActivities)

	// Members grouped by section
	s.db.WithContext(ctx).Model(&models.Member{}).
		Select("section, COUNT(*) AS count").Group("section").Scan(&d.MembersBySection)

	// Top members by XP
	s.db.WithContext(ctx).Order("xp_total DESC").Limit(5).Find(&d.TopMembers)

	// Recent badge awards
	s.db.WithContext(ctx).
		Preload("Badge").Preload("Member").
		Order("awarded_at DESC").Limit(6).Find(&d.RecentBadges)

	return d, nil
}

type MemberReport struct {
	Member        *models.Member       `json:"member"`
	AttendanceRate float64             `json:"attendance_rate"`
	TotalActivities int64              `json:"total_activities"`
	TotalPresent    int64              `json:"total_present"`
	Evaluations    []models.Evaluation `json:"evaluations"`
	XPEvents       []models.XPEvent    `json:"xp_events"`
}

func (s *ReportService) MemberReport(ctx context.Context, memberID uint) (*MemberReport, error) {
	var member models.Member
	if err := s.db.WithContext(ctx).Preload("Badges.Badge").First(&member, memberID).Error; err != nil {
		return nil, ErrNotFound
	}
	r := &MemberReport{Member: &member}

	s.db.WithContext(ctx).Model(&models.ActivityAttendance{}).Where("member_id = ?", memberID).Count(&r.TotalActivities)
	s.db.WithContext(ctx).Model(&models.ActivityAttendance{}).Where("member_id = ? AND status = 'present'", memberID).Count(&r.TotalPresent)
	if r.TotalActivities > 0 {
		r.AttendanceRate = float64(r.TotalPresent) / float64(r.TotalActivities) * 100
	}
	s.db.WithContext(ctx).Where("member_id = ?", memberID).Preload("Evaluator").Find(&r.Evaluations)
	s.db.WithContext(ctx).Where("member_id = ?", memberID).Order("created_at DESC").Limit(20).Find(&r.XPEvents)
	return r, nil
}

type UnitReport struct {
	Unit            *models.Unit    `json:"unit"`
	MemberCount     int64           `json:"member_count"`
	AttendanceRate  float64         `json:"attendance_rate"`
	TopMembers      []models.Member `json:"top_members"`
}

func (s *ReportService) UnitReport(ctx context.Context, unitID uint) (*UnitReport, error) {
	var unit models.Unit
	if err := s.db.WithContext(ctx).First(&unit, unitID).Error; err != nil {
		return nil, ErrNotFound
	}
	r := &UnitReport{Unit: &unit}

	s.db.WithContext(ctx).Model(&models.UnitMember{}).Where("unit_id = ?", unitID).Count(&r.MemberCount)

	var memberIDs []uint
	s.db.WithContext(ctx).Model(&models.UnitMember{}).Where("unit_id = ?", unitID).Pluck("member_id", &memberIDs)
	if len(memberIDs) > 0 {
		var total, present int64
		s.db.WithContext(ctx).Model(&models.ActivityAttendance{}).Where("member_id IN ?", memberIDs).Count(&total)
		s.db.WithContext(ctx).Model(&models.ActivityAttendance{}).Where("member_id IN ? AND status = 'present'", memberIDs).Count(&present)
		if total > 0 {
			r.AttendanceRate = float64(present) / float64(total) * 100
		}
		s.db.WithContext(ctx).Where("id IN ?", memberIDs).Order("xp_total DESC").Limit(5).Find(&r.TopMembers)
	}
	return r, nil
}

type MonthlyReport struct {
	Month           string  `json:"month"`
	TotalActivities int64   `json:"total_activities"`
	AttendanceRate  float64 `json:"attendance_rate"`
	NewMembers      int64   `json:"new_members"`
	XPDistributed   int64   `json:"xp_distributed"`
}

func (s *ReportService) Monthly(ctx context.Context, year, month int) (*MonthlyReport, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	r := &MonthlyReport{Month: start.Format("2006-01")}
	s.db.WithContext(ctx).Model(&models.Activity{}).
		Where("starts_at >= ? AND starts_at < ?", start, end).Count(&r.TotalActivities)

	var total, present int64
	s.db.WithContext(ctx).Model(&models.ActivityAttendance{}).
		Joins("JOIN activities ON activities.id = activity_attendance.activity_id").
		Where("activities.starts_at >= ? AND activities.starts_at < ?", start, end).Count(&total)
	s.db.WithContext(ctx).Model(&models.ActivityAttendance{}).
		Joins("JOIN activities ON activities.id = activity_attendance.activity_id").
		Where("activities.starts_at >= ? AND activities.starts_at < ? AND activity_attendance.status = 'present'", start, end).Count(&present)
	if total > 0 {
		r.AttendanceRate = float64(present) / float64(total) * 100
	}

	s.db.WithContext(ctx).Model(&models.Member{}).
		Where("join_date >= ? AND join_date < ?", start, end).Count(&r.NewMembers)

	var totalXP struct{ Sum int64 }
	s.db.WithContext(ctx).Model(&models.XPEvent{}).
		Where("created_at >= ? AND created_at < ?", start, end).
		Select("COALESCE(SUM(points),0) AS sum").Scan(&totalXP)
	r.XPDistributed = totalXP.Sum

	return r, nil
}
