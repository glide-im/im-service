package action_handler

import (
	"github.com/glide-im/glide/pkg/gate"
	"github.com/glide-im/glide/pkg/logger"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/glide/pkg/messaging"
	"github.com/glide-im/im-service/internal/offline_message"
)

type InternalActionHandler struct {
}

func (o *InternalActionHandler) Handle(h *messaging.MessageInterfaceImpl, cliInfo *gate.Info, m *messages.GlideMessage) bool {
	if m.GetAction().IsInternal() {
		if !cliInfo.ID.IsTemp() {
			switch m.GetAction() {
			case messages.ActionInternalOffline:
			case messages.ActionInternalOnline:
				go func() {
					defer func() {
						err, ok := recover().(error)
						if err != nil && ok {
							logger.ErrE("push offline message error", err)
						}
					}()
					offline_message.PushOfflineMessage(h, cliInfo.ID.UID())
				}()
			}
		}
		return true
	}
	return false
}
