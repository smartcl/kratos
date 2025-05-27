package adminjwt

import (
	"fmt"
	jwtV4 "github.com/golang-jwt/jwt/v4"
	"github.com/ory/kratos/selfmodule/cv"
	"time"
)

// Token 授权码
type Token struct {
	UserName       string
	StandardClaims jwtV4.RegisteredClaims
}

// JsonWebToken token对象
type JsonWebToken struct {
	TokenObj    *Token
	TokenString string // token 本身
}

func NewJsonWebToken(token *Token, tokenString string) *JsonWebToken {
	return &JsonWebToken{
		TokenObj:    token,
		TokenString: tokenString,
	}
}

// Valid  必须实现的接口
//
//	@receiver t
//	@return error
func (t *Token) Valid() error {
	return nil
}

// GenerateJwtToken  生成token
//
//	@receiver t
//	@param key
//	@param issuer
//	@param expiresTime
//	@return error
func (t *JsonWebToken) GenerateJwtToken() error {
	var err error
	hmacSampleSecret := []byte(cv.AdminJwtKey) // 密钥，不能泄露
	token := jwtV4.New(jwtV4.SigningMethodHS256)
	nowTime := time.Now()
	expirseAt := nowTime.Add(time.Second * time.Duration(cv.AdminJwtExpTime))
	t.TokenObj.StandardClaims = jwtV4.RegisteredClaims{
		NotBefore: jwtV4.NewNumericDate(nowTime),   // 签名生效时间
		ExpiresAt: jwtV4.NewNumericDate(expirseAt), // 签名过期时间
		Issuer:    cv.AdminJwtIssue,                // 签名颁发者
	}
	token.Claims = t.TokenObj
	t.TokenString, err = token.SignedString(hmacSampleSecret)
	return err
}

// ParseJwtToken  解token
//
//	@receiver t
//	@param key
//	@return error
func (t *JsonWebToken) ParseJwtToken() error {
	var hmacSampleSecret = []byte(cv.AdminJwtKey)
	// 前面例子生成的token
	token, err := jwtV4.ParseWithClaims(t.TokenString, &Token{}, func(t *jwtV4.Token) (interface{}, error) {
		return hmacSampleSecret, nil
	})

	if err != nil {
		return err
	}
	tmp, _ := token.Claims.(*Token)
	t.TokenObj = tmp
	return nil
}

func (t *JsonWebToken) VerifyTime() bool {
	return t.TokenObj.StandardClaims.ExpiresAt.After(time.Now())
}

func VerifyAdminToken(tokenString string) error {
	t := &JsonWebToken{
		TokenString: tokenString,
	}
	err := t.ParseJwtToken()
	if err != nil {
		fmt.Println("err: ", err)
		return fmt.Errorf("请重新登录")
	}
	res := t.VerifyTime()
	if !res {
		fmt.Println("登录过期")
		return fmt.Errorf("请重新登录")
	}
	if t.TokenObj.UserName != cv.AdminUserName {
		fmt.Println("用户名错误")
		return fmt.Errorf("无此权限")
	}
	return nil
}
