package data

import (
	"context"
	v1 "review-service/api/review/v1"
	"review-service/internal/biz"
	"review-service/internal/data/model"
	"review-service/internal/data/query"

	"github.com/go-kratos/kratos/v2/log"
)

type reviewRepo struct {
	data *Data
	log  *log.Helper
}

// NewReviewRepo .
func NewReviewRepo(data *Data, logger log.Logger) biz.ReviewRepo {
	return &reviewRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *reviewRepo) SaveReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	err := r.data.query.ReviewInfo.WithContext(ctx).Save(review)
	return review, err
}

func (r *reviewRepo) GetReviewByOrderID(ctx context.Context, orderID int64) ([]*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.OrderID.Eq(orderID)).Find()
}

func (r *reviewRepo) GetReview(ctx context.Context, reviewID int64) (*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.ReviewID.Eq(reviewID)).First()
}

func (r *reviewRepo) AuditReview(ctx context.Context, param *biz.AuditParam) error {
	return nil
}

func (r *reviewRepo) AppealReview(ctx context.Context, param *biz.AppealParam) error {
	return nil
}

func (r *reviewRepo) AuditAppeal(ctx context.Context, param *biz.AuditAppealParam) error {
	return nil
}

func (r *reviewRepo) SaveReply(ctx context.Context, reply *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error) {
	review, err := r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.ReviewID.Eq(reply.ReviewID)).First()
	if err != nil {
		return nil, err
	}
	if review.HasReply == 1 {
		return nil, v1.ErrorDuplicateReply("review %d has already been replied", reply.ReviewID)
	}
	if review.StoreID != reply.StoreID {
		return nil, v1.ErrorHorizontalOversteppingOfAuthority("store %d has no authority to reply review %d", reply.StoreID, reply.ReviewID)
	}
	err = r.data.query.Transaction(func(tx *query.Query) error {
		if err := tx.ReviewReplyInfo.WithContext(ctx).Save(reply); err != nil {
			r.log.WithContext(ctx).Errorf("failed to save reply: %v, error: %v", reply, err)
			return err
		}
		if _, err := tx.ReviewInfo.WithContext(ctx).Where(tx.ReviewInfo.ReviewID.Eq(reply.ReviewID)).Update(tx.ReviewInfo.HasReply, 1); err != nil {
			r.log.WithContext(ctx).Errorf("update review failed: %v, error: %v", reply.ReviewID, err)
			return err
		}
		return nil
	})
	return reply, err
}

func (r *reviewRepo) ListReviewByUserID(ctx context.Context, userID int64, offset, limit int) ([]*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewInfo.UserID.Eq(userID)).
		Order(r.data.query.ReviewInfo.ID.Desc()).
		Limit(limit).
		Offset(offset).
		Find()
}

func (r *reviewRepo) GetReviewReply(ctx context.Context, reviewID int64) (*model.ReviewReplyInfo, error) {
	return r.data.query.ReviewReplyInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewReplyInfo.ReviewID.Eq(reviewID)).
		First()
}
