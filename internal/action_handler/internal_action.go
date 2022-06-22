package action_handler

import (
	"github.com/glide-im/glide/pkg/gate"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/glide/pkg/messaging"
)

type InternalActionHandler struct {
}

func (o *InternalActionHandler) Handle(h *messaging.MessageInterfaceImpl, cliInfo *gate.Info, m *messages.GlideMessage) bool {
	if m.GetAction().IsInternal() {
		if m.GetAction() == messages.ActionInternalOffline && !cliInfo.ID.IsTemp() {
			// TODO notify friends
		}
		return true
	}
	return false
}
