package services

import (
	"context"
	"kashfi/internal/models"
	"kashfi/internal/ws"

	"gorm.io/gorm"
)

type ChatService struct {
	db  *gorm.DB
	hub *ws.Hub
}

func NewChatService(db *gorm.DB, hub *ws.Hub) *ChatService {
	return &ChatService{db: db, hub: hub}
}

func (s *ChatService) UserChannels(ctx context.Context, userID uint) ([]models.Channel, error) {
	// Return unit channels for units the user leads + direct channels
	var channels []models.Channel
	s.db.WithContext(ctx).Where("type = 'unit'").Preload("Unit").Find(&channels)
	return channels, nil
}

func (s *ChatService) Messages(ctx context.Context, channelID uint, page, pageSize int) ([]models.Message, int64, error) {
	var msgs []models.Message
	var total int64
	s.db.WithContext(ctx).Model(&models.Message{}).Where("channel_id = ?", channelID).Count(&total)
	err := s.db.WithContext(ctx).Where("channel_id = ?", channelID).
		Preload("Sender").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&msgs).Error
	return msgs, total, err
}

func (s *ChatService) EnsureUnitChannels(ctx context.Context) error {
	var units []models.Unit
	s.db.WithContext(ctx).Where("is_active = true").Find(&units)
	for _, u := range units {
		ch := models.Channel{
			Type:   models.ChannelUnit,
			UnitID: &u.ID,
			Name:   u.Name,
		}
		s.db.WithContext(ctx).Where("type = 'unit' AND unit_id = ?", u.ID).FirstOrCreate(&ch)
	}
	return nil
}

func (s *ChatService) SaveMessage(ctx context.Context, channelID, senderID uint, body string) (*models.Message, error) {
	msg := &models.Message{
		ChannelID: channelID,
		SenderID:  senderID,
		Body:      body,
	}
	if err := s.db.WithContext(ctx).Create(msg).Error; err != nil {
		return nil, err
	}
	return msg, nil
}
