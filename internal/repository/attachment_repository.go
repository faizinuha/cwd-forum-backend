package repository

import (
	"context"
	"gin-quickstart/internal/model"
	"gin-quickstart/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var ctx = context.Background()

type AttachmentRepository struct {
	log         *logger.Logger
	GormDB      *gorm.DB
	RedisClient *redis.Client
}

func NewAttachmentRepository(logger *logger.Logger, db *gorm.DB, redis *redis.Client) *AttachmentRepository {
	return &AttachmentRepository{
		log:         logger,
		GormDB:      db,
		RedisClient: redis,
	}
}

// GETTER

func (r AttachmentRepository) GetAllAttachments(ctx *gin.Context) ([]*model.Attachment, error) {
	var attachments []*model.Attachment
	err := r.GormDB.Find(&attachments).Error
	if err != nil {
		return nil, err
	}

	return attachments, nil
}

func (r AttachmentRepository) GetAttachmentByID(ctx *gin.Context, id uint64) (*model.Attachment, error) {
	var attachment model.Attachment
	err := r.GormDB.First(&attachment, id).Error
	if err != nil {
		return nil, err
	}
	return &attachment, nil
}

func (r AttachmentRepository) GetAttachmentsByPostID(ctx *gin.Context, postID uint64) ([]*model.Attachment, error) {
	var attachments []*model.Attachment
	err := r.GormDB.Where("post_id = ?", postID).Find(&attachments).Error
	if err != nil {
		return nil, err
	}
	return attachments, nil
}

// SETTER
func (r *AttachmentRepository) Delete(ctx *gin.Context, attachment *model.Attachment) error {
	return r.GormDB.Delete(attachment).Error
}
