package message_handler

import (
	"github.com/glide-im/glide/pkg/gate"
	"github.com/glide-im/glide/pkg/messages"
)

const (
	NotifyKickOut messages.Action = "notify.kickout"
	AckOffline                    = "ack.offline"
)

func createKickOutMessage(c *gate.Info) *messages.GlideMessage {
	return messages.NewMessage(0, NotifyKickOut, "")
}
