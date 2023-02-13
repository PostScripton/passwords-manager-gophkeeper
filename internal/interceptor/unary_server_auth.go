package interceptor

import (
	"context"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strconv"
)

var unprotectedRoutes = []string{"/proto.User/Register", "/proto.User/Login"}

type UnaryServerAuthInterceptor struct {
	authService services.Auth
}

func NewUnaryServerAuthInterceptor(authService services.Auth) *UnaryServerAuthInterceptor {
	return &UnaryServerAuthInterceptor{
		authService: authService,
	}
}

func (i *UnaryServerAuthInterceptor) Handle() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		newCtx := ctx
		if i.isProtectedRoute(info.FullMethod) {
			newCtx, err = i.authorize(ctx)
			if err != nil {
				return nil, err
			}
		}

		return handler(newCtx, req)
	}
}

func (i *UnaryServerAuthInterceptor) isProtectedRoute(routeName string) bool {
	for _, unprotectedRoute := range unprotectedRoutes {
		if unprotectedRoute == routeName {
			return false
		}
	}

	return true
}

func (i *UnaryServerAuthInterceptor) authorize(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
	}

	tokenString := authHeaders[0]
	token, err := i.authService.ParseJWT(tokenString)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "Authorization failed")
	}

	if !token.Valid {
		return ctx, status.Errorf(codes.Unauthenticated, "JWT is invalid")
	}

	id, err := i.authService.GetIDFromJWT(token)
	if err != nil {
		return ctx, status.Errorf(codes.Unimplemented, "Unable get user ID from JWT claims")
	}

	ctx = metadata.NewIncomingContext(ctx, metadata.New(map[string]string{"user_id": strconv.Itoa(id)}))

	return ctx, nil
}
