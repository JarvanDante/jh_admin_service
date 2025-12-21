package service

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"

	"jh_user_service/internal/dao"
	"jh_user_service/internal/model/do"
	"jh_user_service/internal/model/entity"
)

type sAdmin struct{}

func Admin() *sAdmin {
	return &sAdmin{}
}

// GetByUsername 根据用户名获取管理员
func (s *sAdmin) GetByUsername(ctx context.Context, username string, siteId int) (*entity.Admin, error) {
	var admin *entity.Admin
	err := dao.Admin.Ctx(ctx).Where(do.Admin{
		Username: username,
		SiteId:   siteId,
		Status:   1, // 正常状态
		DeleteAt: 0, // 未删除
	}).Scan(&admin)

	return admin, err
}

// GetById 根据ID获取管理员
func (s *sAdmin) GetById(ctx context.Context, id uint) (*entity.Admin, error) {
	var admin *entity.Admin
	err := dao.Admin.Ctx(ctx).Where(do.Admin{
		Id:       id,
		Status:   1, // 正常状态
		DeleteAt: 0, // 未删除
	}).Scan(&admin)

	return admin, err
}

// UpdateLastLoginInfo 更新最后登录信息
func (s *sAdmin) UpdateLastLoginInfo(ctx context.Context, adminId uint, ip string) error {
	_, err := dao.Admin.Ctx(ctx).Where(do.Admin{Id: adminId}).Update(do.Admin{
		LastLoginIp:   ip,
		LastLoginTime: gtime.Now(),
	})
	return err
}

// AddLog 添加管理员操作日志
func (s *sAdmin) AddLog(ctx context.Context, admin *entity.Admin, message string, ip string) error {
	_, err := dao.AdminLog.Ctx(ctx).Insert(do.AdminLog{
		SiteId:        admin.SiteId,
		AdminId:       int(admin.Id),
		AdminUsername: admin.Username,
		Ip:            ip,
		Remark:        message,
		CreatedAt:     gtime.Now(),
	})
	return err
}
