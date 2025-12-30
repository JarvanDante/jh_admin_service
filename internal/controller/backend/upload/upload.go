package upload

import (
	"context"
	v1 "jh_admin_service/api/backend/upload/v1"
	"jh_admin_service/internal/service/backend"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedUploadServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterUploadServer(s.Server, &Controller{})
}

// UploadImage 上传图片
func (*Controller) UploadImage(ctx context.Context, req *v1.UploadImageReq) (res *v1.UploadImageRes, err error) {
	return backend.Upload().UploadImage(ctx, req)
}
