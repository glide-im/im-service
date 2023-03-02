package action_handler

import (
	"github.com/glide-im/glide/pkg/gate"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/glide/pkg/messaging"
	messages2 "github.com/glide-im/im-service/pkg/messages"
)

type ClientCustomMessageHandler struct {
}

func (c *ClientCustomMessageHandler) Handle(h *messaging.MessageInterfaceImpl, ci *gate.Info, m *messages.GlideMessage) bool {
	if m.Action != messages2.ActionClientCustom {
		return false
	}
	dispatch2AllDevice(h, m.To, m)
	return true
}
