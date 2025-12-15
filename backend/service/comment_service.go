package service

import (
	"github.com/xzwsloser/software_design/backend/model"
)

type CommentService struct {

}

func (*CommentService) QueryPositiveCommentByPage(siteIndex int32, pageIndex int32, pageSize int32) ([]model.CommentPositive, error) {
	offset := (pageIndex - 1) * pageSize
	limit  := pageSize

	c := &model.CommentPositive{}
	c.SiteIndex = siteIndex

	comments, err := c.QueryCommentsByPage(offset, limit)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (*CommentService) QueryNegativeCommentByPage(siteIndex int32, pageIndex int32, pageSize int32) ([]model.CommentNegative, error) {
	offset := (pageIndex - 1) * pageSize
	limit  := pageSize

	c := &model.CommentNegative{}
	c.SiteIndex = siteIndex

	comments, err := c.QueryCommentsByPage(offset, limit)
	if err != nil {
		return nil, err
	}

	return comments, err
}

func (*CommentService) CountPositiveComment(siteIndex int32) (int, error) {
	c := &model.CommentPositive{}
	c.SiteIndex = siteIndex

	count, err := c.CountPositiveComment()

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (*CommentService) CountNegativeComment(siteIndex int32) (int, error) {
	c := &model.CommentNegative{}
	c.SiteIndex = siteIndex

	count, err := c.CountNegativeComment()

	if err != nil {
		return 0, err
	}

	return int(count), nil
}




