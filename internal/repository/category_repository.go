package repository

import (
	"gin-quickstart/internal/model"
	"gin-quickstart/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	log         *logger.Logger
	GormDB      *gorm.DB
	RedisClient *redis.Client
}

func NewCategoryRepository(log *logger.Logger, db *gorm.DB, redis *redis.Client) *CategoryRepository {
	return &CategoryRepository{
		log:         log,
		GormDB:      db,
		RedisClient: redis,
	}
}

// GETTER
func (r CategoryRepository) GetAllCategories(ctx *gin.Context) ([]model.Category, error) {
	var categories []model.Category
	err := r.GormDB.
		Preload("Categories").
		Preload("Threads").
		Find(&categories).Error

	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r CategoryRepository) GetCategoryByID(ctx *gin.Context, id uint64) (*model.Category, error) {
	var category model.Category
	err := r.GormDB.
		Preload("Categories").
		Preload("Threads").
		First(&category, id).Error

	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r CategoryRepository) GetCategoryBySlug(ctx *gin.Context, slug string) (*model.Category, error) {
	var category model.Category
	err := r.GormDB.
		Preload("Categories").
		Preload("Threads").
		Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// SETTER
func (r *CategoryRepository) Create(ctx *gin.Context, category *model.Category) error {
	return r.GormDB.Create(category).Error
}

func (r *CategoryRepository) Update(ctx *gin.Context, category *model.Category) error {
	return r.GormDB.Save(category).Error
}

func (r *CategoryRepository) Delete(ctx *gin.Context, category *model.Category) error {
	return r.GormDB.Delete(category).Error
}
