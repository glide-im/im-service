package server

import (
	"context"
	"github.com/glide-im/im-service/pkg/proto"
)

type IMRpcServer interface {
	SetClientID(ctx context.Context, request *proto.SetIDRequest, response *proto.Response) error

	ExitClient(ctx context.Context, request *proto.ExitClientRequest, response *proto.Response) error

	IsOnline(ctx context.Context, request *proto.IsOnlineRequest, response *proto.IsOnlineResponse) error

	EnqueueMessage(ctx context.Context, request *proto.EnqueueMessageRequest, response *proto.Response) error
}
