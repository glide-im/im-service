package message_handler

import (
	"github.com/glide-im/glide/pkg/gate"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/im-service/internal/offline_message"
)

func (d *MessageHandler) handleHeartbeat(cInfo *gate.Info, msg *messages.GlideMessage) error {
	return nil
}

// handleAckRequest 处理接收者收到消息发回来的确认消息
func (d *MessageHandler) handleAckRequest(c *gate.Info, msg *messages.GlideMessage) error {
	ackMsg := new(messages.AckRequest)
	if !d.unmarshalData(c, msg, ackMsg) {
		return nil
	}
	ackNotify := messages.NewMessage(0, messages.ActionAckNotify, ackMsg)

	// 通知发送者, 对方已收到消息
	d.dispatchAllDevice(ackMsg.From, ackNotify)
	return nil
}

func (d *MessageHandler) handleAckOffline(c *gate.Info, msg *messages.GlideMessage) error {
	if c.ID.IsTemp() {
		return nil
	}
	ackMsg := new(messages.AckRequest)
	if !d.unmarshalData(c, msg, ackMsg) {
		return nil
	}
	offline_message.AckOfflineMessage(c.ID.UID())
	return nil
}
