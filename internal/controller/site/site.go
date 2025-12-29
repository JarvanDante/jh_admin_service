package site

import (
	"context"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	v1 "jh_admin_service/api/site/v1"
	"jh_admin_service/internal/service"
)

type Controller struct {
	v1.UnimplementedSiteServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterSiteServer(s.Server, &Controller{})
}

// GetBasicSetting 获取站点基本设置 (gRPC)
func (*Controller) GetBasicSetting(ctx context.Context, req *v1.GetBasicSettingReq) (res *v1.GetBasicSettingRes, err error) {
	return service.Site().GetBasicSetting(ctx, req)
}

// UpdateBasicSetting 更新站点基本设置 (gRPC)
func (*Controller) UpdateBasicSetting(ctx context.Context, req *v1.UpdateBasicSettingReq) (res *v1.UpdateBasicSettingRes, err error) {
	return service.Site().UpdateBasicSetting(ctx, req)
}
