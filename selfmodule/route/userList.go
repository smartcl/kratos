package selfroute

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/ory/kratos/selfmodule/cv"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	Total int             `json:"total"`
	Data  []*UserListData `json:"data"`
}

type UserListData struct {
	ID         string  `json:"id"`
	Username   string  `json:"user_name"`
	Type       float64 `json:"type"`
	Phone      string  `json:"phone"`
	CreateAt   string  `json:"createAt"`
	AuthStatus float64 `json:"auth_status"`
}

func (h *Handler) UserList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		pageSize   = 10
		page       = 1
		authStatus int
		err        error
	)
	pageSizeStr := r.URL.Query().Get("page_size")
	pageStr := r.URL.Query().Get("page")
	userName := r.URL.Query().Get("user_name")
	authStatusStr := r.URL.Query().Get("auth_status")

	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			h.r.Writer().WriteError(w, r, fmt.Errorf("page_size must be int"))
			return
		}
	}
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			h.r.Writer().WriteError(w, r, fmt.Errorf("page must be int"))
			return
		}
	}
	if authStatusStr != "" {
		authStatus, err = strconv.Atoi(authStatusStr)
		if err != nil {
			h.r.Writer().WriteError(w, r, fmt.Errorf("auth_status must be int"))
			return
		}
	}

	ctx := r.Context()
	list, count, err := h.r.IdentityPool().ListIdentitiesByUserNameOrAuth(ctx, pageSize, page, authStatus, userName)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
	data := make([]*UserListData, 0)
	for _, v := range list {
		traitsByteStr, err := v.Traits.MarshalJSON()
		if err != nil {
			h.r.Writer().WriteError(w, r, err)
			return
		}
		traitsMap := make(map[string]interface{})
		err = json.Unmarshal(traitsByteStr, &traitsMap)
		if err != nil {
			h.r.Writer().WriteError(w, r, err)
			return
		}
		userNameRes := traitsMap[cv.UserName].(string)
		phone := traitsMap[cv.Phone].(string)

		publicMetadataByteStr, err := v.MetadataPublic.MarshalJSON()
		if err != nil {
			h.r.Writer().WriteError(w, r, err)
			return
		}
		publicMetadataMap := make(map[string]interface{})
		err = json.Unmarshal(publicMetadataByteStr, &publicMetadataMap)
		if err != nil {
			h.r.Writer().WriteError(w, r, err)
			return
		}
		authStatusRes := publicMetadataMap[cv.AuthStatus].(float64)
		typeRes := publicMetadataMap[cv.Type].(float64)

		data = append(data, &UserListData{
			ID:         v.ID.String(),
			Username:   userNameRes,
			Type:       typeRes,
			Phone:      phone,
			CreateAt:   v.CreatedAt.String(),
			AuthStatus: authStatusRes,
		})
	}
	h.r.Writer().Write(w, r, &UserListResponse{
		Total: count,
		Data:  data,
	})
}
