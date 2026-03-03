package middleware

import (
	"context"

	"github.com/eatdetey/letterboxd-replica/source/go-common/pkg/ctxmetadata"
	"google.golang.org/grpc"
)

func RequestIdMiddleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx, _ = ctxmetadata.EnsureReqId(ctx)

		return handler(ctx, req)
	}
}
