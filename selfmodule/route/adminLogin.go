package selfroute

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/ory/kratos/selfmodule/adminjwt"
	"github.com/ory/kratos/selfmodule/cv"
	"io"
	"net/http"
)

type AdminLoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func (h *Handler) AdminLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &AdminLoginRequest{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
	err = json.Unmarshal(body, req)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
	if req.UserName != cv.AdminUserName || req.Password != cv.AdminPassword {
		h.r.Writer().WriteError(w, r, errors.New("用户名或密码错误"))
		return
	}
	jwt := adminjwt.NewJsonWebToken(&adminjwt.Token{
		UserName: cv.AdminUserName,
	}, "")
	err = jwt.GenerateJwtToken()
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
	cookie := &http.Cookie{
		Name:     cv.AdminJwtCookie,
		Value:    jwt.TokenString,
		Path:     "/",
		HttpOnly: false,
	}
	http.SetCookie(w, cookie)
	h.r.Writer().Write(w, r, "登录成功")
	return
}
