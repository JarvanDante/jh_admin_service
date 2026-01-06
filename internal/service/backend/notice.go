// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ================================================================================

package backend

import (
	"context"
	"jh_app_service/api/backend/notice/v1"
)

type (
	INotice interface {
		GetNoticeList(ctx context.Context, req *v1.GetNoticeListReq) (*v1.GetNoticeListRes, error)
		CreateNotice(ctx context.Context, req *v1.CreateNoticeReq) (*v1.CreateNoticeRes, error)
		UpdateNotice(ctx context.Context, req *v1.UpdateNoticeReq) (*v1.UpdateNoticeRes, error)
		DeleteNotice(ctx context.Context, req *v1.DeleteNoticeReq) (*v1.DeleteNoticeRes, error)
	}
)

var (
	localNotice INotice
)

func Notice() INotice {
	if localNotice == nil {
		panic("implement not found for interface INotice, forgot register?")
	}
	return localNotice
}

func RegisterNotice(i INotice) {
	localNotice = i
}
