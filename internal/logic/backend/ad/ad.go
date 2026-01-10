package ad

import (
	"context"
	"fmt"
	"jh_app_service/api/backend/ad/v1"
	"jh_app_service/internal/dao"
	"jh_app_service/internal/middleware"
	"jh_app_service/internal/model/do"
	"jh_app_service/internal/model/entity"
	"jh_app_service/internal/service/backend"
	"jh_app_service/internal/util"

	"github.com/gogf/gf/v2/os/gtime"
)

type (
	sAd struct{}
)

func init() {
	backend.RegisterAd(&sAd{})
}

// GetAdList 获取广告列表
func (s *sAd) GetAdList(ctx context.Context, req *v1.GetAdListReq) (*v1.GetAdListRes, error) {
	middleware.LogWithTrace(ctx, "info", "获取广告列表请求 - Page: %d, Size: %d", req.Page, req.Size)

	// 默认站点ID为1
	siteId := int32(1)

	page := req.Page
	if page <= 0 {
		page = 1
	}
	size := req.Size
	if size <= 0 {
		size = 50
	}

	// 构建查询条件
	query := dao.Ad.Ctx(ctx).Where(do.Ad{
		SiteId: siteId,
	})
	if req.Status >= 0 {
		query = query.Where("status", req.Status)
	}
	if req.Position > 0 {
		query = query.Where("position", req.Position)
	}

	// 获取总数
	total, err := query.Count()
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "获取广告总数失败: %v", err)
		return nil, err
	}

	// 分页查询
	var ads []*entity.Ad
	offset := int((page - 1) * size)
	err = query.Fields("id, name, image, url, position, status, sort, created_at, updated_at, start_time, expired_time").
		Order("sort ASC, created_at DESC").
		Limit(offset, int(size)).
		Scan(&ads)
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "获取广告列表失败: %v", err)
		return nil, err
	}

	// 位置映射（写死数据）
	positionMap := map[int]string{
		1: "首页轮播",
		2: "首页弹窗",
		3: "侧边栏",
		4: "底部横幅",
		5: "浮动广告",
	}

	// 状态映射
	statusMap := map[int]string{
		0: "禁用",
		1: "启用",
	}

	// 转换为响应格式
	adList := make([]*v1.AdItem, 0, len(ads))
	for _, ad := range ads {
		// 使用修复后的工具函数处理时间格式化
		startTime := util.FormatTime(ad.StartTime)
		expiredTime := util.FormatTime(ad.ExpiredTime)
		createdAt := util.FormatTime(ad.CreatedAt)
		updatedAt := util.FormatTime(ad.UpdatedAt)

		adList = append(adList, &v1.AdItem{
			Id:           int32(ad.Id),
			Title:        ad.Name,
			Image:        ad.Image,
			Link:         ad.Url,
			Position:     int32(ad.Position),
			PositionName: positionMap[ad.Position],
			Status:       int32(ad.Status),
			StatusName:   statusMap[ad.Status],
			Sort:         int32(ad.Sort),
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
			StartTime:    startTime,
			ExpiredTime:  expiredTime,
		})
	}

	// 位置列表
	positionList := []*v1.PositionItem{
		{Value: 1, Name: "首页轮播"},
		{Value: 2, Name: "首页弹窗"},
		{Value: 3, Name: "侧边栏"},
		{Value: 4, Name: "底部横幅"},
		{Value: 5, Name: "浮动广告"},
	}

	// 状态列表
	statusList := []*v1.StatusItem{
		{Value: 0, Name: "禁用"},
		{Value: 1, Name: "启用"},
	}

	middleware.LogWithTrace(ctx, "info", "获取广告列表成功 - 总数: %d, 当前页: %d", total, page)

	return &v1.GetAdListRes{
		List:         adList,
		Total:        int32(total),
		CurrentPage:  page,
		PerPage:      size,
		PositionList: positionList,
		StatusList:   statusList,
	}, nil
}

