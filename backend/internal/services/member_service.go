package services

import (
	"context"
	"kashfi/internal/config"
	"kashfi/internal/models"
	"kashfi/internal/utils"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type MemberService struct {
	db     *gorm.DB
	gamify *GamificationService
	cfg    *config.Config
}

func NewMemberService(db *gorm.DB, gamify *GamificationService, cfg *config.Config) *MemberService {
	return &MemberService{db: db, gamify: gamify, cfg: cfg}
}

func (s *MemberService) UploadPath() string  { return s.cfg.Upload.Dir }
func (s *MemberService) UploadMaxMB() int64   { return s.cfg.Upload.MaxSizeMB }
func (s *MemberService) UploadTypes() string  { return s.cfg.Upload.AllowedTypes }

// SetPhoto stores the public URL for an uploaded member photo.
func (s *MemberService) SetPhoto(ctx context.Context, memberID uint, filename string) (*models.Member, error) {
	var m models.Member
	if err := s.db.WithContext(ctx).First(&m, memberID).Error; err != nil {
		return nil, ErrNotFound
	}
	url := filepath.ToSlash(filepath.Join(s.cfg.Upload.PublicPath, filename))
	if err := s.db.WithContext(ctx).Model(&m).Update("photo_url", url).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

type MemberFilter struct {
	UnitID  *uint
	Section *string
	Status  *string
	Search  *string
}

func (s *MemberService) List(ctx context.Context, f MemberFilter, page, pageSize int) ([]*models.Member, int64, error) {
	q := s.db.WithContext(ctx).Model(&models.Member{})
	if f.UnitID != nil {
		q = q.Joins("JOIN unit_members um ON um.member_id = members.id AND um.deleted_at IS NULL").
			Where("um.unit_id = ?", *f.UnitID)
	}
	if f.Section != nil {
		q = q.Where("members.section = ?", *f.Section)
	}
	if f.Status != nil {
		q = q.Where("members.status = ?", *f.Status)
	}
	if f.Search != nil {
		like := "%" + *f.Search + "%"
		q = q.Where("members.full_name LIKE ? OR members.parent_phone LIKE ?", like, like)
	}

	var total int64
	q.Count(&total)

	var members []*models.Member
	err := q.Offset((page - 1) * pageSize).Limit(pageSize).Find(&members).Error
	return members, total, err
}

func (s *MemberService) Get(ctx context.Context, id uint) (*models.Member, error) {
	var m models.Member
	err := s.db.WithContext(ctx).
		Preload("Badges.Badge").
		Preload("Skills.Skill").
		Preload("Units.Unit").
		First(&m, id).Error
	if err != nil {
		return nil, ErrNotFound
	}
	return &m, nil
}

type CreateMemberInput struct {
	FullName    string
	BirthDate   time.Time
	Gender      string
	Section     models.Section
	RankStage   string
	JoinDate    time.Time
	ParentName  string
	ParentPhone string
	SecondPhone string
	Address     string
	UserID      *uint
}

func (s *MemberService) Create(ctx context.Context, in CreateMemberInput) (*models.Member, error) {
	m := &models.Member{
		FullName:    in.FullName,
		BirthDate:   in.BirthDate,
		Gender:      in.Gender,
		Section:     in.Section,
		RankStage:   in.RankStage,
		JoinDate:    in.JoinDate,
		ParentName:  in.ParentName,
		ParentPhone: utils.NormalizePhone(in.ParentPhone),
		SecondPhone: in.SecondPhone,
		Address:     in.Address,
		UserID:      in.UserID,
		Status:      models.MemberStatusActive,
		Level:       1,
	}
	if err := s.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

type UpdateMemberInput struct {
	FullName    *string
	BirthDate   *time.Time
	Gender      *string
	Section     *models.Section
	RankStage   *string
	JoinDate    *time.Time
	ParentName  *string
	ParentPhone *string
	SecondPhone *string
	Address     *string
	Status      *models.MemberStatus
}

func (s *MemberService) Update(ctx context.Context, id uint, in UpdateMemberInput) (*models.Member, error) {
	var m models.Member
	if err := s.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, ErrNotFound
	}
	updates := map[string]interface{}{}
	if in.FullName != nil { updates["full_name"] = *in.FullName }
	if in.BirthDate != nil { updates["birth_date"] = *in.BirthDate }
	if in.Gender != nil { updates["gender"] = *in.Gender }
	if in.Section != nil { updates["section"] = *in.Section }
	if in.RankStage != nil { updates["rank_stage"] = *in.RankStage }
	if in.JoinDate != nil { updates["join_date"] = *in.JoinDate }
	if in.ParentName != nil { updates["parent_name"] = *in.ParentName }
	if in.ParentPhone != nil { updates["parent_phone"] = utils.NormalizePhone(*in.ParentPhone) }
	if in.SecondPhone != nil { updates["secondary_phone"] = *in.SecondPhone }
	if in.Address != nil { updates["address"] = *in.Address }
	if in.Status != nil { updates["status"] = *in.Status }

	if err := s.db.WithContext(ctx).Model(&m).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// CanManageMember reports whether the user may create/update/delete the given
// member. super_admin can manage anyone; a leader/assistant can manage a member
// only if that member belongs to a unit the user is assigned to lead.
func (s *MemberService) CanManageMember(ctx context.Context, userID uint, role string, memberID uint) bool {
	if role == "super_admin" {
		return true
	}
	if role != "leader" && role != "assistant" {
		return false
	}
	var count int64
	s.db.WithContext(ctx).
		Table("unit_members AS um").
		Joins("JOIN unit_leaders ul ON ul.unit_id = um.unit_id").
		Where("um.member_id = ? AND um.deleted_at IS NULL AND ul.user_id = ?", memberID, userID).
		Count(&count)
	return count > 0
}

func (s *MemberService) Delete(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&models.Member{}, id)
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}

func (s *MemberService) GetMedical(ctx context.Context, memberID uint) (*models.MemberMedical, error) {
	var med models.MemberMedical
	err := s.db.WithContext(ctx).Where("member_id = ?", memberID).First(&med).Error
	if err != nil {
		return nil, ErrNotFound
	}
	return &med, nil
}

type MedicalInput struct {
	BloodType         string
	Allergies         string
	ChronicConditions string
	Medications       string
	EmergencyNotes    string
}

func (s *MemberService) UpsertMedical(ctx context.Context, memberID uint, in MedicalInput) (*models.MemberMedical, error) {
	med := models.MemberMedical{
		MemberID:          memberID,
		BloodType:         in.BloodType,
		Allergies:         in.Allergies,
		ChronicConditions: in.ChronicConditions,
		Medications:       in.Medications,
		EmergencyNotes:    in.EmergencyNotes,
	}
	err := s.db.WithContext(ctx).Where("member_id = ?", memberID).Assign(med).FirstOrCreate(&med).Error
	if err != nil {
		return nil, err
	}
	return &med, nil
}

type EvalInput struct {
	Period        string
	Discipline    int
	Participation int
	Leadership    int
	Skill         int
	Overall       int
	Notes         string
}

func (s *MemberService) CreateEvaluation(ctx context.Context, memberID, evaluatorID uint, in EvalInput) (*models.Evaluation, error) {
	eval := &models.Evaluation{
		MemberID:      memberID,
		EvaluatorID:   evaluatorID,
		Period:        in.Period,
		Discipline:    in.Discipline,
		Participation: in.Participation,
		Leadership:    in.Leadership,
		Skill:         in.Skill,
		Overall:       in.Overall,
		Notes:         in.Notes,
	}
	if err := s.db.WithContext(ctx).Create(eval).Error; err != nil {
		return nil, err
	}
	return eval, nil
}

type Timeline struct {
	Attendances []models.ActivityAttendance `json:"attendances"`
	Badges      []models.MemberBadge        `json:"badges"`
	XPEvents    []models.XPEvent            `json:"xp_events"`
	Evaluations []models.Evaluation         `json:"evaluations"`
}

func (s *MemberService) GetTimeline(ctx context.Context, memberID uint) (*Timeline, error) {
	tl := &Timeline{}
	s.db.WithContext(ctx).Where("member_id = ?", memberID).
		Preload("Activity").Order("created_at DESC").Limit(20).Find(&tl.Attendances)
	s.db.WithContext(ctx).Where("member_id = ?", memberID).
		Preload("Badge").Order("awarded_at DESC").Find(&tl.Badges)
	s.db.WithContext(ctx).Where("member_id = ?", memberID).
		Order("created_at DESC").Limit(50).Find(&tl.XPEvents)
	s.db.WithContext(ctx).Where("member_id = ?", memberID).
		Preload("Evaluator").Order("created_at DESC").Find(&tl.Evaluations)
	return tl, nil
}

func (s *MemberService) GetQRToken(ctx context.Context, memberID uint, secret string) (string, error) {
	var m models.Member
	if err := s.db.WithContext(ctx).First(&m, memberID).Error; err != nil {
		return "", ErrNotFound
	}
	return utils.GenerateQRToken(memberID, secret), nil
}
