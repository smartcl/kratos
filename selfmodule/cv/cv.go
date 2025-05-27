package cv

import (
	"fmt"
	"os"
	"strconv"
)

var (
	AdminUserName = "admin"
	AdminPassword = "123456"

	AdminJwtKey = "qaplwo..0_0"

	AdminJwtExpTime = 3600
)

const (
	AdminJwtCookie = "admin_jwt"
	AdminJwtIssue  = "Chainmaker"

	UserName   = "user_name"
	Phone      = "phone"
	AuthStatus = "auth_status"
	Type       = "type"
)

func init() {
	var err error
	if adminUserName, exists := os.LookupEnv("ADMIN_USERNAME"); exists {
		AdminUserName = adminUserName
	}
	if adminPassword, exists := os.LookupEnv("ADMIN_PASSWORD"); exists {
		AdminPassword = adminPassword
	}
	if adminJwtKey, exists := os.LookupEnv("ADMIN_JWT_KEY"); exists {
		AdminJwtKey = adminJwtKey
	}
	if adminJwtExpTime, exists := os.LookupEnv("ADMIN_JWT_EXP_TIME"); exists {
		AdminJwtExpTime, err = strconv.Atoi(adminJwtExpTime)
		if err != nil {
			panic(fmt.Errorf("invalid admin jwt exp time: %s", err))
		}
	}
}
