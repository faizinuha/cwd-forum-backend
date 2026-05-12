package service

import (
	"gin-quickstart/internal/model"
	"gin-quickstart/internal/repository"
)

type AttachmentService struct {
	r *repository.AttachmentRepository
}

func NewAttachmentService(r *repository.AttachmentRepository) *AttachmentService {
	return &AttachmentService{
		r: r,
	}
}

// GETTER
func (s AttachmentService) GetAllAttachments() ([]*model.Attachment, error) {
	return s.r.GetAllAttachments()
}

func (s AttachmentService) GetAttachmentByID(id uint64) (*model.Attachment, error) {
	return s.r.GetAttachmentByID(id)
}

func (s AttachmentService) GetAttachmentsByPostID(postID uint64) ([]*model.Attachment, error) {
	return s.r.GetAttachmentsByPostID(postID)
}

// SETTER
func (s *AttachmentService) Delete(attachment *model.Attachment) error {
	return s.r.Delete(attachment)
}
