package dto

import "github.com/xzwsloser/software_design/backend/model"

type BasicUserInfo struct {
	Id 			int32	`json:"id"`
	Username	string 	`json:"username"`
}

type UserInfoInRecSys struct {
	model.User
	Update		bool 	`json:"update"`
	Limit		int		`json:"limit"`
}
