package repository

import (
	"gin-quickstart/internal/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	GormDB      *gorm.DB
	RedisClient *redis.Client
}

func NewNotificationRepository(db *gorm.DB, redis *redis.Client) *NotificationRepository {
	return &NotificationRepository{
		GormDB:      db,
		RedisClient: redis,
	}
}

func (r NotificationRepository) GetNotificationsByUserID(userID uint64) ([]model.Notification, error) {
	var notifications []model.Notification
	err := r.GormDB.Where("user_id = ?", userID).
		Order("created_at desc").
		Preload("User").
		Preload("Thread").
		Preload("Post").
		Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r NotificationRepository) GetNotificationByID(id uint64) (*model.Notification, error) {
	var notification model.Notification
	err := r.GormDB.Preload("User").Preload("Thread").Preload("Post").First(&notification, id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *NotificationRepository) Create(notification *model.Notification) error {
	return r.GormDB.Create(notification).Error
}

func (r *NotificationRepository) Update(notification *model.Notification) error {
	return r.GormDB.Save(notification).Error
}

func (r *NotificationRepository) Delete(notification *model.Notification) error {
	return r.GormDB.Delete(notification).Error
}
