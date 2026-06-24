package dto

type CreateAnnouncementRequest struct {
	Title    string `json:"title" validate:"required"`
	Body     string `json:"body" validate:"required"`
	Audience string `json:"audience" validate:"required,oneof=all unit leaders"`
	UnitID   *uint  `json:"unit_id"`
	Pinned   bool   `json:"pinned"`
}

type UpdateAnnouncementRequest struct {
	Title    *string `json:"title"`
	Body     *string `json:"body"`
	Audience *string `json:"audience"`
	UnitID   *uint   `json:"unit_id"`
	Pinned   *bool   `json:"pinned"`
}

type SendMessageRequest struct {
	ChannelID uint   `json:"channel_id" validate:"required"`
	Body      string `json:"body" validate:"required"`
}
