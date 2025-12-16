package cache

import (
	"fmt"
	"strconv"

	"github.com/xzwsloser/software_design/backend/utils"
)

const (
	LIKE_USER_TO_SITE_PREFIX = "like:user:site:"
	LIKE_SITE_TO_USER_PREFIX = "like:site:user:"
)

type LikeCacheService struct {

}

// @Description: 存储点赞信息到 Redis 中
func (*LikeCacheService) Like(userId int, siteIndex int) (error) {
	userKey := fmt.Sprintf("%s%d", LIKE_USER_TO_SITE_PREFIX, userId)
	err := redisHandler.Client.SAdd(redisHandler.Ctx, 
							        userKey,
							  	 	siteIndex).Err()
	
	if err != nil {
		utils.GetLogger().Error(err)
		return err
	}

	siteKey := fmt.Sprintf("%s%d", LIKE_SITE_TO_USER_PREFIX, siteIndex)
	err = redisHandler.Client.SAdd(redisHandler.Ctx,
								   siteKey,
								   userId).Err()
	return err
}

// @Description: 查询对应用户的点赞记录
func (*LikeCacheService) QueryLikeOfUser(userId int) ([]int, error) {
	userKey := fmt.Sprintf("%s%d", LIKE_USER_TO_SITE_PREFIX, userId)
	resultStrs, err := redisHandler.Client.SMembers(redisHandler.Ctx, userKey).Result()
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return nil, err
	}

	results := make([]int, 0, len(resultStrs))
	for _, resultStr := range resultStrs {
		result, err := strconv.Atoi(resultStr)
		if err != nil {
			utils.GetLogger().Error()
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

// @Description: 查询对应的点赞用户
func (*LikeCacheService) QueryLikeOfSite(siteIndex int) ([]int, error) {
	siteKey := fmt.Sprintf("%s%d", LIKE_SITE_TO_USER_PREFIX, siteIndex)
	resultStrs, err := redisHandler.Client.SMembers(redisHandler.Ctx, siteKey).Result()
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return nil, err
	}

	results := make([]int, 0, len(resultStrs))
	for _, resultStr := range resultStrs {
		result, err := strconv.Atoi(resultStr)
		if err != nil {
			utils.GetLogger().Error()
			return nil, err
		}

		results = append(results, result)
	}

	return results, err
}

// @Description: 查询某一个景点是否被用户点赞
func (*LikeCacheService) QuerySiteIsLikedByUser(userId int, siteIndex int) (bool, error) {
	userKey := fmt.Sprintf("%s%d", LIKE_USER_TO_SITE_PREFIX, userId)
	result, err := redisHandler.Client.SIsMember(redisHandler.Ctx, userKey, siteIndex).Result()
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return false, err
	}

	return result, nil
}

// @Description: 取消点赞
func (*LikeCacheService) CancelLike(userId int, siteIndex int) error {
	userKey := fmt.Sprintf("%s%d", LIKE_USER_TO_SITE_PREFIX, userId)
	siteKey := fmt.Sprintf("%s%d", LIKE_SITE_TO_USER_PREFIX, siteIndex)

	err := redisHandler.Client.SRem(redisHandler.Ctx, userKey, siteIndex).Err()
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return err
	}

	err = redisHandler.Client.SRem(redisHandler.Ctx, siteKey, userId).Err()
	if err != nil {
		utils.GetLogger().Error(err.Error())
	}

	return err
}
