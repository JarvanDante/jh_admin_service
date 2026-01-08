// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ================================================================================

package backend

import (
	"context"
	"jh_app_service/api/backend/option/v1"
)

type (
	IOption interface {
		GetUserGradeList(ctx context.Context, req *v1.UserGradeListRequest) (*v1.UserGradeListResponse, error)
		GetAdminRoleList(ctx context.Context, req *v1.AdminRoleListRequest) (*v1.AdminRoleListResponse, error)
	}
)

var (
	localOption IOption
)

func Option() IOption {
	if localOption == nil {
		panic("implement not found for interface IOption, forgot register?")
	}
	return localOption
}

func RegisterOption(i IOption) {
	localOption = i
}
