package notice

import (
	"context"
	v1 "jh_app_service/api/backend/notice/v1"
	"jh_app_service/internal/service/backend"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedNoticeServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterNoticeServer(s.Server, &Controller{})
}

// GetNoticeList 获取公告列表
func (*Controller) GetNoticeList(ctx context.Context, req *v1.GetNoticeListReq) (res *v1.GetNoticeListRes, err error) {
	return backend.Notice().GetNoticeList(ctx, req)
}

// CreateNotice 创建公告
func (*Controller) CreateNotice(ctx context.Context, req *v1.CreateNoticeReq) (res *v1.CreateNoticeRes, err error) {
	return backend.Notice().CreateNotice(ctx, req)
}

// UpdateNotice 更新公告
func (*Controller) UpdateNotice(ctx context.Context, req *v1.UpdateNoticeReq) (res *v1.UpdateNoticeRes, err error) {
	return backend.Notice().UpdateNotice(ctx, req)
}

// DeleteNotice 删除公告
func (*Controller) DeleteNotice(ctx context.Context, req *v1.DeleteNoticeReq) (res *v1.DeleteNoticeRes, err error) {
	return backend.Notice().DeleteNotice(ctx, req)
}
