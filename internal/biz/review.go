package biz

import (
	"context"
	v1 "review-service/api/review/v1"
	"review-service/internal/data/model"
	"review-service/pkg/snowflake"

	"github.com/go-kratos/kratos/v2/log"
)

// ReviewRepo is a Review repo.
type ReviewRepo interface {
	SaveReview(context.Context, *model.ReviewInfo) (*model.ReviewInfo, error)
	GetReviewByOrderID(context.Context, int64) ([]*model.ReviewInfo, error)
	GetReview(context.Context, int64) (*model.ReviewInfo, error)
	SaveReply(context.Context, *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error)
	GetReviewReply(context.Context, int64) (*model.ReviewReplyInfo, error)
	AuditReview(context.Context, *AuditParam) error
	AppealReview(context.Context, *AppealParam) error
	AuditAppeal(context.Context, *AuditAppealParam) error
	ListReviewByUserID(ctx context.Context, userID int64, offset, limit int) ([]*model.ReviewInfo, error)
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
	reviews, err := uc.repo.GetReviewByOrderID(ctx, review.OrderID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to get reviews by order ID: %v, error: %v", review.OrderID, err)
		return nil, v1.ErrorDbFailed("db error")
	}
	if len(reviews) > 0 {
		return nil, v1.ErrorOrderReviewed("order %d has already been reviewed", review.OrderID)
	}
	review.ReviewID = snowflake.GenID()
	return uc.repo.SaveReview(ctx, review)
}

func (uc *ReviewUsecase) GetReview(ctx context.Context, reviewID int64) (*model.ReviewInfo, error) {
	return uc.repo.GetReview(ctx, reviewID)
}

func (uc *ReviewUsecase) AuditReview(ctx context.Context, param *AuditParam) error {
	return uc.repo.AuditReview(ctx, param)
}

func (uc *ReviewUsecase) CreateReply(ctx context.Context, param *ReplyParam) (*model.ReviewReplyInfo, error) {
	reply := &model.ReviewReplyInfo{
		ReplyID:   snowflake.GenID(),
		ReviewID:  param.ReviewID,
		StoreID:   param.StoreID,
		Content:   param.Content,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	}
	return uc.repo.SaveReply(ctx, reply)
}
