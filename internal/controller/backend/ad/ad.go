package ad

import (
	"context"
	v1 "jh_app_service/api/backend/ad/v1"
	"jh_app_service/internal/service/backend"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedAdServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterAdServer(s.Server, &Controller{})
}

// GetAdList 获取广告列表
func (*Controller) GetAdList(ctx context.Context, req *v1.GetAdListReq) (res *v1.GetAdListRes, err error) {
	return backend.Ad().GetAdList(ctx, req)
}

// CreateAd 创建广告
func (*Controller) CreateAd(ctx context.Context, req *v1.CreateAdReq) (res *v1.CreateAdRes, err error) {
	return backend.Ad().CreateAd(ctx, req)
}

// UpdateAd 更新广告
func (*Controller) UpdateAd(ctx context.Context, req *v1.UpdateAdReq) (res *v1.UpdateAdRes, err error) {
	return backend.Ad().UpdateAd(ctx, req)
}

// DeleteAd 删除广告
func (*Controller) DeleteAd(ctx context.Context, req *v1.DeleteAdReq) (res *v1.DeleteAdRes, err error) {
	return backend.Ad().DeleteAd(ctx, req)
}
