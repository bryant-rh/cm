package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//token的过期时长
const TokenExpireDuration = time.Hour * 2

//const TokenExpireDuration = time.Minute * 1

//secret,签名时使用
var MySecret = []byte("bryant-rh")

//用来生成token的struct
type MyClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

//创建token
func GenToken(username, password string) (string, error) {
	c := MyClaims{
		username, // 自定义字段
		password, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "bryant-rh",                                // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// 解析token
func ParseToken(tokenString string) (*MyClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	// if err != nil {
	// 	return nil, err
	// }
	// if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
	// 	// fmt.Println("jwt ok")
	// 	// fmt.Println(claims.Username)
	// 	return claims, nil
	// }
	if token.Valid { //服务端验证token是否有效
		return token.Claims.(*MyClaims), nil

	} else if ve, ok := err.(*jwt.ValidationError); ok { //
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, errors.New("invalid token")
		} else if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
			return nil, errors.New("token is expired")
		} else if ve.Errors&(jwt.ValidationErrorNotValidYet) != 0 {
			return nil, errors.New("token is not valid yet")
		}

	}
	return nil, err
	//return nil, errors.New("invalid token")
}

// // 更新token
// func RefreshToken(tokenString string) (string, error) {
// 	jwt.TimeFunc = func() time.Time {
// 		return time.Unix(0, 0)
// 	}
// 	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		return j.SigningKey, nil
// 	})
// 	if err != nil {
// 		return "", err
// 	}
// 	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
// 		jwt.TimeFunc = time.Now
// 		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
// 		return j.CreateToken(*claims)
// 	}
// 	return "", TokenInvalid
// }
