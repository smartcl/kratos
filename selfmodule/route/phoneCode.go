package selfroute

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/ory/kratos/selfmodule/auth_img"
	"github.com/ory/kratos/selfmodule/auth_phone"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

type PhoneCodeRequest struct {
	Phone   string `json:"phone"`
	ImgCode string `json:"img_code"`
	ImgId   string `json:"id"`
}

type PhoneCodeResponse struct {
	Id  string `json:"id"`
	Img string `json:"code"`
}

func (h *Handler) PhoneCode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &PhoneCodeRequest{}
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
	if req.Phone == "" {
		h.r.Writer().WriteError(w, r, errors.New("手机号不能为空"))
		return
	}
	if req.ImgCode == "" {
		h.r.Writer().WriteError(w, r, errors.New("验证码不能为空"))
		return
	}
	if req.ImgId == "" {
		h.r.Writer().WriteError(w, r, errors.New("验证码id不能为空"))
		return
	}
	// 检查图片验证码
	verifyRes := auth_img.VerifyCaptcha(req.ImgId, req.ImgCode)
	if !verifyRes {
		h.r.Writer().WriteError(w, r, errors.New("验证码错误"))
		return
	}
	_, err = auth_phone.AuthPhoneGlobal.GenerateAuthCode(req.Phone)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
	h.r.Writer().Write(w, r, "发送成功")
	return
}
