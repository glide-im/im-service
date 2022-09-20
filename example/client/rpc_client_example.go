package main

import (
	"github.com/glide-im/glide/pkg/gate"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/glide/pkg/rpc"
	"github.com/glide-im/glide/pkg/subscription/subscription_impl"
	"github.com/glide-im/im-service/pkg/client"
)

func main() {
	// 消息网关
	//RpcGatewayClientExample()

	// 发布订阅(群聊)
	RpcSubscriberClientExample()
}

func RpcGatewayClientExample() {

	// 消息网关配置
	options := &rpc.ClientOptions{
		Addr: "127.0.0.1",
		Port: 8092,
		Name: "im_rpc_server",
	}

	// 创建消息网关接口客户端
	cli, err := client.NewGatewayRpcImpl(options)
	defer cli.Close()
	if err != nil {
		panic(err)
	}

	// 给网关中指定 id 的链接推送一条消息
	err = cli.EnqueueMessage(gate.NewID2("1"), messages.NewEmptyMessage())
	if err != nil {
		panic(err)
	}

	// 设置网关中连接新 id
	err = cli.SetClientID(gate.NewID2("1"), gate.NewID2("2"))
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
