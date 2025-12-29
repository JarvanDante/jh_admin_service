package role

import (
	"context"
	v1 "jh_user_service/api/role/v1"
	"jh_user_service/internal/middleware"
	"jh_user_service/internal/service"
)

type (
	sRole struct{}
)

func init() {
	service.RegisterRole(&sRole{})
}

func (s *sRole) GetRoleList(ctx context.Context, req *v1.GetRoleListReq) (*v1.GetRoleListRes, error) {
	middleware.LogWithTrace(ctx, "info", "角色列表请求 - SiteId: %d", req.SiteId)

	// 默认站点ID为1，如果请求中有指定则使用指定的
	siteId := int32(1)
	if req.SiteId > 0 {
		siteId = req.SiteId
	}

	middleware.LogWithTrace(ctx, "info", "获取角色列表请求成功 - SiteId: %d", siteId)
	return nil, nil
}
