package im_server

import (
	"context"
	"encoding/json"
	"github.com/glide-im/glide/pkg/gate"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/im-service/internal/rpc"
	"github.com/glide-im/im-service/pkg/proto"
)

type RpcServer struct {
	s gate.Gateway
}

func RunRpcServer(options *rpc.ServerOptions, gate gate.Gateway) error {
	server := rpc.NewBaseServer(options)
	rpcServer := RpcServer{
		s: gate,
	}
	server.Register(options.Name, &rpcServer)
	return server.Run()
}

func (r *RpcServer) SetClientID(ctx context.Context, request *proto.SetIDRequest, response *proto.Response) error {
	err := r.s.SetClientID(gate.ID(request.OldId), gate.ID(request.NewId))
	if err != nil {
		response.Code = int32(proto.Response_ERROR)
		response.Msg = err.Error()
	}
	return nil
}

func (r *RpcServer) ExitClient(ctx context.Context, request *proto.ExitClientRequest, response *proto.Response) error {
	err := r.s.ExitClient(gate.ID(request.Id))
	if err != nil {
		response.Code = int32(proto.Response_ERROR)
		response.Msg = err.Error()
	}
	return nil
}

func (r *RpcServer) IsOnline(ctx context.Context, request *proto.IsOnlineRequest, response *proto.IsOnlineResponse) error {
	response.Online = r.s.IsOnline(gate.ID(request.Id))
	return nil
}

func (r *RpcServer) EnqueueMessage(ctx context.Context, request *proto.EnqueueMessageRequest, response *proto.Response) error {

	msg := messages.GlideMessage{}
	err := json.Unmarshal(request.Msg, &msg)
	if err != nil {
		response.Code = int32(proto.Response_ERROR)
		response.Msg = err.Error()
		return nil
	}

	err = r.s.EnqueueMessage(gate.ID(request.Id), &msg)
	if err != nil {
		response.Code = int32(proto.Response_ERROR)
		response.Msg = err.Error()
	}
	return nil
}
