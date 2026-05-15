package repository

import (
	"gin-quickstart/internal/model"
	"gin-quickstart/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	log         *logger.Logger
	GormDB      *gorm.DB
	RedisClient *redis.Client
}

func NewNotificationRepository(log *logger.Logger, db *gorm.DB, redis *redis.Client) *NotificationRepository {
	return &NotificationRepository{
		log:         log,
		GormDB:      db,
		RedisClient: redis,
	}
}

func (r NotificationRepository) GetNotificationsByUserID(ctx *gin.Context, userID uint64) ([]model.Notification, error) {
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

func (r NotificationRepository) GetNotificationByID(ctx *gin.Context, id uint64) (*model.Notification, error) {
	var notification model.Notification
	err := r.GormDB.Preload("User").Preload("Thread").Preload("Post").First(&notification, id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *NotificationRepository) Create(ctx *gin.Context, notification *model.Notification) error {
	return r.GormDB.Create(notification).Error
}

func (r *NotificationRepository) Update(ctx *gin.Context, notification *model.Notification) error {
	return r.GormDB.Save(notification).Error
}

func (r *NotificationRepository) Delete(ctx *gin.Context, notification *model.Notification) error {
	return r.GormDB.Delete(notification).Error
}
