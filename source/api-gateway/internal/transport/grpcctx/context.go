package grpcctx

import (
	"context"
	"fmt"

	"github.com/eatdetey/letterboxd-replica/source/go-common/pkg/ctxmetadata"
	"github.com/gofiber/fiber/v3"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "Authorization"
	accessCookieName    = "access_token"
	refreshCookieName   = "refresh_token"
)

// FromFiber builds outgoing gRPC context with request-id and auth metadata based on Fiber context.
func FromFiber(c fiber.Ctx, reqID string) context.Context {
	ctx := context.Background()

	if reqID != "" {
		ctx = ctxmetadata.WithReqId(ctx, reqID)
	} else {
		ctx, _ = ctxmetadata.EnsureReqId(ctx)
	}

	md := metadata.MD{}

	if rid := ctxmetadata.GetReqIdFromContext(ctx); rid != "" {
		md.Set(ctxmetadata.RequestIDKey, rid)
	}

	if header := c.Get(authorizationHeader); header != "" {
		md.Set(ctxmetadata.AuthorizationKey, header)
	} else if token := c.Cookies(accessCookieName); token != "" {
		md.Set(ctxmetadata.AuthorizationKey, fmt.Sprintf("Bearer %s", token))
	} else if token := c.Cookies(refreshCookieName); token != "" {
		md.Set(ctxmetadata.AuthorizationKey, fmt.Sprintf("Bearer %s", token))
	}

	return metadata.NewOutgoingContext(ctx, md)
}
