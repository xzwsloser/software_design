package dto

import "github.com/xzwsloser/software_design/backend/model"

type SiteBasicInfo struct {
	Id			int32		`json:"id"`
	Name 		string 		`json:"name"`
	Score 		float32		`json:"score"`
	HotDegree	float32		`json:"hogDegree"`
	Images		string		`json:"images"`
	Address 	string		`json:"address"`
	SiteIndex   int32		`json:"siteIndex"`
}

type SiteDetailInfo struct {
	model.Site
	PositiveCommentWCPic	string		`json:"positiveCommentPic"`
	NegativeCommentWCPic	string		`json:"negativeCommentPic"`
	TouristTypePiePic		string		`json:"touristTypePiePic"`
}




