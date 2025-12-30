package backend

import (
	"context"
	v1 "jh_admin_service/api/backend/upload/v1"
)

type (
	IUpload interface {
		UploadImage(ctx context.Context, req *v1.UploadImageReq) (*v1.UploadImageRes, error)
	}
)

var (
	localUpload IUpload
)

func Upload() IUpload {
	if localUpload == nil {
		panic("implement not found for interface IUpload, forgot register?")
	}
	return localUpload
}

func RegisterUpload(i IUpload) {
	localUpload = i
}
