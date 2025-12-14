package middleware

import (
	"testing"

	"github.com/xzwsloser/software_design/backend/dto"
	"github.com/xzwsloser/software_design/backend/utils"
)

func TestJwt(t *testing.T) {
	utils.LoadConfig("../config.json")
	utils.InitLogger()

	j := NewJwt()
	clamis := j.CreatClaims(dto.BasicUserInfo{
		Id: 1,
		Username: "xzw",
	})

	t.Log("-----Create JWT Token-----")
	jwtToken, err := j.CreateJwtToken(clamis)
	if err != nil {
		t.Error(err)
	}

	t.Log("jwt Token: ", jwtToken)

	t.Log("-----Parse JWT Token-----")
	parsedClaims, err := j.ParseToken(jwtToken)
	if err != nil {
		t.Error(err)
	}

	t.Log(parsedClaims.BasicUserInfo)
}