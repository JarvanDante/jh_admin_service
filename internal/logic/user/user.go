package user

import (
	"jh_user_service/internal/service"
)

type (
	sUser struct{}
)

func init() {
	service.RegisterUser(&sUser{})
}

// 用户相关的服务方法已删除
