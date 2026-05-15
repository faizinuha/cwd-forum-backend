package service

import (
	"encoding/json"
	"errors"
	"gin-quickstart/internal/model"
	"gin-quickstart/internal/repository"
	"gin-quickstart/pkg/logger"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type TagService struct {
	log *logger.Logger
	r   *repository.TagRepository
}

func NewTagService(log *logger.Logger, r *repository.TagRepository) *TagService {
	return &TagService{
		log: log,
		r:   r,
	}
}

// GETTER
func (s TagService) GetAllTags(ctx *gin.Context) ([]model.Tag, error) {
	getStatus := s.r.RedisClient.Get(ctx, "tags")

	if getStatus.Err() == nil {
		var tags []model.Tag
		err := json.Unmarshal([]byte(getStatus.Val()), &tags)

		if err != nil {
			return nil, err
		}

		return tags, nil
	}

	if getStatus.Err() != nil && getStatus.Err() != redis.Nil {
		return nil, getStatus.Err()
	}

	tags, err := s.r.GetAllTags(ctx)

	if err != nil {
		return nil, err
	}

	json, err := json.Marshal(tags)

	if err != nil {
		return nil, err
	}

	cmdStatus := s.r.RedisClient.Set(ctx, "tags", json, time.Hour)

	if cmdStatus.Err() != nil {
		return nil, cmdStatus.Err()
	}

	return tags, nil
}

func (s TagService) GetTagByID(ctx *gin.Context, id uint64) (*model.Tag, error) {
	getStatus := s.r.RedisClient.Get(ctx, "tag:"+strconv.FormatUint(id, 10))

	if getStatus.Err() == nil {
		var tag model.Tag
		err := json.Unmarshal([]byte(getStatus.Val()), &tag)

		if err != nil {
			return nil, err
		}

		return &tag, nil
	}

	if getStatus.Err() != nil && getStatus.Err() != redis.Nil {
		return nil, getStatus.Err()
	}

	tag, err := s.r.GetTagByID(ctx, id)

	if err != nil {
		return nil, err
	}

	json, err := json.Marshal(tag)

	if err != nil {
		return nil, err
	}

	cmdStatus := s.r.RedisClient.Set(ctx, "tag:"+strconv.FormatUint(id, 10), json, time.Hour)

	if cmdStatus.Err() != nil {
		return nil, cmdStatus.Err()
	}

	return tag, nil
}

func (s TagService) GetTagBySlug(ctx *gin.Context, slug string) (*model.Tag, error) {
	getStatus := s.r.RedisClient.Get(ctx, "tag:slug:"+slug)

	if getStatus.Err() == nil {
		var tag model.Tag
		err := json.Unmarshal([]byte(getStatus.Val()), &tag)

		if err != nil {
			return nil, err
		}

		return &tag, nil
	}

	if getStatus.Err() != nil && getStatus.Err() != redis.Nil {
		return nil, getStatus.Err()
	}

	tag, err := s.r.GetTagBySlug(ctx, slug)

	if err != nil {
		return nil, err
	}

	json, err := json.Marshal(tag)

	if err != nil {
		return nil, err
	}

	cmdStatus := s.r.RedisClient.Set(ctx, "tag:slug:"+slug, json, time.Hour)

	if cmdStatus.Err() != nil {
		return nil, cmdStatus.Err()
	}

	return tag, nil
}

// SETTER
func (s *TagService) Create(
	ctx *gin.Context,
	Name string,
	Slug string,
	Color string,
) (*model.Tag, error) {
	tag := &model.Tag{
		Name:  Name,
		Slug:  Slug,
		Color: Color,
	}

	slugExists, _ := s.r.GetTagBySlug(ctx, Slug)
	if slugExists != nil {
		return nil, errors.New("tag with the same slug already exists")
	}

	err := s.r.Create(ctx, tag)
	if err != nil {
		return nil, err
	}

	delStatus := s.r.RedisClient.Del(ctx, "tags")

	if delStatus.Err() != nil {
		return nil, delStatus.Err()
	}

	return tag, nil
}

func (s *TagService) Update(
	ctx *gin.Context,
	ID uint64,
	Name *string,
	Slug *string,
	Color *string,
) (*model.Tag, error) {
	tag, err := s.GetTagByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	if tag == nil {
		return nil, errors.New("tag not found")
	}

	if Name != nil {
		tag.Name = *Name
	}
	if Slug != nil {
		slugExists, _ := s.r.GetTagBySlug(ctx, *Slug)
		if slugExists != nil && slugExists.ID != uint(ID) {
			return nil, errors.New("tag with the same slug already exists")
		}
		tag.Slug = *Slug
	}
	if Color != nil {
		tag.Color = *Color
	}

	err = s.r.Update(ctx, tag)
	if err != nil {
		return nil, err
	}

	delStatus := s.r.RedisClient.Del(ctx, "tags", "tag:"+strconv.FormatUint(ID, 10), "tag:slug:"+tag.Slug)

	if delStatus.Err() != nil {
		return nil, delStatus.Err()
	}

	return tag, nil
}

func (s *TagService) Delete(ctx *gin.Context, id uint64) error {
	tag, err := s.r.GetTagByID(ctx, id)
	if err != nil {
		return err
	}

	if tag == nil {
		return errors.New("tag not found")
	}

	pruneErr := s.r.GormDB.Model(tag).Association("Threads").Clear()

	if pruneErr != nil {
		return pruneErr
	}

	err = s.r.Delete(ctx, id)
	if err != nil {
		return err
	}

	delStatus := s.r.RedisClient.Del(ctx, "tags", "tag:"+strconv.FormatUint(id, 10), "tag:slug:"+tag.Slug)

	if delStatus.Err() != nil {
		return delStatus.Err()
	}

	return nil
}
