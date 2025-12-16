package cache

import (
	"fmt"
	"strconv"

	"github.com/xzwsloser/software_design/backend/utils"
)

type ViewCacheService struct {

}

const (
	VIEW_USER_TO_SITE_PREFIX = "view:user:site:"
	VIEW_SITE_TO_USER_PREFIX = "view:site:user:"
)

// @Description: 记录一条浏览记录
func (*ViewCacheService) View(userId int, siteIndex int) error {
	userKey := fmt.Sprintf("%s%d", VIEW_USER_TO_SITE_PREFIX, userId)
	siteKey := fmt.Sprintf("%s%d", VIEW_SITE_TO_USER_PREFIX, siteIndex)

	err := redisHandler.Client.SAdd(redisHandler.Ctx, userKey, siteIndex).Err()
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return err
	}

	err = redisHandler.Client.SAdd(redisHandler.Ctx, siteKey, userId).Err()
	if err != nil {
		utils.GetLogger().Error(err.Error())
	}

	return err
}


// @Description: 查询一个用户所有的浏览记录
func (*ViewCacheService) QueryViewOfUser(userId int) ([]int, error) {
	userKey := fmt.Sprintf("%s%d", VIEW_USER_TO_SITE_PREFIX, userId)

	resultStrs, err := redisHandler.Client.SMembers(redisHandler.Ctx, userKey).Result()
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return nil, err
	}

	results := make([]int, 0, len(resultStrs))
	for _, resultStr := range resultStrs {
		result, err := strconv.Atoi(resultStr)
		if err != nil {
			utils.GetLogger().Error(err.Error())
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

// @Description: 查询一个景点所有的浏览用户
func (*ViewCacheService) QueryViewOfSite(siteIndex int) ([]int, error) {
	siteKey := fmt.Sprintf("%s%d", VIEW_SITE_TO_USER_PREFIX, siteIndex)

	resultStrs, err := redisHandler.Client.SMembers(redisHandler.Ctx, siteKey).Result()
	if err != nil  {
		utils.GetLogger().Error(err.Error())
		return nil, err
	}

	results := make([]int, 0, len(resultStrs))	
	for _, resultStr := range resultStrs {
		result, err := strconv.Atoi(resultStr)
		if err != nil {
			utils.GetLogger().Error(err.Error())
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

// @Description: 查询某一个景点是否被用户浏览
func (*ViewCacheService) QueryIsViewByUser(userId int, siteIndex int) (bool , error) {
	userKey := fmt.Sprintf("%s%d", VIEW_USER_TO_SITE_PREFIX, userId)
	result, err := redisHandler.Client.SIsMember(redisHandler.Ctx, userKey, siteIndex).Result()
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return false, err
	}

	return result, nil
}

// @Description: 查询某一个用户是否浏览某一个景点
func (*ViewCacheService) QueryIsViewSite(userId int, siteIndex int) (bool , error) {
	siteKey := fmt.Sprintf("%s%d", VIEW_SITE_TO_USER_PREFIX, siteIndex)
	result, err := redisHandler.Client.SIsMember(redisHandler.Ctx, siteKey, userId).Result()
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return false, err
	}

	return result, nil
}


