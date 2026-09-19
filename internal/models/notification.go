package models

type NotificationMessage struct {
	ChannelID string `json:"channelId"`
	GameID    string `json:"gameId"`
	TurnID    string `json:"turnID"`
}
