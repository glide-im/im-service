package main

import (
	"flag"
	"github.com/glide-im/glide/pkg/auth/jwt_auth"
	"github.com/glide-im/glide/pkg/bootstrap"
	"github.com/glide-im/glide/pkg/logger"
	"github.com/glide-im/glide/pkg/rpc"
	"github.com/glide-im/glide/pkg/store"
	"github.com/glide-im/glide/pkg/subscription/subscription_impl"
	"github.com/glide-im/im-service/internal/action_handler"
	"github.com/glide-im/im-service/internal/config"
	"github.com/glide-im/im-service/internal/im_server"
	"github.com/glide-im/im-service/internal/message_handler"
	"github.com/glide-im/im-service/internal/message_store_db"
	"github.com/glide-im/im-service/internal/pkg/db"
	"github.com/glide-im/im-service/internal/server_state"
)

var state *string

func init() {
	state = flag.String("state", "", "show im server run state")
	flag.Parse()
}

func main() {

	if *state != "" {
		server_state.ShowServerState("localhost:9091")
		return
	}

	config.MustLoad()

	err := db.Init(&db.MySQLConfig{
		Host:     config.MySql.Host,
		Port:     config.MySql.Port,
		User:     config.MySql.Username,
		Password: config.MySql.Password,
		Database: config.MySql.Db,
		Charset:  config.MySql.Charset,
	}, &db.RedisConfig{
		Host:     config.Redis.Host,
		Port:     config.Redis.Port,
		Password: config.Redis.Password,
		Db:       config.Redis.Db,
	})
	if err != nil {
		panic(err)
	}

	gateway, err := im_server.NewServer(config.WsServer.ID, config.WsServer.Addr, config.WsServer.Port)
	if err != nil {
		panic(err)
	}

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

	handler, err := message_handler.NewHandlerWithOptions(&message_handler.Options{
		MessageStore:           cStore,
		Auth:                   jwt_auth.NewAuthorizeImpl(config.WsServer.JwtSecret),
		DontInitDefaultHandler: true,
		NotifyOnErr:            true,
	})
	if err != nil {
		panic(err)
	}
	if config.Common.StoreOfflineMessage {
		message_handler.Enable = true
		handler.SetOfflineMessageHandler(message_handler.GetHandleFn())
	}
	action_handler.Setup(handler)
	handler.InitDefaultHandler(nil)

	subscription := subscription_impl.NewSubscription(sStore, seqStore)
	options := bootstrap.Options{
		Messaging:    handler,
		Gate:         gateway,
		Subscription: subscription,
	}

	go func() {
		logger.D("websocket listening on %s:%d", config.WsServer.Addr, config.WsServer.Port)
		err = bootstrap.Bootstrap(&options)

		if err != nil {
			panic(err)
		}
	}()

	go func() {
		logger.D("state server is listening on 0.0.0.0:%d", 9091)
		server_state.StartSrv(9091, gateway)
	}()

	rpcOpts := rpc.ServerOptions{
		Name:    config.IMService.Name,
		Network: config.IMService.Network,
		Addr:    config.IMService.Addr,
		Port:    config.IMService.Port,
	}
	logger.D("rpc %s listening on %s %s:%d", rpcOpts.Name, rpcOpts.Network, rpcOpts.Addr, rpcOpts.Port)
	err = im_server.RunRpcServer(&rpcOpts, gateway, subscription)
	if err != nil {
		panic(err)
	}
}
