package handler

import (
	"context"
	"strconv"
	"strings"
	"time"

	userpb "github.com/eatdetey/letterboxd-replica/source/api-gateway/gen/go/user/v1"
	"github.com/eatdetey/letterboxd-replica/source/api-gateway/internal/mapper"
	"github.com/eatdetey/letterboxd-replica/source/api-gateway/internal/transport/grpcctx"
	"github.com/gofiber/fiber/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	client userpb.UserServiceClient
}

func NewUserHandler(client userpb.UserServiceClient) *UserHandler {
	return &UserHandler{client: client}
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type getUsersRequest struct {
	IDs       []int64  `json:"ids"`
	Usernames []string `json:"usernames"`
	Limit     int32    `json:"limit"`
	Offset    int32    `json:"offset"`
}

func (h *UserHandler) Register(c fiber.Ctx) error {
	var req registerRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	reqID, _ := c.Locals("request_id").(string)
	ctx := grpcctx.FromFiber(c, reqID)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := h.client.Register(ctx, &userpb.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		return grpcErrorToFiber(err)
	}

	setRefreshCookie(c, resp.GetTokens().GetRefreshToken())

	return c.JSON(map[string]any{
		"user":         mapper.UserFromPB(resp.User),
		"access_token": resp.GetTokens().GetAccessToken(),
	})
}

func (h *UserHandler) Login(c fiber.Ctx) error {
	var req loginRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	reqID, _ := c.Locals("request_id").(string)
	ctx := grpcctx.FromFiber(c, reqID)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := h.client.Login(ctx, &userpb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return grpcErrorToFiber(err)
	}

	setRefreshCookie(c, resp.GetTokens().GetRefreshToken())

	return c.JSON(map[string]any{
		"user":         mapper.UserFromPB(resp.User),
		"access_token": resp.GetTokens().GetAccessToken(),
	})
}

func (h *UserHandler) Refresh(c fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return fiber.NewError(fiber.StatusBadRequest, "refresh_token cookie is required")
	}

	reqID, _ := c.Locals("request_id").(string)
	ctx := grpcctx.FromFiber(c, reqID)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := h.client.Refresh(ctx, &userpb.RefreshRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return grpcErrorToFiber(err)
	}

	setRefreshCookie(c, resp.GetTokens().GetRefreshToken())

	return c.JSON(map[string]any{
		"access_token": resp.GetTokens().GetAccessToken(),
	})
}

func (h *UserHandler) Logout(c fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken != "" {
		reqID, _ := c.Locals("request_id").(string)
		ctx := grpcctx.FromFiber(c, reqID)
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if _, err := h.client.Logout(ctx, &userpb.LogoutRequest{
			RefreshToken: refreshToken,
		}); err != nil {
			return grpcErrorToFiber(err)
		}
	}

	clearRefreshCookie(c)
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandler) GetUsers(c fiber.Ctx) error {
	var reqBody getUsersRequest
	if err := c.Bind().Body(&reqBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	usernames := make([]string, 0, len(reqBody.Usernames))
	for _, username := range reqBody.Usernames {
		if trimmed := strings.TrimSpace(username); trimmed != "" {
			usernames = append(usernames, trimmed)
		}
	}

	if len(reqBody.IDs) == 0 && len(usernames) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "ids or usernames must be provided")
	}

	reqID, _ := c.Locals("request_id").(string)
	ctx := grpcctx.FromFiber(c, reqID)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := h.client.GetUsers(ctx, &userpb.GetUsersRequest{
		Ids:       reqBody.IDs,
		Usernames: usernames,
		Limit:     reqBody.Limit,
		Offset:    reqBody.Offset,
	})
	if err != nil {
		return grpcErrorToFiber(err)
	}

	users := make([]mapper.UserResponse, 0, len(resp.Users))
	for _, u := range resp.Users {
		users = append(users, mapper.UserFromPB(u))
	}

	return c.JSON(map[string]any{
		"users": users,
	})
}

func (h *UserHandler) GetUserByID(c fiber.Ctx) error {
	idRaw := strings.TrimSpace(c.Params("id"))
	if idRaw == "" {
		return fiber.NewError(fiber.StatusBadRequest, "user id is required")
	}
	id, err := strconv.ParseInt(idRaw, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "user id must be integer")
	}

	reqID, _ := c.Locals("request_id").(string)
	ctx := grpcctx.FromFiber(c, reqID)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := h.client.GetUsers(ctx, &userpb.GetUsersRequest{
		Ids:   []int64{id},
		Limit: 1,
	})
	if err != nil {
		return grpcErrorToFiber(err)
	}
	if len(resp.Users) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	return c.JSON(map[string]any{
		"user": mapper.UserFromPB(resp.Users[0]),
	})
}

func parseStrings(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func grpcErrorToFiber(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return fiber.NewError(fiber.StatusBadRequest, st.Message())
	case codes.NotFound:
		return fiber.NewError(fiber.StatusNotFound, st.Message())
	case codes.AlreadyExists:
		return fiber.NewError(fiber.StatusConflict, st.Message())
	case codes.Unauthenticated:
		return fiber.NewError(fiber.StatusUnauthorized, st.Message())
	case codes.PermissionDenied:
		return fiber.NewError(fiber.StatusForbidden, st.Message())
	default:
		return fiber.NewError(fiber.StatusInternalServerError, st.Message())
	}
}

func setRefreshCookie(c fiber.Ctx, refreshToken string) {
	if refreshToken == "" {
		return
	}
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
	})
}

func clearRefreshCookie(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}
