package service

import (
	"github.com/xzwsloser/software_design/backend/cache"
	"github.com/xzwsloser/software_design/backend/dto"
	"github.com/xzwsloser/software_design/backend/model"
	"github.com/xzwsloser/software_design/backend/rpc"
	"github.com/xzwsloser/software_design/backend/utils"
)

type RecSysService struct {
	userToRecChan		chan *dto.UserInfoInRecSys
	maxConcurrentNum	int
}

var (
	recSysCacheService = new(cache.RecSysCacheService)
	recSysService *RecSysService = nil
)

func NewRecSysPipeline() {
	recSysService = &RecSysService{
		userToRecChan: make(chan *dto.UserInfoInRecSys, 10),
		// 最大并发数量
		maxConcurrentNum: 3,
	}

	for i := 0 ; i < recSysService.maxConcurrentNum ; i ++ {
		go recSysService.SyncGetRecResult()
	}
}

func AddRecTaskToPipeline(info *dto.UserInfoInRecSys) {
	recSysService.userToRecChan <- info
}


func (r *RecSysService) SyncGetRecResult() {
	for {
		user := <- r.userToRecChan

		// 获取推荐结果
		userId := int(user.Id)
		addressId := user.AddressId
		touristType := user.TouristType
		priceSensitive := user.PriceSensitive
		update := user.Update
		limit := user.Limit

		likeType, err := utils.ParseFromStrToArray(user.LikeType)
		targetType, err := utils.ParseFromStrToArray(user.Targets)
		attentionType, err := utils.ParseFromStrToArray(user.Attention)

		if err != nil {
			utils.GetLogger().Error(err.Error())
			continue
		}

		// 获取用户喜欢景点
		likedSiteIdxList, err := likeCacheService.QueryLikeOfUser(userId)
		if err != nil {
			utils.GetLogger().Error(err.Error())
			continue
		}

		recResult, err := rpc.GetRecSysClient().GetRecResult(
			userId,
			addressId,
			touristType,
			priceSensitive,
			likeType,
			targetType,
			attentionType,
			update,
			limit,
			likedSiteIdxList,
		)

		// 存入 redis
		err = recSysCacheService.StoreRecSiteIdList(userId, recResult)
		if err != nil {
			utils.GetLogger().Error(err.Error())
			continue
		}

	}
}

func (r *RecSysService) Stop() {
	close(r.userToRecChan)
}

func (*RecSysService) QueryRecommandSiteIdxs(userId int) ([]int, error) {
	siteIdxList, err := recSysCacheService.GetRecSiteIdList(userId)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return nil, err
	}

	if len(siteIdxList) == 0 {
		// 没有更新到 Redis
		utils.GetLogger().Info("Begin to update recommand info in redis")
		u := &model.User{
			Id: int32(userId),
		}

		matchedUser, err := u.QueryByUserId()
		if err != nil {
			utils.GetLogger().Error(err.Error())
			return nil, err
		}

		// 手动更新
		addressId := matchedUser.AddressId
		touristType := matchedUser.TouristType
		priceSensitive := matchedUser.PriceSensitive
		update := true
		limit := 200

		likeType, err := utils.ParseFromStrToArray(matchedUser.LikeType)
		targetType, err := utils.ParseFromStrToArray(matchedUser.Targets)
		attentionType, err := utils.ParseFromStrToArray(matchedUser.Attention)

		if err != nil {
			utils.GetLogger().Error(err.Error())
			return nil, err
		}

		likedSiteIdxList, err := likeCacheService.QueryLikeOfUser(userId)
		if err != nil {
			utils.GetLogger().Error(err.Error())
			return nil, err
		}

		recResult, err := rpc.GetRecSysClient().GetRecResult(
			userId,
			addressId,
			touristType,
			priceSensitive,
			likeType,
			targetType,
			attentionType,
			update,
			limit,
			likedSiteIdxList,
		)

		// 存入 redis
		err = recSysCacheService.StoreRecSiteIdList(userId, recResult)
		if err != nil {
			utils.GetLogger().Error(err.Error())
			return nil, err
		}

		return recResult, nil
	}

	return siteIdxList, nil
}
