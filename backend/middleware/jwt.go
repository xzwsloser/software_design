package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xzwsloser/software_design/backend/dto"
	"github.com/xzwsloser/software_design/backend/utils"
)

const (
	JWT_HEADER_KEY = "Authorization"
)

type CustomClamis struct {
	dto.BasicUserInfo
	// jwt 官方字段
	jwt.RegisteredClaims
}

type JWT struct {
	// jwt 签名
	SigningKey []byte
}

var (
	TokenExpired 		= errors.New("Token is expired")
	TokenNotValidYet	= errors.New("Token is active yet") 
	TokenMalFormed 		= errors.New("Not a token")
	TokenInValid 		= errors.New("Cannnot Handle the token")
)

func NewJwt() *JWT {
	return &JWT{
		[]byte(utils.GetJwtConfig().SerectKey),
	}
}

// @Description: 利用用户信息生成 jwt payload
func (j *JWT) CreatClaims(userInfo dto.BasicUserInfo) CustomClamis {
	now := time.Now().Unix()

	claims := CustomClamis{
		BasicUserInfo: userInfo,

		// jwt payload 官方字段信息
		RegisteredClaims: jwt.RegisteredClaims{
			// jwt token 生效事件
			NotBefore: jwt.NewNumericDate(time.Unix(now-1000,0)),
			// jwt token 过期事件(7天)
			ExpiresAt: jwt.NewNumericDate(time.Unix(now+604800,0)),
			// jwt 签发者
			Issuer: utils.GetJwtConfig().Issuer,
			// jwt 签发时间
			IssuedAt: jwt.NewNumericDate(time.Unix(now, 0)),
		},
	}

	return claims
}

// @Description: 生成 jwt token
func (j *JWT) CreateJwtToken(claims CustomClamis) (string, error) {
	// 使用 HS256 Hash 算法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// @Description: 从 jwt token 中解析出负载信息
func (j *JWT) ParseToken(jwtToken string) (*CustomClamis, error) {
	// token 本体对象
	// 需要 jwt 密钥
	token, err := jwt.ParseWithClaims(jwtToken, &CustomClamis{}, func(token *jwt.Token) (any , error) {
		return j.SigningKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, TokenMalFormed
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, TokenExpired
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, TokenNotValidYet
		} else {
			return nil, TokenInValid
		}
	}

	if token.Valid {
		if clamis, ok := token.Claims.(*CustomClamis) ; ok {
			return clamis, nil
		}
		return nil, TokenInValid
	}

	return nil, TokenInValid
}

// @Descipriont: JWT 鉴权中间件
func JwtAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jwtToken := ctx.Request.Header.Get(JWT_HEADER_KEY)
		if jwtToken == "" {
			utils.GetLogger().Error("JWT Token is Empty")
			ctx.JSON(http.StatusUnauthorized, dto.Fail("JWT Token is Empty"))
			ctx.Abort()
			return 
		}

		j := NewJwt()
		claims, err := j.ParseToken(jwtToken)
		if err != nil {
			if errors.Is(err, TokenExpired) {
				utils.GetLogger().Warn("Token is Expired")
				ctx.JSON(http.StatusOK, dto.Fail("Token is Expired"))
				ctx.Abort()
				return
			}
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}

// @Description: 直接从 jwt token 中解析出 clamis
func GetClamisWithCtx(c *gin.Context) (*CustomClamis, error) {
	jwtToken := c.Request.Header.Get(JWT_HEADER_KEY)
	j := NewJwt()
	clamis, err := j.ParseToken(jwtToken)
	if err != nil {
		utils.GetLogger().Error("Failed to get parsed clamis")
	}

	return clamis, err
}

// @Description: 获取到用户信息(JWT Token or 用户信息)
func GetBasicUserInfo(c *gin.Context) (dto.BasicUserInfo, error) {
	if clamis, exists := c.Get("claims") ; !exists {
		if cl, err := GetClamisWithCtx(c) ; err != nil {
			return dto.BasicUserInfo{}, nil
		} else {
			return cl.BasicUserInfo, nil
		}
	} else {
		cl, _ := clamis.(*CustomClamis)
		return cl.BasicUserInfo, nil
	}
}

