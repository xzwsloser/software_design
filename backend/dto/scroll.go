package dto

type ScrollRequest struct {
	PageIndex	 int32 	`json:"pageIndex"`
	PageSize 	 int32 	`json:"pageSize"`
}

type ScrollResp[T any] struct {
	Data  	[]T		`json:"data"`
	Total	int32 	`json:"total"` 
}
