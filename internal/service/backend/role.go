// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ================================================================================

package backend

import (
	"context"
	"jh_app_service/api/backend/role/v1"
)

type (
	IRole interface {
		GetRoleList(ctx context.Context, req *v1.GetRoleListReq) (*v1.GetRoleListRes, error)
		CreateRole(ctx context.Context, req *v1.CreateRoleReq) (*v1.CreateRoleRes, error)
		UpdateRole(ctx context.Context, req *v1.UpdateRoleReq) (*v1.UpdateRoleRes, error)
		DeleteRole(ctx context.Context, req *v1.DeleteRoleReq) (*v1.DeleteRoleRes, error)
		GetPermissions(ctx context.Context, req *v1.GetPermissionsReq) (*v1.GetPermissionsRes, error)
		SavePermission(ctx context.Context, req *v1.SavePermissionReq) (*v1.SavePermissionRes, error)
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
