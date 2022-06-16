package client

import (
	"context"
	"github.com/glide-im/glide/pkg/rpc"
	"github.com/glide-im/im-service/pkg/proto"
	"github.com/glide-im/im-service/pkg/server"
)

var _ server.IMRpcServer = &imRpcClient{}

type imRpcClient struct {
	cli *rpc.BaseClient
}

func newIMRpcClient(options *rpc.ClientOptions) (*imRpcClient, error) {
	client, err := rpc.NewBaseClient(options)
	if err != nil {
		return nil, err
	}
	return &imRpcClient{
		cli: client,
	}, nil
}

func (I *imRpcClient) SetClientID(ctx context.Context, request *proto.SetIDRequest, response *proto.Response) error {
	return I.cli.Call(ctx, "SetClientID", request, response)
}

func (I *imRpcClient) ExitClient(ctx context.Context, request *proto.ExitClientRequest, response *proto.Response) error {
	return I.cli.Call(ctx, "ExitClient", request, response)
}

func (I *imRpcClient) IsOnline(ctx context.Context, request *proto.IsOnlineRequest, response *proto.IsOnlineResponse) error {
	return I.cli.Call(ctx, "IsOnline", request, response)
}

func (I *imRpcClient) EnqueueMessage(ctx context.Context, request *proto.EnqueueMessageRequest, response *proto.Response) error {
	return I.cli.Call(ctx, "EnqueueMessage", request, response)
}
