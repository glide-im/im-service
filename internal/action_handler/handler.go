package action_handler

import (
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/glide/pkg/messaging"
)

func Setup(handler messaging.Interface) {
	handler.AddHandler(&InternalActionHandler{})
	handler.AddHandler(messaging.NewActionHandler(messages.ActionHeartbeat, handleHeartbeat))
}
