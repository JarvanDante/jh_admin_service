package option

import (
	"context"

	v1 "jh_app_service/api/backend/option/v1"
	"jh_app_service/internal/service/backend"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedOptionServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterOptionServer(s.Server, &Controller{})
}

// GetUserGradeList 获取会员等级列表
func (*Controller) GetUserGradeList(ctx context.Context, req *v1.UserGradeListRequest) (*v1.UserGradeListResponse, error) {
	return backend.Option().GetUserGradeList(ctx, req)
}

// GetAdminRoleList 获取后台角色列表
func (*Controller) GetAdminRoleList(ctx context.Context, req *v1.AdminRoleListRequest) (*v1.AdminRoleListResponse, error) {
	return backend.Option().GetAdminRoleList(ctx, req)
}
