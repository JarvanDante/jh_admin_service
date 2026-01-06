// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ================================================================================

package backend

import (
	"context"
	"jh_app_service/api/backend/ad/v1"
)

type (
	IAd interface {
		GetAdList(ctx context.Context, req *v1.GetAdListReq) (*v1.GetAdListRes, error)
		CreateAd(ctx context.Context, req *v1.CreateAdReq) (*v1.CreateAdRes, error)
		UpdateAd(ctx context.Context, req *v1.UpdateAdReq) (*v1.UpdateAdRes, error)
		DeleteAd(ctx context.Context, req *v1.DeleteAdReq) (*v1.DeleteAdRes, error)
	}
)

var (
	localAd IAd
)

func Ad() IAd {
	if localAd == nil {
		panic("implement not found for interface IAd, forgot register?")
	}
	return localAd
}

func RegisterAd(i IAd) {
	localAd = i
}
