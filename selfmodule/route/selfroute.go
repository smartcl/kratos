package selfroute

import (
	"github.com/ory/kratos/cipher"
	"github.com/ory/kratos/driver/config"
	"github.com/ory/kratos/hash"
	"github.com/ory/kratos/identity"
	"github.com/ory/kratos/selfmodule/auth_phone"
	"github.com/ory/kratos/x"
	"github.com/ory/x/decoderx"
)

const (
	ImgCode    = "/imgCode"
	PhoneCode  = "/phoneCode"
	AdminLogin = "/login"
	UserList   = "/userList"
)

type (
	handlerDependencies interface {
		identity.PoolProvider
		x.WriterProvider
		config.Provider
		x.CSRFProvider
		cipher.Provider
		hash.HashProvider
	}
	Handler struct {
		r  handlerDependencies
		dx *decoderx.HTTP
	}
)

func NewHandler(r handlerDependencies) *Handler {
	auth_phone.NewAuthPhone()
	return &Handler{
		r:  r,
		dx: decoderx.NewHTTP(),
	}
}

func (h *Handler) RegisterAdminRoutes(admin *x.RouterAdmin) {
	admin.GET(ImgCode, h.ImgCode)
	admin.POST(PhoneCode, h.PhoneCode)
	admin.POST(AdminLogin, h.AdminLogin)
	admin.GET(UserList, h.UserList)
}
