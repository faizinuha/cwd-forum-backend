package repository

import (
	"gin-quickstart/internal/model"
	"gin-quickstart/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PostRepository struct {
	log         *logger.Logger
	GormDB      *gorm.DB
	RedisClient *redis.Client
}

func NewPostRepository(log *logger.Logger, db *gorm.DB, redis *redis.Client) *PostRepository {
	return &PostRepository{
		log:         log,
		GormDB:      db,
		RedisClient: redis,
	}
}

// GETTER
func (r PostRepository) GetAllPosts(ctx *gin.Context) ([]model.Post, error) {
	var posts []model.Post
	err := r.GormDB.
		Preload("Thread").
		Preload("Author").
		Find(&posts).Error

	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r PostRepository) GetPostByID(ctx *gin.Context, id uint64) (*model.Post, error) {
	var post model.Post
	err := r.GormDB.
		Preload("Thread").
		Preload("Author").
		Preload("Posts").
		First(&post, id).Error

	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r PostRepository) GetPostsByThreadID(ctx *gin.Context, threadID uint64) ([]model.Post, error) {
	var posts []model.Post
	err := r.GormDB.
		Preload("Thread").
		Preload("Author").
		Where("thread_id = ?", threadID).Find(&posts).Error

	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r PostRepository) GetPostsByAuthorID(ctx *gin.Context, authorID uint64) ([]model.Post, error) {
	var posts []model.Post
	err := r.GormDB.
		Preload("Thread").
		Preload("Author").
		Where("author_id = ?", authorID).Find(&posts).Error

	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r PostRepository) GetPostsByParentID(ctx *gin.Context, parentID uint64) ([]model.Post, error) {
	var posts []model.Post
	err := r.GormDB.
		Preload("Thread").
		Preload("Author").
		Where("parent_id = ?", parentID).Find(&posts).Error

	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r PostRepository) GetPostVotes(ctx *gin.Context, postID uint64) ([]model.Vote, error) {
	var votes []model.Vote
	err := r.GormDB.Where("post_id = ?", postID).Find(&votes).Error

	if err != nil {
		return nil, err
	}
	return votes, nil
}

// SETTER

func (r *PostRepository) Create(ctx *gin.Context, post *model.Post) error {
	return r.GormDB.Create(post).Error
}

func (r *PostRepository) Update(ctx *gin.Context, post *model.Post) error {
	post.IsEdited = true
	return r.GormDB.Save(post).Error
}

func (r *PostRepository) Delete(ctx *gin.Context, post *model.Post) error {
	return r.GormDB.Delete(post).Error
}

func (r *PostRepository) CreateAttachment(ctx *gin.Context, postID uint64,
	attachment *model.Attachment) (*model.Attachment, error) {

	r.GormDB.Model(&model.Post{ID: uint(postID)}).Association("Attachments").Append(attachment)

	return attachment, nil
}
