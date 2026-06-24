package dto

type CreateMemberRequest struct {
	FullName    string `json:"full_name" validate:"required,min=2"`
	BirthDate   string `json:"birth_date" validate:"required"`
	Gender      string `json:"gender" validate:"required,oneof=male female"`
	Section     string `json:"section" validate:"required,oneof=ashbal kashaf jawala mukashe"`
	RankStage   string `json:"rank_stage"`
	JoinDate    string `json:"join_date" validate:"required"`
	ParentName  string `json:"parent_name"`
	ParentPhone string `json:"parent_phone"`
	SecondPhone string `json:"secondary_phone"`
	Address     string `json:"address"`
	UserID      *uint  `json:"user_id"`
}

type UpdateMemberRequest struct {
	FullName    *string `json:"full_name"`
	BirthDate   *string `json:"birth_date"`
	Gender      *string `json:"gender" validate:"omitempty,oneof=male female"`
	Section     *string `json:"section" validate:"omitempty,oneof=ashbal kashaf jawala mukashe"`
	RankStage   *string `json:"rank_stage"`
	JoinDate    *string `json:"join_date"`
	ParentName  *string `json:"parent_name"`
	ParentPhone *string `json:"parent_phone"`
	SecondPhone *string `json:"secondary_phone"`
	Address     *string `json:"address"`
	Status      *string `json:"status" validate:"omitempty,oneof=active inactive"`
}

type UpsertMedicalRequest struct {
	BloodType         string `json:"blood_type"`
	Allergies         string `json:"allergies"`
	ChronicConditions string `json:"chronic_conditions"`
	Medications       string `json:"medications"`
	EmergencyNotes    string `json:"emergency_notes"`
}

type CreateEvaluationRequest struct {
	Period        string `json:"period" validate:"required"`
	Discipline    int    `json:"discipline" validate:"min=0,max=10"`
	Participation int    `json:"participation" validate:"min=0,max=10"`
	Leadership    int    `json:"leadership" validate:"min=0,max=10"`
	Skill         int    `json:"skill" validate:"min=0,max=10"`
	Overall       int    `json:"overall" validate:"min=0,max=10"`
	Notes         string `json:"notes"`
}

type MemberFilter struct {
	UnitID  *uint   `query:"unit_id"`
	Section *string `query:"section"`
	Status  *string `query:"status"`
	Search  *string `query:"search"`
}
