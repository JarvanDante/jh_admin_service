package role

import (
	"context"
	v1 "jh_user_service/api/role/v1"
	"jh_user_service/internal/service"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedRoleServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterRoleServer(s.Server, &Controller{})
}

func (*Controller) GetRoleList(ctx context.Context, req *v1.GetRoleListReq) (res *v1.GetRoleListRes, err error) {
	return service.Role().GetRoleList(ctx, req)
}
