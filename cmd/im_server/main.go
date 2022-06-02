package main

import (
	"github.com/glide-im/glide/pkg/auth/jwt_auth"
	"github.com/glide-im/glide/pkg/bootstrap"
	"github.com/glide-im/glide/pkg/messaging/message_handler"
	"github.com/glide-im/glide/pkg/subscription/group_subscription"
	"github.com/glide-im/im-service/internal/config"
	"github.com/glide-im/im-service/internal/im_server"
	"github.com/glide-im/im-service/internal/message_store_db"
	"github.com/glide-im/im-service/pkg/rpc"
)

func main() {

	config.MustLoad()

	gateway, err := im_server.NewServer(config.WsServer.ID, config.WsServer.Addr, config.WsServer.Port)
	if err != nil {
		panic(err)
	}

	auth := jwt_auth.NewAuthorizeImpl(config.WsServer.JwtSecret)
	dbStore, err := message_store_db.New(config.MySql)
	if err != nil {
		panic(err)
	}
	handler, err := message_handler.NewHandler(dbStore, auth)
	if err != nil {
		panic(err)
	}

	store := &message_store_db.SubscriptionMessageStore{}
	options := bootstrap.Options{
		Messaging:    handler,
		Gate:         gateway,
		Subscription: group_subscription.NewSubscription(store),
	}

	go func() {
		err = bootstrap.Bootstrap(&options)

		if err != nil {
			panic(err)
		}
	}()

	rpcOpts := rpc.ServerOptions{
		Name:    config.IMService.Name,
		Network: config.IMService.Network,
		Addr:    config.IMService.Addr,
		Port:    config.IMService.Port,
	}
	err = im_server.RunRpcServer(&rpcOpts, gateway)
	if err != nil {
		panic(err)
	}
}
