package dto

type SiteBasicInfo struct {
	Id			int32		`json:"id"`
	Name 		string 		`json:"name"`
	Score 		float32		`json:"score"`
	HotDegree	float32		`json:"hogDegree"`
	Images		string		`json:"images"`
	Address 	string		`json:"address"`
	SiteIndex   int32		`json:"siteIndex"`
}

