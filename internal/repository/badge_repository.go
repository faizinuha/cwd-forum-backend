package repository

import (
	"gin-quickstart/internal/model"
	"gin-quickstart/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BadgeRepository struct {
	log         *logger.Logger
	GormDB      *gorm.DB
	RedisClient *redis.Client
}

func NewBadgeRepository(log *logger.Logger, db *gorm.DB, redis *redis.Client) *BadgeRepository {
	return &BadgeRepository{
		log:         log,
		GormDB:      db,
		RedisClient: redis,
	}
}

// GETTER
func (r BadgeRepository) GetAllBadges(ctx *gin.Context) ([]*model.Badge, error) {
	var badges []*model.Badge
	err := r.GormDB.Find(&badges).Error
	if err != nil {
		return nil, err
	}
	return badges, nil
}

func (r BadgeRepository) GetBadgeByID(ctx *gin.Context, id uint64) (*model.Badge, error) {
	var badge model.Badge
	err := r.GormDB.First(&badge, id).Error
	if err != nil {
		return nil, err
	}
	return &badge, nil
}

// SETTER
func (r *BadgeRepository) Create(ctx *gin.Context, badge *model.Badge) error {
	return r.GormDB.Create(badge).Error
}

func (r *BadgeRepository) Update(ctx *gin.Context, badge *model.Badge) error {
	return r.GormDB.Save(badge).Error
}

func (r *BadgeRepository) Delete(ctx *gin.Context, badge *model.Badge) error {
	return r.GormDB.Delete(badge).Error
}
