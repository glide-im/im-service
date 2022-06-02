package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/glide-im/glide/pkg/gate"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/im-service/pkg/proto"
	"github.com/glide-im/im-service/pkg/rpc"
	"strings"
)

const (
	errRpcInvocation = "rpc invocation error: "
)

// IsRpcInvocationError
// Rpc invocation failed errors are returned by the rpc client when the rpc call fails.
func IsRpcInvocationError(err error) bool {
	return err != nil && strings.HasPrefix(err.Error(), errRpcInvocation)
}

// IMServiceClient is the client for the IM service.
type IMServiceClient struct {
	rpc *imRpcClient
}

// NewIMServiceClient .
func NewIMServiceClient(options *rpc.ClientOptions) (*IMServiceClient, error) {
	client, err := newIMRpcClient(options)
	if err != nil {
		return nil, err
	}
	return &IMServiceClient{
		rpc: client,
	}, nil
}

func (i *IMServiceClient) SetClientID(old gate.ID, new_ gate.ID) error {
	ctx := context.TODO()
	request := proto.SetIDRequest{
		OldId: string(old),
		NewId: string(new_),
	}
	response := proto.Response{}
	err := i.rpc.SetClientID(ctx, &request, &response)
	if err != nil {
		return errors.New(errRpcInvocation + err.Error())
	}
	return getResponseError(&response)
}

func (i *IMServiceClient) ExitClient(id gate.ID) error {
	ctx := context.TODO()
	request := proto.ExitClientRequest{
		Id: string(id),
	}
	response := proto.Response{}
	err := i.rpc.ExitClient(ctx, &request, &response)
	if err != nil {
		return errors.New(errRpcInvocation + err.Error())
	}
	return getResponseError(&response)
}

func (i *IMServiceClient) IsOnline(id gate.ID) bool {
	ctx := context.TODO()
	request := proto.IsOnlineRequest{
		Id: string(id),
	}
	response := proto.IsOnlineResponse{}
	err := i.rpc.IsOnline(ctx, &request, &response)
	if err != nil {
		return false
	}
	return response.GetOnline()
}

func (i *IMServiceClient) EnqueueMessage(id gate.ID, message *messages.GlideMessage) error {

	marshal, err := json.Marshal(message)
	if err != nil {
		return err
	}
	ctx := context.TODO()
	request := proto.EnqueueMessageRequest{
		Id:  string(id),
		Msg: marshal,
	}
	response := proto.Response{}
	err = i.rpc.EnqueueMessage(ctx, &request, &response)
	if err != nil {
		return errors.New(errRpcInvocation + err.Error())
	}
	return getResponseError(&response)
}

func getResponseError(response *proto.Response) error {
	if proto.Response_ResponseCode(response.GetCode()) != proto.Response_OK {
		return fmt.Errorf("%d: %s", response.GetCode(), response.GetMsg())
	}
	return nil
}
