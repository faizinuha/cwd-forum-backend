package repository

import (
	"gin-quickstart/internal/model"
	"gin-quickstart/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ThreadRepository struct {
	log         *logger.Logger
	GormDB      *gorm.DB
	RedisClient *redis.Client
}

func NewThreadRepository(log *logger.Logger, db *gorm.DB, redis *redis.Client) *ThreadRepository {
	return &ThreadRepository{
		log:         log,
		GormDB:      db,
		RedisClient: redis,
	}
}

// GETTER
func (r ThreadRepository) GetAllThreads(ctx *gin.Context) ([]model.Thread, error) {
	var threads []model.Thread
	err := r.GormDB.
		Preload("Category").
		Preload("Posts").
		Preload("Tags").
		Preload("Author").
		Preload("PinnedByUser").
		Find(&threads).Error

	if err != nil {
		return nil, err
	}
	return threads, nil
}

func (r ThreadRepository) GetThreadByID(ctx *gin.Context, id uint64) (*model.Thread, error) {
	var thread model.Thread
	err := r.GormDB.
		Preload("Category").
		Preload("Posts").
		Preload("Tags").
		Preload("Author").
		Preload("PinnedByUser").
		First(&thread, id).
		Error

	if err != nil {
		return nil, err
	}
	return &thread, nil
}

func (r ThreadRepository) GetThreadBySlug(ctx *gin.Context, slug string) (*model.Thread, error) {
	var thread model.Thread
	err := r.GormDB.Where("slug = ?", slug).First(&thread).Error
	if err != nil {
		return nil, err
	}
	return &thread, nil
}

func (r ThreadRepository) GetThreadsByCategoryID(ctx *gin.Context, categoryID uint) ([]model.Thread, error) {
	var threads []model.Thread
	err := r.GormDB.Where("category_id = ?", categoryID).Find(&threads).Error
	if err != nil {
		return nil, err
	}
	return threads, nil
}

func (r ThreadRepository) GetThreadsByAuthorID(ctx *gin.Context, authorID uint) ([]model.Thread, error) {
	var threads []model.Thread
	err := r.GormDB.Where("author_id = ?", authorID).Find(&threads).Error
	if err != nil {
		return nil, err
	}
	return threads, nil
}

func (r ThreadRepository) GetThreadsByTagID(ctx *gin.Context, tagID uint) ([]model.Thread, error) {
	var threads []model.Thread
	err := r.GormDB.
		Joins("JOIN thread_tags ON thread_tags.thread_id = threads.id").
		Where("thread_tags.tag_id = ?", tagID).
		Find(&threads).Error

	if err != nil {
		return nil, err
	}
	return threads, nil
}

// SETTER
func (r *ThreadRepository) Create(ctx *gin.Context, thread *model.Thread) error {
	return r.GormDB.Create(thread).Error
}

func (r *ThreadRepository) Update(ctx *gin.Context, thread *model.Thread) error {
	return r.GormDB.Save(thread).Error
}

func (r *ThreadRepository) Delete(ctx *gin.Context, thread *model.Thread) error {
	return r.GormDB.Delete(thread).Error
}

func (r *ThreadRepository) CreatePostAttachment(
	ctx *gin.Context,
	post *model.Post,
	attachment *model.Attachment,
) error {
	if err := r.GormDB.Create(attachment).Error; err != nil {
		return err
	}

	if err := r.GormDB.Save(post).Error; err != nil {
		return err
	}

	if err := r.GormDB.Save(attachment).Error; err != nil {
		return err
	}

	return nil
}
