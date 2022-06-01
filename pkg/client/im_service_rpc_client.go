package client

import (
	"github.com/glide-im/glide/pkg/gate"
	"github.com/glide-im/glide/pkg/messages"
)

type IMRpcClient struct {
}

func (I *IMRpcClient) SetClientID(old gate.ID, new_ gate.ID) error {

	return nil
}

func (I *IMRpcClient) ExitClient(id gate.ID) error {
	return nil
}

func (I *IMRpcClient) IsOnline(id gate.ID) bool {
	return true
}

func (I *IMRpcClient) EnqueueMessage(id gate.ID, message *messages.GlideMessage) error {
	return nil
}
