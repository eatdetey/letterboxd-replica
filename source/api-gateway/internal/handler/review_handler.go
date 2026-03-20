package handler

import (
	"strings"

	reviewpb "github.com/eatdetey/letterboxd-replica/source/api-gateway/gen/go/review/v1"
	"github.com/gofiber/fiber/v3"
)

type ReviewHandler struct {
	client reviewpb.ReviewServiceClient
}

type reviewResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}

type addReviewRequest struct {
	Text string `json:"text"`
}

func NewReviewHandler(client reviewpb.ReviewServiceClient) *ReviewHandler {
	return &ReviewHandler{client: client}
}

func (h *ReviewHandler) GetMovieReviews(c fiber.Ctx) error {
	movieID := strings.TrimSpace(c.Params("id"))
	if movieID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "movie id is required")
	}

	ctx, cancel := buildGRPCContext(c)
	defer cancel()

	resp, err := h.client.GetReviews(ctx, &reviewpb.GetReviewsRequest{
		MovieId: movieID,
	})
	if err != nil {
		return grpcErrorToFiber(err)
	}

	items := make([]reviewResponse, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, mapReviewResponse(item))
	}

	return c.JSON(fiber.Map{"items": items})
}

func (h *ReviewHandler) GetMovieReviewByID(c fiber.Ctx) error {
	movieID := strings.TrimSpace(c.Params("id"))
	if movieID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "movie id is required")
	}
	reviewID := strings.TrimSpace(c.Params("review_id"))
	if reviewID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "review id is required")
	}

	ctx, cancel := buildGRPCContext(c)
	defer cancel()

	resp, err := h.client.GetReviews(ctx, &reviewpb.GetReviewsRequest{
		MovieId: movieID,
	})
	if err != nil {
		return grpcErrorToFiber(err)
	}

	for _, review := range resp.Items {
		if review.GetId() == reviewID {
			return c.JSON(fiber.Map{
				"review": mapReviewResponse(review),
			})
		}
	}

	return fiber.NewError(fiber.StatusNotFound, "review not found")
}

func (h *ReviewHandler) AddMovieReview(c fiber.Ctx) error {
	movieID := strings.TrimSpace(c.Params("id"))
	if movieID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "movie id is required")
	}

	var reqBody addReviewRequest
	if err := c.Bind().Body(&reqBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	ctx, cancel := buildGRPCContext(c)
	defer cancel()

	resp, err := h.client.AddReview(ctx, &reviewpb.AddReviewRequest{
		MovieId: movieID,
		Text:    reqBody.Text,
	})
	if err != nil {
		return grpcErrorToFiber(err)
	}

	review := resp.GetReview()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"review": mapReviewResponse(review),
	})
}

func mapReviewResponse(item *reviewpb.Review) reviewResponse {
	return reviewResponse{
		ID:        item.GetId(),
		UserID:    item.GetUserId(),
		Username:  item.GetUsername(),
		Text:      item.GetText(),
		CreatedAt: item.GetCreatedAt(),
	}
}
