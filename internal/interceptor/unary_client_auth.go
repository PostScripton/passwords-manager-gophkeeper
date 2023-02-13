package interceptor

import (
	"context"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type UnaryClientAuthInterceptor struct {
	settingsRepo repository.Settings
}

func NewUnaryClientAuthInterceptor(settingsRepo repository.Settings) *UnaryClientAuthInterceptor {
	return &UnaryClientAuthInterceptor{
		settingsRepo: settingsRepo,
	}
}

func (i *UnaryClientAuthInterceptor) Handle() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		jwt, ok, err := i.settingsRepo.Get(ctx, "jwt")
		if err != nil {
			return err
		}
		if ok {
			ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{"authorization": jwt}))
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
