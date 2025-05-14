package service

import (
	"context"
	"errors"
	"time"
	
	"shop/backend/profile/internal/domain/entity"
)

// ErrFeedbackNotFound 反馈不存在错误
var ErrFeedbackNotFound = errors.New("feedback not found")

// FeedbackServiceImpl 用户反馈服务实现
type FeedbackServiceImpl struct {
	repo ProfileRepository
}

// NewFeedbackService 创建用户反馈服务实例
func NewFeedbackService(repo ProfileRepository) FeedbackService {
	return &FeedbackServiceImpl{
		repo: repo,
	}
}

// ListFeedbacks 获取用户反馈列表
func (s *FeedbackServiceImpl) ListFeedbacks(ctx context.Context, userID int64, page, pageSize int) ([]*entity.UserFeedback, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.GetFeedbacksByUserID(ctx, userID, offset, pageSize)
}

// GetFeedback 获取反馈详情
func (s *FeedbackServiceImpl) GetFeedback(ctx context.Context, id int64) (*entity.UserFeedback, error) {
	feedback, err := s.repo.GetFeedbackByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if feedback == nil {
		return nil, ErrFeedbackNotFound
	}
	
	return feedback, nil
}

// SubmitFeedback 提交反馈
func (s *FeedbackServiceImpl) SubmitFeedback(ctx context.Context, feedback *entity.UserFeedback) error {
	now := time.Now()
	feedback.Status = 0 // 待处理状态
	feedback.CreatedAt = now
	feedback.UpdatedAt = now
	
	return s.repo.CreateFeedback(ctx, feedback)
}

// UpdateFeedback 更新反馈
func (s *FeedbackServiceImpl) UpdateFeedback(ctx context.Context, feedback *entity.UserFeedback) error {
	// 检查反馈是否存在
	existingFeedback, err := s.repo.GetFeedbackByID(ctx, feedback.ID)
	if err != nil {
		return err
	}
	
	if existingFeedback == nil {
		return ErrFeedbackNotFound
	}
	
	// 检查是否为同一用户
	if existingFeedback.UserID != feedback.UserID {
		return ErrUserNotMatch
	}
	
	// 只有待处理状态的反馈才能被用户修改
	if existingFeedback.Status != 0 {
		return errors.New("feedback cannot be modified in current status")
	}
	
	// 用户只能修改特定字段
	existingFeedback.Subject = feedback.Subject
	existingFeedback.Content = feedback.Content
	existingFeedback.FileURLs = feedback.FileURLs
	existingFeedback.UpdatedAt = time.Now()
	
	return s.repo.UpdateFeedback(ctx, existingFeedback)
}

// DeleteFeedback 删除反馈
func (s *FeedbackServiceImpl) DeleteFeedback(ctx context.Context, id int64) error {
	// 检查反馈是否存在
	feedback, err := s.repo.GetFeedbackByID(ctx, id)
	if err != nil {
		return err
	}
	
	if feedback == nil {
		return ErrFeedbackNotFound
	}
	
	return s.repo.DeleteFeedback(ctx, id)
}

// BrowsingHistoryServiceImpl 浏览历史服务实现
type BrowsingHistoryServiceImpl struct {
	repo ProfileRepository
}

// NewBrowsingHistoryService 创建浏览历史服务实例
func NewBrowsingHistoryService(repo ProfileRepository) BrowsingHistoryService {
	return &BrowsingHistoryServiceImpl{
		repo: repo,
	}
}

// GetHistories 获取用户浏览历史
func (s *BrowsingHistoryServiceImpl) GetHistories(ctx context.Context, userID int64, page, pageSize int) ([]*entity.BrowsingHistory, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.GetBrowsingHistories(ctx, userID, offset, pageSize)
}

// AddHistory 添加浏览记录
func (s *BrowsingHistoryServiceImpl) AddHistory(ctx context.Context, userID, goodsID int64, source string, stayTime int) error {
	// 创建浏览记录
	history := &entity.BrowsingHistory{
		UserID:    userID,
		GoodsID:   goodsID,
		Source:    source,
		StayTime:  stayTime,
		CreatedAt: time.Now(),
	}
	
	return s.repo.AddBrowsingHistory(ctx, history)
}

// RemoveHistories 删除浏览记录
func (s *BrowsingHistoryServiceImpl) RemoveHistories(ctx context.Context, userID int64, ids []int64) error {
	// 确保只删除该用户自己的浏览记录
	return s.repo.DeleteBrowsingHistory(ctx, userID, ids)
}

// ClearHistories 清空浏览记录
func (s *BrowsingHistoryServiceImpl) ClearHistories(ctx context.Context, userID int64) error {
	return s.repo.ClearBrowsingHistory(ctx, userID)
}
