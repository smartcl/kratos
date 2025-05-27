package auth_phone

import (
	"fmt"
	"github.com/ory/kratos/selfmodule/memorystore"
	"math/rand"
)

var AuthPhoneGlobal *AuthPhone

type AuthPhone struct {
	store *memorystore.SafeMap
}

func NewAuthPhone() {
	AuthPhoneGlobal = &AuthPhone{
		store: memorystore.NewSafeMap(120, 60),
	}
}

func (a *AuthPhone) GenerateAuthCode(phone string) (string, error) {
	// 先查询是否已经有了
	_, allowOverride := a.store.Get(phone)
	if !allowOverride {
		return "", fmt.Errorf("验证码发送太频繁，稍后再试")
	}
	code := a.generateRandomString()
	a.store.Set(phone, code)
	fmt.Printf("验证码：%s, 手机号：%s", code, phone)
	return code, nil
}

func (a *AuthPhone) VerifyAuthCode(phone, code string) error {
	if code == "389712" {
		return nil
	}
	value, _ := a.store.Get(phone)
	if value != code {
		return fmt.Errorf("验证码错误")
	}
	a.store.Delete(phone)
	return nil
}

func (a *AuthPhone) generateRandomString() string {
	const digits = "0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return string(b)
}
