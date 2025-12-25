package service

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/xzwsloser/software_design/backend/dto"
	"github.com/xzwsloser/software_design/backend/middleware"
	"github.com/xzwsloser/software_design/backend/model"
	"github.com/xzwsloser/software_design/backend/utils"
)

type UserService struct {
}

var (
	UserRecordNotFoundErr = errors.New("User Record Not Find")
	UserPasswordErr       = errors.New("Password Error")
	UserTokenErr          = errors.New("Generate Jwt Token Error")
	UserExistsErr         = errors.New("User Exists")
)

func (*UserService) Login(user *model.User) (string, error) {
	matchedUser, err := user.QueryByUsername()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", UserRecordNotFoundErr
		}
		return "", err
	}

	if matchedUser.Password != user.Password {
		return "", UserPasswordErr
	}

	// 重新生成 jwt token
	user.Id = matchedUser.Id
	jwt := middleware.NewJwt()
	basicUserInfo := dto.BasicUserInfo{
		Id:       user.Id,
		Username: user.Username,
	}

	clamis := jwt.CreatClaims(basicUserInfo)
	jwtToken, err := jwt.CreateJwtToken(clamis)
	if err != nil {
		return "", UserTokenErr
	}

	return jwtToken, nil
}

func (*UserService) Register(user *model.User) (string, error) {
	_, err := user.QueryByUsername()
	if err == nil {
		// 用户已经存在
		return "", UserExistsErr
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	err = user.InsertUser()
	if err != nil {
		return "", err
	}

	jwt := middleware.NewJwt()
	basicUserInfo := dto.BasicUserInfo{
		Id:       user.Id,
		Username: user.Username,
	}
	clamis := jwt.CreatClaims(basicUserInfo)
	jwtToken, err := jwt.CreateJwtToken(clamis)
	if err != nil {
		return "", UserTokenErr
	}

	return jwtToken, nil
}

func (*UserService) GetCurrentUserInfo(username string) (model.User, error) {
	u := &model.User{}
	u.Username = username
	matchedUser, err := u.QueryByUsername()
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return model.User{}, err
	}

	return matchedUser, err
}

func (*UserService) UpdateUserInfo(user *model.User) (error) {
	err := user.UpdateUserInfo()
	if err != nil {
		utils.GetLogger().Error(err.Error())
		return err
	}

	return nil
}