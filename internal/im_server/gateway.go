package im_server

import (
	"github.com/glide-im/glide/pkg/conn"
	"github.com/glide-im/glide/pkg/gate"
	"github.com/glide-im/glide/pkg/gate/gateway"
	"github.com/glide-im/glide/pkg/logger"
	"github.com/glide-im/glide/pkg/messages"
	"sync/atomic"
	"time"
)

type GatewayState struct {
	ServerId             string    `json:"server_id"`
	Addr                 string    `json:"addr"`
	Port                 int       `json:"port"`
	StartAt              time.Time `json:"start_at"`
	RunningHours         float64   `json:"running_hours"`
	ConnectedClientCount int32     `json:"connected_client_count"`
	OnlineClients        int32     `json:"online_clients"`
	OnlineTempClients    int32     `json:"online_temp_clients"`
	DeliveredMessages    int32     `json:"delivered_messages"`
	DeliverMessageFails  int32     `json:"deliver_message_fails"`
	ReceivedMessages     int32     `json:"received_messages"`
}

type GatewayServer struct {
	*gateway.Impl

	server conn.Server
	h      gate.MessageHandler

	gateID string
	addr   string
	port   int

	state *GatewayState
}

func NewServer(id string, addr string, port int) (*GatewayServer, error) {
	srv := GatewayServer{}
	srv.Impl, _ = gateway.NewServer(
		&gateway.Options{
			ID:                    id,
			MaxMessageConcurrency: 30_0000,
		},
	)
	srv.state = &GatewayState{
		ServerId: id,
		Addr:     addr,
		Port:     port,
	}
	srv.addr = addr
	srv.port = port
	srv.gateID = id

	options := &conn.WsServerOptions{
		ReadTimeout:  time.Minute * 3,
		WriteTimeout: time.Minute * 3,
	}
	srv.server = conn.NewWsServer(options)
	return &srv, nil
}

func (c *GatewayServer) Run() error {

	c.state.StartAt = time.Now()

	c.server.SetConnHandler(func(conn conn.Connection) {
		c.HandleConnection(conn)
		atomic.AddInt32(&c.state.ConnectedClientCount, 1)
	})
	return c.server.Run(c.addr, c.port)
}

func (c *GatewayServer) SetMessageHandler(h gate.MessageHandler) {
	handler := func(id *gate.Info, msg *messages.GlideMessage) {
		atomic.AddInt32(&c.state.ReceivedMessages, 1)
		h(id, msg)
	}
	c.h = handler
	c.Impl.SetMessageHandler(handler)
}

// HandleConnection 当一个用户连接建立后, 由该方法创建 Client 实例 Client 并管理该连接, 返回该由连接创建客户端的标识 id
// 返回的标识 id 是一个临时 id, 后续连接认证后会改变
func (c *GatewayServer) HandleConnection(conn conn.Connection) gate.ID {

	// 获取一个临时 uid 标识这个连接
	id, err := gate.GenTempID(c.gateID)
	if err != nil {
		logger.E("[gateway] gen temp id error: %v", err)
		return ""
	}
	ret := gateway.NewClientWithConfig(conn, c, c.h, &gateway.ClientConfig{
		HeartbeatLostLimit:      3,
		ClientHeartbeatDuration: time.Second * 30,
		ServerHeartbeatDuration: time.Second * 30,
		CloseImmediately:        false,
	})
	ret.SetID(id)
	c.Impl.AddClient(ret)

	// 开始处理连接的消息
	ret.Run()

	hello := messages.ServerHello{
		TempID:            id.UID(),
		HeartbeatInterval: 30,
	}

	m := messages.NewMessage(0, messages.ActionHello, hello)
	_ = ret.EnqueueMessage(m)

	return id
}

func (c *GatewayServer) EnqueueMessage(id gate.ID, msg *messages.GlideMessage) error {
	err := c.Impl.EnqueueMessage(id, msg)
	if err != nil {
		atomic.AddInt32(&c.state.DeliverMessageFails, 1)
	} else {
		atomic.AddInt32(&c.state.DeliveredMessages, 1)
	}
	return err
}

func (c *GatewayServer) GetState() GatewayState {
	all := c.Impl.GetAll()
	temp := 0
	for id := range all {
		if id.IsTemp() {
			temp++
		}
	}
	c.state.OnlineTempClients = int32(temp)
	c.state.OnlineClients = int32(len(all))
	span := time.Now().Unix() - c.state.StartAt.Unix()
	c.state.RunningHours = float64(span) / 60.0 / 60.0
	return *c.state
}
