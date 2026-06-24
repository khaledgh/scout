package dto

type CreateActivityRequest struct {
	Title             string   `json:"title" validate:"required,min=2"`
	Description       string   `json:"description"`
	Type              string   `json:"type" validate:"required,oneof=camp hike training meeting service"`
	Location          string   `json:"location"`
	LocationLat       *float64 `json:"location_lat"`
	LocationLng       *float64 `json:"location_lng"`
	StartsAt          string   `json:"starts_at" validate:"required"`
	EndsAt            string   `json:"ends_at" validate:"required"`
	UnitID            *uint    `json:"unit_id"`
}

type UpdateActivityRequest struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	Type        *string  `json:"type" validate:"omitempty,oneof=camp hike training meeting service"`
	Location    *string  `json:"location"`
	LocationLat *float64 `json:"location_lat"`
	LocationLng *float64 `json:"location_lng"`
	StartsAt    *string  `json:"starts_at"`
	EndsAt      *string  `json:"ends_at"`
	Status      *string  `json:"status" validate:"omitempty,oneof=planned ongoing completed cancelled"`
	UnitID      *uint    `json:"unit_id"`
}

type RecordAttendanceRequest struct {
	Records []AttendanceRecord `json:"records" validate:"required,min=1"`
}

type AttendanceRecord struct {
	MemberID uint   `json:"member_id" validate:"required"`
	Status   string `json:"status" validate:"required,oneof=present absent excused late"`
}

type CheckInRequest struct {
	Method   string   `json:"method" validate:"required,oneof=qr gps"`
	QRToken  *string  `json:"qr_token"`
	Lat      *float64 `json:"lat"`
	Lng      *float64 `json:"lng"`
}

type FeedbackRequest struct {
	Rating        int    `json:"rating" validate:"required,min=1,max=5"`
	WhatWentWell  string `json:"what_went_well"`
	WhatToImprove string `json:"what_to_improve"`
	Comment       string `json:"comment"`
}

type ActivityFilter struct {
	Type   *string `query:"type"`
	From   *string `query:"from"`
	To     *string `query:"to"`
	UnitID *uint   `query:"unit_id"`
	Status *string `query:"status"`
}
