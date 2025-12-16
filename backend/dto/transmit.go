package dto

type Connection struct {
	UserId 		int		`json:"id"`
	SiteIndex 	int		`json:"siteIndex"`
}

type SiteQueryReq struct {
	SiteIndexList 	[]int 	`json:"siteIndexList"`
}




