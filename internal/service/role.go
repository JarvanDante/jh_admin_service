// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ================================================================================

package service

import (
	"context"

	v1 "jh_admin_service/api/role/v1"
)

type (
	IRole interface {
		GetRoleList(ctx context.Context, req *v1.GetRoleListReq) (*v1.GetRoleListRes, error)
	}
)

var (
	localRole IRole
)

func Role() IRole {
	if localRole == nil {
		panic("implement not found for interface IRole, forgot register?")
	}
	return localRole
}

func RegisterRole(i IRole) {
	localRole = i
}
