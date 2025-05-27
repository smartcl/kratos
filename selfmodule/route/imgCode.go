package selfroute

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ory/kratos/selfmodule/auth_img"
	"net/http"
)

type ImgCodeRequest struct {
	CodeType int `json:"code_type"`
}

type ImgCodeResponse struct {
	Id  string `json:"id"`
	Img string `json:"img"`
}

func (h *Handler) ImgCode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &ImgCodeRequest{
		CodeType: auth_img.AuthImgCaptchaTypeString,
	}
	id, b64s, err := auth_img.CreateCode(req.CodeType)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
	h.r.Writer().Write(w, r, &ImgCodeResponse{
		Id:  id,
		Img: b64s,
	})
	return
}
