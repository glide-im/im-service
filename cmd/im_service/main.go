package main

import (
	"github.com/glide-im/glide/pkg/auth/jwt_auth"
	"github.com/glide-im/glide/pkg/bootstrap"
	"github.com/glide-im/glide/pkg/logger"
	"github.com/glide-im/glide/pkg/messaging/message_handler"
	"github.com/glide-im/glide/pkg/rpc"
	"github.com/glide-im/glide/pkg/store"
	"github.com/glide-im/glide/pkg/subscription/subscription_impl"
	"github.com/glide-im/im-service/internal/config"
	"github.com/glide-im/im-service/internal/im_server"
	"github.com/glide-im/im-service/internal/message_store_db"
)

func main() {

	config.MustLoad()

	gateway, err := im_server.NewServer(config.WsServer.ID, config.WsServer.Addr, config.WsServer.Port)
	if err != nil {
		panic(err)
	}

	auth := jwt_auth.NewAuthorizeImpl(config.WsServer.JwtSecret)

	var cStore store.MessageStore = &message_store_db.IdleChatMessageStore{}
	var sStore store.SubscriptionStore = &message_store_db.IdleSubscriptionStore{}
	var seqStore subscription_impl.ChannelSequenceStore = &message_store_db.IdleSubscriptionStore{}

	if config.Common.StoreMessageHistory {
		dbStore, err := message_store_db.New(config.MySql)
		if err != nil {
			panic(err)
		}
		cStore = dbStore
		sStore = &message_store_db.SubscriptionMessageStore{}
		seqStore = &message_store_db.SubscriptionMessageStore{}
	} else {
		logger.D("Common.StoreMessageHistory is false, message history will not be stored")
	}

	handler, err := message_handler.NewHandler(cStore, auth)
	if err != nil {
		panic(err)
	}

	options := bootstrap.Options{
		Messaging:    handler,
		Gate:         gateway,
		Subscription: subscription_impl.NewSubscription(sStore, seqStore),
	}

	go func() {
		logger.D("websocket listening on %s:%d", config.WsServer.Addr, config.WsServer.Port)
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
	logger.D("rpc listening on %s:%d", rpcOpts.Addr, rpcOpts.Port)
	err = im_server.RunRpcServer(&rpcOpts, gateway)
	if err != nil {
		panic(err)
	}
}
