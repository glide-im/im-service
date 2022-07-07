package action_handler

import (
	"github.com/glide-im/glide/pkg/gate"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/glide/pkg/messaging"
)

const (
	ActionClientCustom = "message.cli"
)

// ClientCustom client custom message, server does not store to database.
type ClientCustom struct {
	Type    string      `json:"type,omitempty"`
	Content interface{} `json:"content,omitempty"`
}

type ClientCustomMessageHandler struct {
}

func (c *ClientCustomMessageHandler) Handle(h *messaging.MessageInterfaceImpl, ci *gate.Info, m *messages.GlideMessage) bool {
	if m.Action != ActionClientCustom {
		return false
	}
	dispatch2AllDevice(h, m.To, m)
	return true
}
