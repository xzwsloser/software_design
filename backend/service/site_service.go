package service

import (
	"github.com/xzwsloser/software_design/backend/dto"
	"github.com/xzwsloser/software_design/backend/model"
	"github.com/xzwsloser/software_design/backend/utils"
)

type SiteService struct {

}

func (*SiteService) QueryByPageParams(pageParams *dto.ScrollRequest) ([]dto.SiteBasicInfo, error) {
	offset := (pageParams.PageIndex - 1)*pageParams.PageSize
	limit  := pageParams.PageSize

	s := &model.Site{}
	sites, err := s.QueryByPage(offset, limit)

	if err != nil {
		return nil, err
	}

	site_infos := make([]dto.SiteBasicInfo, 0, len(sites))

	for _, site := range sites {
		site_infos = append(site_infos, dto.SiteBasicInfo{
			Id: site.Id,
			Name: site.Name,
			Score: site.Score,
			HotDegree: site.HotDegree,
			Address: site.Address,
			Images: site.Images,
			SiteIndex: site.SiteIndex,
		})
	}

	return site_infos, nil
}

func (*SiteService) QueryByIndex(siteIndex int32) (model.Site, error) {
	site := &model.Site{}
	site.SiteIndex = siteIndex
	cur_site, err := site.QueryByIndex()
	return cur_site, err
}

func (*SiteService) QueryBySiteIndexList(siteIndexList []int) ([]dto.SiteBasicInfo, error) {
	site := &model.Site{}
	siteIndexListTrans := make([]int32, 0, len(siteIndexList))
	for _, siteIndex := range siteIndexList {
		siteIndexListTrans = append(siteIndexListTrans, int32(siteIndex))
	}

	sites, err := site.QueryBySiteIndexes(siteIndexListTrans)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return nil, err
	}

	site_infos := make([]dto.SiteBasicInfo, 0, len(sites))

	for _, site := range sites {
		site_infos = append(site_infos, dto.SiteBasicInfo{
			Id: site.Id,
			Name: site.Name,
			Score: site.Score,
			HotDegree: site.HotDegree,
			Address: site.Address,
			Images: site.Images,
			SiteIndex: site.SiteIndex,
		})
	}

	return site_infos, nil
}



