package service

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/xzwsloser/software_design/backend/cache"
	"github.com/xzwsloser/software_design/backend/model"
	"github.com/xzwsloser/software_design/backend/utils"
)

type ViewService struct {

}

var (
	viewCacheService = new(cache.ViewCacheService)
)

// @Description: 记录景点浏览记录(返回值 -> 是否插入成功)
func (*ViewService) View(userId int, siteIndex int) (bool, error) {
	isVisited, err := viewCacheService.QueryIsViewByUser(userId, siteIndex)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return false, err
	}

	if (isVisited) {
		utils.GetLogger().Debugf("User: %d has visited the Site %d", userId, siteIndex)	
		return false, nil
	}

	v := &model.View{}
	v.UserId = userId
	v.SiteIndex = siteIndex
	_, err = v.QueryViewConnection()

	if err == nil {
		return false, nil
	} else {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.GetLogger().Error(err.Error())
			return false, err
		}
	}

	err = v.CreateViewRecord()
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return false, err
	}

	err = viewCacheService.View(userId, siteIndex)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return false, err
	}

	return true, nil
}


// @Description: 查询某一个用户的景点浏览列表
func (*ViewService) QueryVisitedSiteList(userId int) ([]int, error) {
	visitedSiteList, err := viewCacheService.QueryViewOfUser(userId)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return nil, err
	}

	if len(visitedSiteList) != 0 {
		return visitedSiteList, nil
	}

	view := &model.View{}
	view.UserId = userId
	viewList, err := view.QueryViewList()
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return nil, err
	}

	visitedSiteList = make([]int, 0, len(viewList))
	for _, cur_view := range viewList {
		visitedSiteList = append(visitedSiteList, cur_view.SiteIndex)
	}

	return visitedSiteList, nil
}

// @Description: 查询浏览过某一个景点的用户列表
func (*ViewService) QueryUserListed(siteIndex int) ([]int, error) {
	userList, err := viewCacheService.QueryViewOfSite(siteIndex)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return nil, err
	}

	if len(userList) != 0 {
		return userList, nil
	}

	view := &model.View{}
	view.SiteIndex = siteIndex
	viewList, err := view.QueryViewUserList()
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return nil, err
	}

	userList = make([]int, 0, len(userList))
	for _, cur_view := range viewList {
		userList = append(userList, cur_view.UserId)
	}

	return userList, nil
}

