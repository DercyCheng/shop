package entity

import "time"

// Message represents a user's message or inquiry
type Message struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	UserID      int64     `json:"user_id" bson:"user_id"`
	MessageType int       `json:"message_type" bson:"message_type"`
	Subject     string    `json:"subject" bson:"subject"`
	Content     string    `json:"content" bson:"content"`
	File        string    `json:"file" bson:"file"`
	Images      []string  `json:"images" bson:"images"`
	Status      int       `json:"status" bson:"status"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// MessageTypeMap maps message type constants to their descriptions
var MessageTypeMap = map[int]string{
	1: "留言",
	2: "投诉",
	3: "询问",
	4: "售后",
	5: "求购",
}

// MessageStatusMap maps message status constants to their descriptions
var MessageStatusMap = map[int]string{
	0: "未处理",
	1: "已处理",
}

// GetMessageTypeText returns the text description of the message type
func (m *Message) GetMessageTypeText() string {
	if text, ok := MessageTypeMap[m.MessageType]; ok {
		return text
	}
	return "未知类型"
}

// GetStatusText returns the text description of the message status
func (m *Message) GetStatusText() string {
	if text, ok := MessageStatusMap[m.Status]; ok {
		return text
	}
	return "未知状态"
}
