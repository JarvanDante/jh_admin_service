package role

import (
	"context"
	"fmt"

	v1 "jh_admin_service/api/role/v1"
	"jh_admin_service/internal/dao"
	"jh_admin_service/internal/middleware"
	"jh_admin_service/internal/model/do"
	"jh_admin_service/internal/model/entity"
	"jh_admin_service/internal/service"
)

type (
	sRole struct{}
)

func init() {
	service.RegisterRole(&sRole{})
}

// GetRoleList 获取角色列表
func (s *sRole) GetRoleList(ctx context.Context, req *v1.GetRoleListReq) (*v1.GetRoleListRes, error) {
	middleware.LogWithTrace(ctx, "info", "获取角色列表请求 - SiteId: %d", req.SiteId)

	// 默认站点ID为1，如果请求中有指定则使用指定的
	siteId := int32(1)
	if req.SiteId > 0 {
		siteId = req.SiteId
	}

	// 查询角色列表
	var roles []*entity.AdminRole
	err := dao.AdminRole.Ctx(ctx).Where(do.AdminRole{
		SiteId: siteId,
		Status: 1, // 只查询启用的角色
	}).OrderAsc(dao.AdminRole.Columns().Id).Scan(&roles)

	if err != nil {
		middleware.LogWithTrace(ctx, "error", "查询角色列表失败: %v", err)
		return nil, fmt.Errorf("查询角色列表失败: %v", err)
	}

	// 构建响应数据
	var roleList []*v1.RoleInfo
	for _, role := range roles {
		roleInfo := &v1.RoleInfo{
			Id:          int32(role.Id),
			SiteId:      int32(role.SiteId),
			Name:        role.Name,
			Status:      int32(role.Status),
			Permissions: role.Permissions,
		}

		// 处理时间字段
		if role.CreatedAt != nil {
			roleInfo.CreatedAt = role.CreatedAt.Unix()
		}
		if role.UpdatedAt != nil {
			roleInfo.UpdatedAt = role.UpdatedAt.Unix()
		}

		roleList = append(roleList, roleInfo)
	}

	middleware.LogWithTrace(ctx, "info", "获取角色列表成功 - SiteId: %d, 角色数量: %d", siteId, len(roleList))

	return &v1.GetRoleListRes{
		Roles: roleList,
	}, nil
}
