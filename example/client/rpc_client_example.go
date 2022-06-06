package main

import (
	"github.com/glide-im/glide/pkg/gate"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/im-service/pkg/client"
	"github.com/glide-im/im-service/pkg/rpc"
)

func main() {
	RpcClientExample()
}

func RpcClientExample() {
	options := &rpc.ClientOptions{
		Addr: "127.0.0.1",
		Port: 8092,
		Name: "im_rpc_server",
	}
	cli, err := client.NewIMServiceClient(options)
	defer cli.Close()
	if err != nil {
		panic(err)
	}
	err = cli.EnqueueMessage(gate.NewID2(1), messages.NewEmptyMessage())
	if err != nil {
		panic(err)
	}
}
