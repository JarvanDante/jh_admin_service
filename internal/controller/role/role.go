package role

import (
	"context"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"

	v1 "jh_admin_service/api/role/v1"
	"jh_admin_service/internal/service"
)

type Controller struct {
	v1.UnimplementedRoleServer
}

// Register 注册角色服务
func Register(s *grpcx.GrpcServer) {
	v1.RegisterRoleServer(s.Server, &Controller{})
}

// GetRoleList 获取角色列表
func (*Controller) GetRoleList(ctx context.Context, req *v1.GetRoleListReq) (*v1.GetRoleListRes, error) {
	return service.Role().GetRoleList(ctx, req)
}
