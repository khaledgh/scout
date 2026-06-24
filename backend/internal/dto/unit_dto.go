package dto

type CreateUnitRequest struct {
	Name    string `json:"name" validate:"required,min=2"`
	Section string `json:"section" validate:"required,oneof=ashbal kashaf jawala mukashe"`
	Motto   string `json:"motto"`
}

type UpdateUnitRequest struct {
	Name    *string `json:"name"`
	Section *string `json:"section" validate:"omitempty,oneof=ashbal kashaf jawala mukashe"`
	Motto   *string `json:"motto"`
	IsActive *bool  `json:"is_active"`
}

type AddUnitMembersRequest struct {
	MemberIDs []uint `json:"member_ids" validate:"required,min=1"`
}

type AssignUnitLeaderRequest struct {
	UserID     uint   `json:"user_id" validate:"required"`
	RoleInUnit string `json:"role_in_unit" validate:"required,oneof=leader assistant"`
}
