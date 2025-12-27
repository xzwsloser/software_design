package cache

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/xzwsloser/software_design/backend/utils"
)

type RecSysCacheService struct {

}

const (
	REC_SYS_USER_KEY_PREFIX = "rec:user:"	
)

func (*RecSysCacheService) StoreRecSiteIdList(userId int, recSiteIdxList []int) error {
	key := fmt.Sprintf("%s%d", REC_SYS_USER_KEY_PREFIX, userId)
	siteStrToStore := utils.ParseFromArrayToStr(recSiteIdxList)

	err := redisHandler.Client.Set(redisHandler.Ctx, key, siteStrToStore, -1).Err()
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return err
	}

	return nil
}

func (*RecSysCacheService) GetRecSiteIdList(userId int) ([]int, error) {
	key := fmt.Sprintf("%s%d", REC_SYS_USER_KEY_PREFIX, userId)

	siteIdxsStr, err := redisHandler.Client.Get(redisHandler.Ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			utils.GetLogger().Error(err.Error())
			return nil, err
		}

		return []int{}, nil
	}

	siteIdxList, err := utils.ParseFromStrToArray(siteIdxsStr)

	if err != nil {
		utils.GetLogger().Error(err.Error())
		return nil, err
	}

	return siteIdxList, nil
}


func (*RecSysCacheService) DelRecResult(userId int) (error) {
	key := fmt.Sprintf("%s%d", REC_SYS_USER_KEY_PREFIX, userId)

	err := redisHandler.Client.Del(redisHandler.Ctx, key).Err()
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return err
	}

	return nil
}
