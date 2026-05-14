package service

import (
	"errors"
	"gin-quickstart/internal/model"
	"gin-quickstart/internal/repository"
	"time"
)

type NotificationService struct {
	Repo *repository.NotificationRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{Repo: repo}
}

func (s NotificationService) GetNotificationsByUserID(userID uint64) ([]model.Notification, error) {
	return s.Repo.GetNotificationsByUserID(userID)
}

func (s NotificationService) GetNotificationByID(id uint64, userID uint64) (*model.Notification, error) {
	notification, err := s.Repo.GetNotificationByID(id)
	if err != nil {
		return nil, err
	}

	if notification.UserId != uint(userID) {
		return nil, errors.New("notification not found")
	}

	return notification, nil
}

func (s NotificationService) CreateNotification(notification *model.Notification) (*model.Notification, error) {
	if notification.UserId == 0 {
		return nil, errors.New("user_id is required")
	}

	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()

	err := s.Repo.Create(notification)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (s NotificationService) MarkNotificationAsRead(id uint64, userID uint64) (*model.Notification, error) {
	notification, err := s.GetNotificationByID(id, userID)
	if err != nil {
		return nil, err
	}

	if notification.IsRead {
		return notification, nil
	}

	now := time.Now()
	notification.IsRead = true
	notification.ReadAt = &now
	notification.UpdatedAt = now

	err = s.Repo.Update(notification)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (s NotificationService) UpdateNotificationReadState(id uint64, userID uint64, isRead bool) (*model.Notification, error) {
	notification, err := s.GetNotificationByID(id, userID)
	if err != nil {
		return nil, err
	}

	notification.IsRead = isRead
	if isRead {
		now := time.Now()
		notification.ReadAt = &now
	} else {
		notification.ReadAt = nil
	}
	notification.UpdatedAt = time.Now()

	err = s.Repo.Update(notification)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (s NotificationService) DeleteNotification(id uint64, userID uint64) error {
	notification, err := s.GetNotificationByID(id, userID)
	if err != nil {
		return err
	}

	return s.Repo.Delete(notification)
}
