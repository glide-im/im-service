package world_channel

import (
	"encoding/json"
	"github.com/glide-im/glide/pkg/gate"
	"github.com/glide-im/glide/pkg/logger"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/glide/pkg/subscription"
	"github.com/glide-im/glide/pkg/subscription/subscription_impl"
	"github.com/glide-im/im-service/internal/message_handler"
	"time"
)

var sub subscription_impl.SubscribeWrap
var chanId = subscription.ChanID("the_world_channel")

func EnableWorldChannel(subscribe subscription_impl.SubscribeWrap) error {
	sub = subscribe
	err := sub.CreateChannel(chanId, &subscription.ChanInfo{
		ID: chanId,
	})
	if err != nil {
		return err
	}
	err = sub.Subscribe(chanId, "system", &subscription_impl.SubscriberOptions{Perm: subscription_impl.PermWrite})
	return err
}

func OnUserOnline(id gate.ID) {
	if id.IsTemp() {
		return
	}
	err := sub.Subscribe(chanId, subscription.SubscriberID(id.UID()),
		&subscription_impl.SubscriberOptions{Perm: subscription_impl.PermRead | subscription_impl.PermWrite})
	if err == nil {

		b, _ := json.Marshal(&messages.ChatMessage{
			Mid:     time.Now().UnixNano(),
			Seq:     0,
			From:    "system",
			To:      string(chanId),
			Type:    100,
			Content: id.UID(),
			SendAt:  time.Now().Unix(),
		})
		_ = sub.Publish(chanId, &subscription_impl.PublishMessage{
			From:    "system",
			Type:    subscription_impl.TypeMessage,
			Message: messages.NewMessage(0, message_handler.ActionGroupMessage, b),
		})
	} else {
		logger.E("$v", err)
	}
}

func OnUserOffline(id gate.ID) {
	if id.IsTemp() {
		return
	}
	_ = sub.UnSubscribe(chanId, subscription.SubscriberID(id.UID()))
	b, _ := json.Marshal(&messages.ChatMessage{
		Mid:     time.Now().UnixNano(),
		Seq:     0,
		From:    "system",
		To:      string(chanId),
		Type:    101,
		Content: id.UID(),
		SendAt:  time.Now().Unix(),
	})
	_ = sub.Publish(chanId, &subscription_impl.PublishMessage{
		From:    "system",
		Type:    subscription_impl.TypeMessage,
		Message: messages.NewMessage(0, message_handler.ActionGroupMessage, b),
	})
}
