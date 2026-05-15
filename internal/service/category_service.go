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

type CategoryService struct {
	log *logger.Logger
	r   *repository.CategoryRepository
}

func NewCategoryService(log *logger.Logger, r *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		log: log,
		r:   r,
	}
}

// GETTER
func (s CategoryService) GetAllCategories(ctx *gin.Context) ([]model.Category, error) {
	getStatus := s.r.RedisClient.Get(ctx, "categories")

	if getStatus.Err() == nil {
		var categories []model.Category
		err := json.Unmarshal([]byte(getStatus.Val()), &categories)

		if err != nil {
			return nil, err
		}

		return categories, nil
	}

	if getStatus.Err() != nil && getStatus.Err() != redis.Nil {
		return nil, getStatus.Err()
	}

	categories, err := s.r.GetAllCategories(ctx)

	if err != nil {
		return nil, err
	}

	json, err := json.Marshal(categories)

	if err != nil {
		return nil, err
	}

	cmdStatus := s.r.RedisClient.Set(ctx, "categories", json, time.Hour)

	if cmdStatus.Err() != nil {
		return nil, cmdStatus.Err()
	}

	return categories, nil
}

func (s CategoryService) GetCategoryByID(ctx *gin.Context, id uint64) (*model.Category, error) {
	getStatus := s.r.RedisClient.Get(ctx, "category:id:"+strconv.FormatUint(id, 10))

	if getStatus.Err() == nil {
		var category model.Category
		err := json.Unmarshal([]byte(getStatus.Val()), &category)

		if err != nil {
			return nil, err
		}

		return &category, nil
	}

	if getStatus.Err() != nil && getStatus.Err() != redis.Nil {
		return nil, getStatus.Err()
	}

	category, err := s.r.GetCategoryByID(ctx, id)

	if err != nil {
		return nil, err
	}

	json, err := json.Marshal(category)

	if err != nil {
		return nil, err
	}

	cmdStatus := s.r.RedisClient.Set(ctx, "category:id:"+strconv.FormatUint(id, 10), json, time.Hour)

	if cmdStatus.Err() != nil {
		return nil, cmdStatus.Err()
	}

	return category, nil
}

func (s CategoryService) GetCategoryBySlug(ctx *gin.Context, slug string) (*model.Category, error) {
	getStatus := s.r.RedisClient.Get(ctx, "category:slug:"+slug)

	if getStatus.Err() == nil {
		var category model.Category
		err := json.Unmarshal([]byte(getStatus.Val()), &category)

		if err != nil {
			return nil, err
		}

		return &category, nil
	}

	if getStatus.Err() != nil && getStatus.Err() != redis.Nil {
		return nil, getStatus.Err()
	}

	category, err := s.r.GetCategoryBySlug(ctx, slug)

	if err != nil {
		return nil, err
	}

	json, err := json.Marshal(category)

	if err != nil {
		return nil, err
	}

	cmdStatus := s.r.RedisClient.Set(ctx, "category:slug:"+slug, json, time.Hour)

	if cmdStatus.Err() != nil {
		return nil, cmdStatus.Err()
	}

	return category, nil
}

// SETTER
func (s *CategoryService) Create(
	ctx *gin.Context,
	ParentID *uint,
	Name string,
	Slug string,
	Description string,
	IconUrl string,
	SortOrder int,
	IsPrivate bool,
) (*model.Category, error) {
	category := &model.Category{
		ParentID:    ParentID,
		Name:        Name,
		Slug:        Slug,
		Description: Description,
		IconUrl:     IconUrl,
		SortOrder:   SortOrder,
		IsPrivate:   IsPrivate,
	}

	if ParentID != nil {
		parentCategory, err := s.r.GetCategoryByID(ctx, uint64(*ParentID))

		if err != nil {
			return nil, errors.New("Parent category not found")
		}

		if parentCategory == nil {
			return nil, errors.New("Parent category not found")
		}
	}

	slugExists, _ := s.r.GetCategoryBySlug(ctx, Slug)

	if slugExists != nil {
		return nil, errors.New("Slug already exists")
	}

	err := s.r.Create(ctx, category)

	if err != nil {
		return nil, err
	}

	delStatus := s.r.RedisClient.Del(ctx, "categories")

	if delStatus.Err() != nil {
		return nil, delStatus.Err()
	}

	return category, nil
}

func (s *CategoryService) Update(
	ctx *gin.Context,
	ID uint64,
	ParentID *uint,
	Name *string,
	Slug *string,
	Description *string,
	IconUrl *string,
	SortOrder *int,
	IsPrivate *bool,
) (*model.Category, error) {
	category, err := s.r.GetCategoryByID(ctx, ID)

	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errors.New("Category not found")
	}

	if ParentID != nil {
		parentCategory, err := s.r.GetCategoryByID(ctx, uint64(*ParentID))

		if err != nil {
			return nil, errors.New("Parent category not found")
		}

		if parentCategory == nil {
			return nil, errors.New("Parent category not found")
		}

		if category.ID == parentCategory.ID {
			return nil, errors.New("Category cannot be its own parent")
		}
	}

	if Name != nil {
		category.Name = *Name
	}

	if Slug != nil {
		slugExists, _ := s.r.GetCategoryBySlug(ctx, *Slug)

		if slugExists != nil && slugExists.ID != category.ID {
			return nil, errors.New("Slug already exists")
		}

		category.Slug = *Slug
	}

	if Description != nil {
		category.Description = *Description
	}

	if IconUrl != nil {
		category.IconUrl = *IconUrl
	}

	if SortOrder != nil {
		category.SortOrder = *SortOrder
	}

	if IsPrivate != nil {
		category.IsPrivate = *IsPrivate
	}

	err = s.r.Update(ctx, category)

	if err != nil {
		return nil, err
	}

	delIdStatus := s.r.RedisClient.Del(ctx, "category:id:"+strconv.FormatUint(ID, 10))

	if delIdStatus.Err() != nil {
		return nil, delIdStatus.Err()
	}

	delSlugStatus := s.r.RedisClient.Del(ctx, "category:slug:"+category.Slug)

	if delSlugStatus.Err() != nil {
		return nil, delSlugStatus.Err()
	}

	return category, nil
}

func (s *CategoryService) Delete(ctx *gin.Context, ID uint64) error {
	category, err := s.r.GetCategoryByID(ctx, ID)

	if err != nil {
		return err
	}

	if category == nil {
		return errors.New("Category not found")
	}

	threads := category.Threads

	if len(threads) > 0 {
		return errors.New("Cannot delete category with existing threads")
	}

	subcategories := category.Categories

	if len(subcategories) > 0 {
		return errors.New("Cannot delete category with existing subcategories")
	}

	delErr := s.r.Delete(ctx, category)

	if delErr != nil {
		return delErr
	}

	delIdStatus := s.r.RedisClient.Del(ctx, "category:id:"+strconv.FormatUint(ID, 10))

	if delIdStatus.Err() != nil {
		return delIdStatus.Err()
	}

	delSlugStatus := s.r.RedisClient.Del(ctx, "category:slug:"+category.Slug)

	if delSlugStatus.Err() != nil {
		return delSlugStatus.Err()
	}

	return nil
}
