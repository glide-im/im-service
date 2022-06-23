package main

import (
	"github.com/glide-im/glide/pkg/gate"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/glide/pkg/rpc"
	"github.com/glide-im/glide/pkg/subscription/subscription_impl"
	"github.com/glide-im/im-service/pkg/client"
)

func main() {
	//RpcGatewayClientExample()
	RpcSubscriberClientExample()
}

func RpcGatewayClientExample() {
	options := &rpc.ClientOptions{
		Addr: "127.0.0.1",
		Port: 8092,
		Name: "im_rpc_server",
	}
	cli, err := client.NewGatewayRpcImpl(options)
	defer cli.Close()
	if err != nil {
		panic(err)
	}
	err = cli.EnqueueMessage(gate.NewID2("1"), messages.NewEmptyMessage())
	if err != nil {
		panic(err)
	}
}

func RpcSubscriberClientExample() {
	options := &rpc.ClientOptions{
		Addr: "127.0.0.1",
		Port: 8092,
		Name: "im_rpc_server",
	}
	cli, err := client.NewSubscriptionRpcImpl(options)
	defer cli.Close()
	if err != nil {
		panic(err)
	}

	//err = cli.CreateChannel("1", &subscription.ChanInfo{
	//	ID:   "1",
	//	Type: 0,
	//})
	//if err != nil {
	//	panic(err)
	//}

	err = cli.Subscribe("1", "1", &subscription_impl.SubscriberOptions{
		Perm: subscription_impl.PermRead | subscription_impl.PermWrite,
	})
	if err != nil {
		panic(err)
	}

	msg := &subscription_impl.PublishMessage{
		From:    "1",
		Seq:     1,
		Type:    subscription_impl.TypeMessage,
		Message: messages.NewMessage(0, "1", &messages.ChatMessage{}),
	}
	err = cli.Publish("1", msg)
	if err != nil {
		panic(err)
	}
}
