package offline_message

import (
	"github.com/glide-im/glide/pkg/gate"
	"github.com/glide-im/glide/pkg/logger"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/im-service/internal/message_handler"
	"github.com/glide-im/im-service/internal/pkg/db"
	"time"
)

const (
	KeyRedisOfflineMsgPrefix = "im:msg:offline:"
)

func GetHandleFn() func(h *message_handler.MessageHandler, ci *gate.Info, m *messages.GlideMessage) {
	return handler
}

func handler(_ *message_handler.MessageHandler, _ *gate.Info, m *messages.GlideMessage) {
	if m.GetAction() == messages.ActionChatMessage || m.GetAction() == messages.ActionChatMessageResend {
		c := messages.ChatMessage{}
		err := m.Data.Deserialize(&c)
		if err != nil {
			logger.E("deserialize chat message error: %v", err)
			return
		}
		c.To = m.To
		c.From = m.From
		storeOfflineMessage(&c)
	}
}

func storeOfflineMessage(c *messages.ChatMessage) {
	key := KeyRedisOfflineMsgPrefix + c.To
	db.Redis.SAdd(key, c.Mid)
	// TODO 2022-6-22 16:56:57 do not reset expire on new offline message arrived
	// use fixed time segment save offline msg reset segment only.
	db.Redis.Expire(key, time.Hour*24*15)
}
