package notice

import (
	"context"
	"fmt"
	"jh_app_service/api/backend/notice/v1"
	"jh_app_service/internal/dao"
	"jh_app_service/internal/middleware"
	"jh_app_service/internal/model/do"
	"jh_app_service/internal/model/entity"
	"jh_app_service/internal/service/backend"

	"github.com/gogf/gf/v2/os/gtime"
)

type (
	sNotice struct{}
)

func init() {
	backend.RegisterNotice(&sNotice{})
}

// GetNoticeList 获取公告列表
func (s *sNotice) GetNoticeList(ctx context.Context, req *v1.GetNoticeListReq) (*v1.GetNoticeListRes, error) {
	middleware.LogWithTrace(ctx, "info", "获取公告列表请求 - Page: %d, Size: %d", req.Page, req.Size)

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
	query := dao.Notice.Ctx(ctx).Where(do.Notice{
		SiteId: siteId,
	})
	if req.Status >= 0 {
		query = query.Where("status", req.Status)
	}
	if req.Type > 0 {
		query = query.Where("type", req.Type)
	}

	// 获取总数
	total, err := query.Count()
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "获取公告总数失败: %v", err)
		return nil, err
	}

	// 分页查询
	var notices []*entity.Notice
	offset := int((page - 1) * size)
	err = query.Fields("id, title, content, type, status, sort, created_at, updated_at").
		Order("sort ASC, created_at DESC").
		Limit(offset, int(size)).
		Scan(&notices)
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "获取公告列表失败: %v", err)
		return nil, err
	}

	// 类型映射（写死数据）
	typeMap := map[int]string{
		1: "系统公告",
		2: "活动公告",
		3: "维护公告",
	}

	// 状态映射
	statusMap := map[int]string{
		0: "禁用",
		1: "启用",
	}

	// 转换为响应格式
	noticeList := make([]*v1.NoticeItem, 0, len(notices))
	for _, notice := range notices {
		// 默认置顶为0，如果实体有该字段则使用
		isTop := int32(0)

		noticeList = append(noticeList, &v1.NoticeItem{
			Id:         int32(notice.Id),
			Title:      notice.Title,
			Content:    notice.Content,
			Type:       int32(notice.Type),
			TypeName:   typeMap[notice.Type],
			Status:     int32(notice.Status),
			StatusName: statusMap[notice.Status],
			Sort:       int32(notice.Sort),
			IsTop:      isTop,
			CreatedAt:  notice.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  notice.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// 类型列表
	typeList := []*v1.TypeItem{
		{Value: 1, Name: "系统公告"},
		{Value: 2, Name: "活动公告"},
		{Value: 3, Name: "维护公告"},
	}

	// 状态列表
	statusList := []*v1.StatusItem{
		{Value: 0, Name: "禁用"},
		{Value: 1, Name: "启用"},
	}

	middleware.LogWithTrace(ctx, "info", "获取公告列表成功 - 总数: %d, 当前页: %d", total, page)

	return &v1.GetNoticeListRes{
		List:        noticeList,
		Total:       int32(total),
		CurrentPage: page,
		PerPage:     size,
		TypeList:    typeList,
		StatusList:  statusList,
	}, nil
}

// CreateNotice 创建公告
func (s *sNotice) CreateNotice(ctx context.Context, req *v1.CreateNoticeReq) (*v1.CreateNoticeRes, error) {
	middleware.LogWithTrace(ctx, "info", "创建公告请求 - Title: %s", req.Title)

	// 默认站点ID为1
	siteId := int32(1)

	// 创建公告记录
	notice := &do.Notice{
		SiteId:    siteId,
		Title:     req.Title,
		Content:   req.Content,
		Type:      int(req.Type),
		Status:    int(req.Status),
		Sort:      int(req.Sort),
		CreatedAt: gtime.Now(),
		UpdatedAt: gtime.Now(),
	}

	_, err := dao.Notice.Ctx(ctx).Data(notice).Insert()
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "创建公告失败: %v", err)
		return nil, err
	}

	middleware.LogWithTrace(ctx, "info", "创建公告成功 - Title: %s", req.Title)

	return &v1.CreateNoticeRes{}, nil
}

// UpdateNotice 更新公告
func (s *sNotice) UpdateNotice(ctx context.Context, req *v1.UpdateNoticeReq) (*v1.UpdateNoticeRes, error) {
	middleware.LogWithTrace(ctx, "info", "更新公告请求 - Id: %d, Title: %s", req.Id, req.Title)

	// 默认站点ID为1
	siteId := int32(1)

	// 检查公告是否存在
	count, err := dao.Notice.Ctx(ctx).Where("id", req.Id).Where("site_id", siteId).Count()
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "查询公告失败: %v", err)
		return nil, err
	}
	if count == 0 {
		middleware.LogWithTrace(ctx, "error", "公告不存在 - Id: %d", req.Id)
		return nil, fmt.Errorf("公告不存在")
	}

	// 更新公告
	updateData := &do.Notice{
		Title:     req.Title,
		Content:   req.Content,
		Type:      int(req.Type),
		Status:    int(req.Status),
		Sort:      int(req.Sort),
		UpdatedAt: gtime.Now(),
	}

	_, err = dao.Notice.Ctx(ctx).Where("id", req.Id).Where("site_id", siteId).Data(updateData).Update()
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "更新公告失败: %v", err)
		return nil, err
	}

	middleware.LogWithTrace(ctx, "info", "更新公告成功 - Id: %d, Title: %s", req.Id, req.Title)

	return &v1.UpdateNoticeRes{}, nil
}

// DeleteNotice 删除公告
func (s *sNotice) DeleteNotice(ctx context.Context, req *v1.DeleteNoticeReq) (*v1.DeleteNoticeRes, error) {
	middleware.LogWithTrace(ctx, "info", "删除公告请求 - Id: %d", req.Id)

	// 默认站点ID为1
	siteId := int32(1)

	// 检查公告是否存在
	count, err := dao.Notice.Ctx(ctx).Where("id", req.Id).Where("site_id", siteId).Count()
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "查询公告失败: %v", err)
		return nil, err
	}
	if count == 0 {
		middleware.LogWithTrace(ctx, "error", "公告不存在 - Id: %d", req.Id)
		return nil, fmt.Errorf("公告不存在")
	}

	// 删除公告
	_, err = dao.Notice.Ctx(ctx).Where("id", req.Id).Where("site_id", siteId).Delete()
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "删除公告失败: %v", err)
		return nil, err
	}

	middleware.LogWithTrace(ctx, "info", "删除公告成功 - Id: %d", req.Id)

	return &v1.DeleteNoticeRes{}, nil
}