// CreateAd 创建广告
func (s *sAd) CreateAd(ctx context.Context, req *v1.CreateAdReq) (*v1.CreateAdRes, error) {
	middleware.LogWithTrace(ctx, "info", "创建广告请求 - Title: %s", req.Title)

	// 默认站点ID为1
	siteId := int32(1)

	// 解析开始时间和过期时间
	startTime, err := gtime.StrToTime(req.StartTime)
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "开始时间格式错误: %v", err)
		return nil, fmt.Errorf("开始时间格式错误")
	}

	expiredTime, err := gtime.StrToTime(req.ExpiredTime)
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "过期时间格式错误: %v", err)
		return nil, fmt.Errorf("过期时间格式错误")
	}

	// 创建广告记录
	ad := &do.Ad{
		SiteId:      siteId,
		Name:        req.Title,
		Image:       req.Image,
		Url:         req.Link,
		Position:    int(req.Position),
		Status:      int(req.Status),
		Sort:        int(req.Sort),
		StartTime:   startTime,
		ExpiredTime: expiredTime,
		CreatedAt:   gtime.Now(),
		UpdatedAt:   gtime.Now(),
	}

	_, err = dao.Ad.Ctx(ctx).Data(ad).Insert()
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "创建广告失败: %v", err)
		return nil, err
	}

	middleware.LogWithTrace(ctx, "info", "创建广告成功 - Title: %s", req.Title)

	return &v1.CreateAdRes{}, nil
}

// UpdateAd 更新广告
func (s *sAd) UpdateAd(ctx context.Context, req *v1.UpdateAdReq) (*v1.UpdateAdRes, error) {
	middleware.LogWithTrace(ctx, "info", "更新广告请求 - Id: %d, Title: %s", req.Id, req.Title)

	// 默认站点ID为1
	siteId := int32(1)

	// 检查广告是否存在
	count, err := dao.Ad.Ctx(ctx).Where("id", req.Id).Where("site_id", siteId).Count()
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "查询广告失败: %v", err)
		return nil, err
	}
	if count == 0 {
		middleware.LogWithTrace(ctx, "error", "广告不存在 - Id: %d", req.Id)
		return nil, fmt.Errorf("广告不存在")
	}

	// 解析开始时间和过期时间
	startTime, err := gtime.StrToTime(req.StartTime)
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "开始时间格式错误: %v", err)
		return nil, fmt.Errorf("开始时间格式错误")
	}

	expiredTime, err := gtime.StrToTime(req.ExpiredTime)
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "过期时间格式错误: %v", err)
		return nil, fmt.Errorf("过期时间格式错误")
	}

	// 更新广告
	updateData := &do.Ad{
		Name:        req.Title,
		Image:       req.Image,
		Url:         req.Link,
		Position:    int(req.Position),
		Status:      int(req.Status),
		Sort:        int(req.Sort),
		StartTime:   startTime,
		ExpiredTime: expiredTime,
		UpdatedAt:   gtime.Now(),
	}

	_, err = dao.Ad.Ctx(ctx).Where("id", req.Id).Where("site_id", siteId).Data(updateData).Update()
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "更新广告失败: %v", err)
		return nil, err
	}

	middleware.LogWithTrace(ctx, "info", "更新广告成功 - Id: %d, Title: %s", req.Id, req.Title)

	return &v1.UpdateAdRes{}, nil
}

// DeleteAd 删除广告
func (s *sAd) DeleteAd(ctx context.Context, req *v1.DeleteAdReq) (*v1.DeleteAdRes, error) {
	middleware.LogWithTrace(ctx, "info", "删除广告请求 - Id: %d", req.Id)

	// 默认站点ID为1
	siteId := int32(1)

	// 检查广告是否存在
	count, err := dao.Ad.Ctx(ctx).Where("id", req.Id).Where("site_id", siteId).Count()
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "查询广告失败: %v", err)
		return nil, err
	}
	if count == 0 {
		middleware.LogWithTrace(ctx, "error", "广告不存在 - Id: %d", req.Id)
		return nil, fmt.Errorf("广告不存在")
	}

	// 删除广告
	_, err = dao.Ad.Ctx(ctx).Where("id", req.Id).Where("site_id", siteId).Delete()
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "删除广告失败: %v", err)
		return nil, err
	}

	middleware.LogWithTrace(ctx, "info", "删除广告成功 - Id: %d", req.Id)

	return &v1.DeleteAdRes{}, nil
}
