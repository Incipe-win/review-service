package biz

import (
	"context"
	v1 "review-service/api/review/v1"
	"review-service/internal/data/model"
	"review-service/pkg/snowflake"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrReviewNotFound is user not found.
	ErrReviewNotFound  = errors.NotFound(v1.ErrorReason_REVIEW_NOT_FOUND.String(), "user not found")
	ErrSQL             = errors.InternalServer(v1.ErrorReason_SQL_ERROR.String(), "internal server error")
	ErrAlreadyReviewed = errors.BadRequest(v1.ErrorReason_INVALID_PARAMETERS.String(), "already reviewed this order")
)

// ReviewRepo is a Review repo.
type ReviewRepo interface {
	SaveReview(context.Context, *model.ReviewInfo) (*model.ReviewInfo, error)
	GetReviewByOrderID(context.Context, int64) ([]*model.ReviewInfo, error)
}

// ReviewUsecase is a Review usecase.
type ReviewUsecase struct {
	repo ReviewRepo
	log  *log.Helper
}

func NewReviewUsecase(repo ReviewRepo, logger log.Logger) *ReviewUsecase {
	return &ReviewUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *ReviewUsecase) CreateReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	uc.log.WithContext(ctx).Debugf("CreateReview, review: %v", review)
	reviews, err := uc.repo.GetReviewByOrderID(ctx, review.OrderID)
	if err != nil {
		return nil, ErrSQL
	}
	if len(reviews) > 0 {
		return nil, ErrAlreadyReviewed
	}
	review.ReviewID = snowflake.GenID()
	return uc.repo.SaveReview(ctx, review)
}
