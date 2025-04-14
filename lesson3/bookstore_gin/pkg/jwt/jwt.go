package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// CustomClaims 自定义声明类型 并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段
// 假设我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type CustomClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

const TokenExpireDuration = time.Minute * 10
const RefreshTokenExpireDuration = time.Hour * 24 * 7

var CustomSecret = []byte("君子自强不息")
var InvalidAuth = errors.New("invalid token")

// GenToken 生成JWT
func GenToken(userID int64, username string) (string, error) {
	customClaims := &CustomClaims{
		userID,
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)), // 过期时间
			Issuer:    "youngking",                                             // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(CustomSecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
	cm := &CustomClaims{}
	// jwt.ParseWithClaims会解析tokenString后创建一个token结构体，并填入数据
	// 回调keyFunc拿到CustomSerec
	// jwt.ParseWithClaims在由tokenString中的Header和payload+CustomSerec生成一个newtokenString
	// 将newtokenString中的签名和tokenString中的签名比较
	token, err := jwt.ParseWithClaims(tokenString, cm, func(token *jwt.Token) (interface{}, error) {
		return CustomSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, InvalidAuth
	}
	return cm, nil
}

func GenAAndRToken(userID int64, username string) (accessToken, refreshToken string, err error) {
	customClaims := &CustomClaims{
		userID,
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)), // 过期时间
			Issuer:    "youngking",                                             // 签发人
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims).SignedString(CustomSecret)
	if err != nil {
		return "", "", err
	}
	refreshClaims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(RefreshTokenExpireDuration).Unix(),
		Issuer:    "youngking",
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(CustomSecret)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func RefreshToken(aToken, rToken string) (newToken, newRToken string, err error) {
	// 解析 Refresh Token
	_, err = jwt.Parse(rToken, func(token *jwt.Token) (interface{}, error) {
		return CustomSecret, nil
	})
	if err != nil { // Refresh Token 无效，返回错误
		return "", "", InvalidAuth
	}

	// 解析 Access Token
	claims := new(CustomClaims)
	_, err = jwt.ParseWithClaims(aToken, claims, func(token *jwt.Token) (interface{}, error) {
		return CustomSecret, nil
	})

	// 判断 Access Token 是否过期
	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			// Access Token 过期，重新生成新 Token
			return GenAAndRToken(claims.UserID, claims.Username)
		}
	}

	// Access Token 未过期，直接返回原 Token
	return aToken, rToken, nil
}
