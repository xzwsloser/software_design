package service

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/xzwsloser/software_design/backend/cache"
	"github.com/xzwsloser/software_design/backend/model"
	"github.com/xzwsloser/software_design/backend/utils"
)

type LikeService struct {

}

var (
	likeCacheService = new(cache.LikeCacheService)
)

var (
	SITE_LIKED_ERROR = errors.New("The Site is Liked By User")
)

func (*LikeService) Like(userId int, siteIndex int) error {
	isLiked, err := likeCacheService.QuerySiteIsLikedByUser(userId, siteIndex)
	if err != nil {
		utils.GetLogger().Error("Failed to query site is liked by user")
		return err
	}

	if isLiked {
		utils.GetLogger().Error("The Site is Liked By User")
		return SITE_LIKED_ERROR
	}

	// 先写 mysql
	l := &model.Like{}
	l.UserId = userId
	l.SiteIndex = siteIndex
	_, err = l.QueryConnection()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = l.CreateRecord()
		if err != nil {
			utils.GetLogger().Error(err.Error())
			return err
		}
	} else {
		err = l.UpdateLikeStatus(1)
		if err != nil {
			utils.GetLogger().Error(err.Error())
			return err
		}
	}

	err = likeCacheService.Like(userId, siteIndex)
	if err != nil {
		utils.GetLogger().Error(err.Error())
	}

	return err
}

func (*LikeService) CancelLike(userId int, siteIndex int) error {
	isLike, err := likeCacheService.QuerySiteIsLikedByUser(userId, siteIndex)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return err
	}

	if !isLike {
		return nil
	}

	l := &model.Like{}
	l.UserId = userId
	l.SiteIndex = siteIndex
	err = l.UpdateLikeStatus(0)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return err
	}

	err = likeCacheService.CancelLike(userId, siteIndex)
	if err != nil {
		utils.GetLogger().Error(err.Error())
	}

	return err
}

func (*LikeService) QueryLikeOfUser(userId int) ([]int, error) {
	// result, err := likeCacheService.QueryLikeOfUser(userId)
	// if err != nil {
	// 	utils.GetLogger().Error(err.Error())
	// 	return nil, err
	// }

	// if len(result) != 0 {
	// 	return result, nil
	// }

	// 查询 mysql
	like := &model.Like{}
	like.UserId = userId
	likeLists, err := like.QueryLikeList()
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return nil, err
	}

	result := make([]int, 0, len(likeLists))
	for _, relation := range likeLists {
		result = append(result, relation.SiteIndex)
	}

	return result, nil
}

func (*LikeService) QueryLikeOfSite(siteIndex int) ([]int, error) {
	// result, err := likeCacheService.QueryLikeOfSite(siteIndex)
	// if err != nil {
	// 	utils.GetLogger().Error(err.Error())
	// 	return nil, err
	// }

	// if len(result) != 0 {
	// 	return result, nil
	// }

	// 查询 mysql
	like := &model.Like{}
	like.SiteIndex = siteIndex
	likeLists, err := like.QueryLikeUserList()
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return nil, err
	}

	result := make([]int, 0, len(likeLists))
	for _, relation := range likeLists {
		result = append(result, 0, relation.UserId)
	}

	return result, nil
}

func (*LikeService) QueryIsLikedByUser(userId int, siteIndex int) (bool , error) {
	// result, err := likeCacheService.QuerySiteIsLikedByUser(userId, siteIndex)
	// if err != nil {
	// 	utils.GetLogger().Error(err.Error())
	// 	return false, err
	// }

	// if result { return result, nil }

	like := &model.Like{}
	like.UserId = userId
	like.SiteIndex = siteIndex
	r, err := like.QueryConnection()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		utils.GetLogger().Error(err.Error())
		return false, err
	}


	return r.Valid == 1, nil
}

