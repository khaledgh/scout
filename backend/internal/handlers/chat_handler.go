package handlers

import (
	"encoding/json"
	"kashfi/internal/config"
	appMiddleware "kashfi/internal/middleware"
	"kashfi/internal/services"
	"kashfi/internal/utils"
	"kashfi/internal/ws"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:    func(r *http.Request) bool { return true },
}

type ChatHandler struct {
	svc *services.ChatService
	hub *ws.Hub
	cfg *config.Config
}

func NewChatHandler(svc *services.ChatService, hub *ws.Hub, cfg *config.Config) *ChatHandler {
	return &ChatHandler{svc: svc, hub: hub, cfg: cfg}
}

func (h *ChatHandler) Channels(c echo.Context) error {
	userID := appMiddleware.GetUserID(c)
	channels, err := h.svc.UserChannels(c.Request().Context(), userID)
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OK(c, channels)
}

func (h *ChatHandler) Messages(c echo.Context) error {
	channelID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return utils.BadRequest(c, "invalid channel id")
	}
	p := utils.ParsePagination(c)
	msgs, total, err := h.svc.Messages(c.Request().Context(), uint(channelID), p.Page, p.PageSize)
	if err != nil {
		return utils.Internal(c)
	}
	return utils.OKWithMeta(c, msgs, utils.BuildMeta(p, total))
}

type incomingMsg struct {
	ChannelID uint   `json:"channel_id"`
	Body      string `json:"body"`
}

func (h *ChatHandler) WebSocket(c echo.Context) error {
	conn, err := wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	userID := appMiddleware.GetUserID(c)
	send := make(chan []byte, 64)

	go func() {
		defer conn.Close()
		for data := range send {
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				break
			}
		}
	}()

	defer close(send)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var in incomingMsg
		if err := json.Unmarshal(msg, &in); err != nil || in.Body == "" {
			continue
		}
		saved, err := h.svc.SaveMessage(c.Request().Context(), in.ChannelID, userID, in.Body)
		if err != nil {
			continue
		}
		out, _ := json.Marshal(saved)
		h.hub.Broadcast(in.ChannelID, out)
	}
	return nil
}
