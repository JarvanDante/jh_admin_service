package site

import (
	"context"
	"fmt"

	v1 "jh_admin_service/api/site/v1"
	"jh_admin_service/internal/dao"
	"jh_admin_service/internal/middleware"
	"jh_admin_service/internal/model/do"
	"jh_admin_service/internal/model/entity"
	"jh_admin_service/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type (
	sSite struct{}
)

func init() {
	service.RegisterSite(&sSite{})
}

// GetBasicSetting 获取站点基本设置
func (s *sSite) GetBasicSetting(ctx context.Context, req *v1.GetBasicSettingReq) (*v1.GetBasicSettingRes, error) {
	middleware.LogWithTrace(ctx, "info", "获取站点基本设置请求 - SiteId: %d", req.SiteId)

	// 默认站点ID为1，如果请求中有指定则使用指定的
	siteId := int32(1)
	if req.SiteId > 0 {
		siteId = req.SiteId
	}

	// 查询站点配置
	var siteConfig *entity.SiteConfig
	err := dao.SiteConfig.Ctx(ctx).Where(do.SiteConfig{
		SiteId: siteId,
	}).Scan(&siteConfig)

	if err != nil {
		middleware.LogWithTrace(ctx, "error", "查询站点配置失败: %v", err)
		return nil, fmt.Errorf("查询站点配置失败: %v", err)
	}

	if siteConfig == nil {
		middleware.LogWithTrace(ctx, "error", "站点配置不存在 - SiteId: %d", siteId)
		return nil, fmt.Errorf("站点配置不存在")
	}

	// 构建响应数据
	res := &v1.GetBasicSettingRes{
		Code:                 fmt.Sprintf("site_%d", siteId), // 站点代码
		Name:                 fmt.Sprintf("站点_%d", siteId),   // 站点名称，这里可以从其他表获取
		RegisterTimeInterval: int32(siteConfig.RegisterTimeInterval),
		SwitchRegister:       siteConfig.SwitchRegister == 1,
		IsClose:              siteConfig.IsClose == 1,
		CloseReason:          siteConfig.CloseReason,
		UrlService:           siteConfig.UrlService,
		UrlAgentPc:           siteConfig.UrlAgentPc,
		UrlMobile:            siteConfig.UrlMobile,
		UrlAgentRegister:     siteConfig.UrlAgentRegister,
		MinWithdraw:          int32(siteConfig.MinWithdraw),
		MaxWithdraw:          int32(siteConfig.MaxWithdraw),
		MobileLogo:           siteConfig.MobileLogo,
		DefaultAgentId:       int32(siteConfig.DefaultAgentId),
		DefaultAgentName:     "", // 这里可以根据 DefaultAgentId 查询代理名称
		Balance:              0,  // 这里可以查询站点余额
		BalanceReset:         0,  // 这里可以查询剩余额度
	}

	// 如果关闭原因为空，设置默认消息
	if res.CloseReason == "" {
		res.CloseReason = "您好，系统正在维护中，估计需要1个小时左右，给您造成不便，敬请谅解！"
	}

	middleware.LogWithTrace(ctx, "info", "获取站点基本设置成功 - SiteId: %d", siteId)
	return res, nil
}

// UpdateBasicSetting 更新站点基本设置
func (s *sSite) UpdateBasicSetting(ctx context.Context, req *v1.UpdateBasicSettingReq) (*v1.UpdateBasicSettingRes, error) {
	middleware.LogWithTrace(ctx, "info", "更新站点基本设置请求 - SiteId: %d", req.SiteId)

	// 默认站点ID为1，如果请求中有指定则使用指定的
	siteId := int32(1)
	if req.SiteId > 0 {
		siteId = req.SiteId
	}

	// 参数验证
	if err := s.validateUpdateBasicSetting(req); err != nil {
		middleware.LogWithTrace(ctx, "error", "参数验证失败: %v", err)
		return nil, fmt.Errorf("参数验证失败: %v", err)
	}

	// 查询现有配置
	var existingConfig *entity.SiteConfig
	err := dao.SiteConfig.Ctx(ctx).Where(do.SiteConfig{
		SiteId: siteId,
	}).Scan(&existingConfig)

	if err != nil {
		middleware.LogWithTrace(ctx, "error", "查询现有配置失败: %v", err)
		return nil, fmt.Errorf("查询现有配置失败: %v", err)
	}

	middleware.LogWithTrace(ctx, "info", "更新请求参数 - RegisterTimeInterval: %d, SwitchRegister: %t, IsClose: %t", req.RegisterTimeInterval, req.SwitchRegister, req.IsClose)
	middleware.LogWithTrace(ctx, "info", "更新请求参数 - CloseReason: '%s', UrlAgentPc: '%s', UrlService: '%s'", req.CloseReason, req.UrlAgentPc, req.UrlService)
	middleware.LogWithTrace(ctx, "info", "更新请求参数 - MinWithdraw: %d, MaxWithdraw: %d", req.MinWithdraw, req.MaxWithdraw)

	// 使用map来精确控制要更新的字段
	updateFields := g.Map{}

	// 注册时间间隔 - 只有在请求中明确指定时才更新
	if req.RegisterTimeInterval > 0 {
		updateFields["register_time_interval"] = req.RegisterTimeInterval
		middleware.LogWithTrace(ctx, "info", "添加更新字段: register_time_interval = %d", req.RegisterTimeInterval)
	}

	// 布尔值字段 - 只有在请求中明确指定时才更新
	// 注意：protobuf中布尔值的默认值是false，我们需要区分"未设置"和"设置为false"
	// 这里我们假设如果传递了请求，就更新这些布尔值
	switchRegisterValue := g.Map{"true": 1, "false": 0}[fmt.Sprintf("%t", req.SwitchRegister)]
	updateFields["switch_register"] = switchRegisterValue
	middleware.LogWithTrace(ctx, "info", "添加更新字段: switch_register = %v (原值: %t)", switchRegisterValue, req.SwitchRegister)

	isCloseValue := g.Map{"true": 1, "false": 0}[fmt.Sprintf("%t", req.IsClose)]
	updateFields["is_close"] = isCloseValue
	middleware.LogWithTrace(ctx, "info", "添加更新字段: is_close = %v (原值: %t)", isCloseValue, req.IsClose)

	// 字符串字段 - 只有非空时才更新
	if req.CloseReason != "" {
		updateFields["close_reason"] = req.CloseReason
		middleware.LogWithTrace(ctx, "info", "添加更新字段: close_reason = '%s'", req.CloseReason)
	}
	if req.UrlAgentPc != "" {
		updateFields["url_agent_pc"] = req.UrlAgentPc
		middleware.LogWithTrace(ctx, "info", "添加更新字段: url_agent_pc = '%s'", req.UrlAgentPc)
	}
	if req.UrlMobile != "" {
		updateFields["url_mobile"] = req.UrlMobile
		middleware.LogWithTrace(ctx, "info", "添加更新字段: url_mobile = '%s'", req.UrlMobile)
	}
	if req.UrlAgentRegister != "" {
		updateFields["url_agent_register"] = req.UrlAgentRegister
		middleware.LogWithTrace(ctx, "info", "添加更新字段: url_agent_register = '%s'", req.UrlAgentRegister)
	}
	if req.UrlService != "" {
		updateFields["url_service"] = req.UrlService
		middleware.LogWithTrace(ctx, "info", "添加更新字段: url_service = '%s'", req.UrlService)
	}
	if req.MobileLogo != "" {
		updateFields["mobile_logo"] = req.MobileLogo
		middleware.LogWithTrace(ctx, "info", "添加更新字段: mobile_logo = '%s'", req.MobileLogo)
	}

	// 提现金额 - 只有大于0时才更新
	middleware.LogWithTrace(ctx, "info", "检查提现金额 - MinWithdraw: %d, MaxWithdraw: %d", req.MinWithdraw, req.MaxWithdraw)

	if req.MinWithdraw > 0 {
		if req.MinWithdraw < 1 {
			updateFields["min_withdraw"] = 1
			middleware.LogWithTrace(ctx, "info", "添加更新字段: min_withdraw = 1 (原值太小)")
		} else {
			updateFields["min_withdraw"] = req.MinWithdraw
			middleware.LogWithTrace(ctx, "info", "添加更新字段: min_withdraw = %d", req.MinWithdraw)
		}
	} else {
		middleware.LogWithTrace(ctx, "info", "跳过 min_withdraw 更新，值为: %d", req.MinWithdraw)
	}

	if req.MaxWithdraw > 0 {
		if req.MaxWithdraw < 1 {
			updateFields["max_withdraw"] = 9999999
			middleware.LogWithTrace(ctx, "info", "添加更新字段: max_withdraw = 9999999 (原值太小)")
		} else {
			updateFields["max_withdraw"] = req.MaxWithdraw
			middleware.LogWithTrace(ctx, "info", "添加更新字段: max_withdraw = %d", req.MaxWithdraw)
		}
	} else {
		middleware.LogWithTrace(ctx, "info", "跳过 max_withdraw 更新，值为: %d", req.MaxWithdraw)
	}

	middleware.LogWithTrace(ctx, "info", "准备更新的字段总数: %d, 字段: %+v", len(updateFields), updateFields)

	// 如果配置不存在，创建新配置
	if existingConfig == nil {
		updateFields["site_id"] = siteId
		middleware.LogWithTrace(ctx, "info", "配置不存在，创建新配置，字段: %+v", updateFields)
		result, err := dao.SiteConfig.Ctx(ctx).Data(updateFields).Insert()
		if err != nil {
			middleware.LogWithTrace(ctx, "error", "创建站点配置失败: %v", err)
			return nil, fmt.Errorf("创建站点配置失败: %v", err)
		}
		middleware.LogWithTrace(ctx, "info", "创建站点配置成功，影响行数: %v", result)
	} else {
		// 更新现有配置 - 只更新指定的字段
		middleware.LogWithTrace(ctx, "info", "配置存在，更新现有配置，SiteId: %d, 字段: %+v", siteId, updateFields)

		// 检查是否有字段需要更新
		if len(updateFields) == 0 {
			middleware.LogWithTrace(ctx, "warning", "没有字段需要更新，跳过数据库操作")
			return &v1.UpdateBasicSettingRes{
				Message: "没有字段需要更新",
			}, nil
		}

		result, err := dao.SiteConfig.Ctx(ctx).Where(do.SiteConfig{
			SiteId: siteId,
		}).Data(updateFields).Update()

		if err != nil {
			middleware.LogWithTrace(ctx, "error", "更新站点配置失败: %v", err)
			return nil, fmt.Errorf("更新站点配置失败: %v", err)
		}

		rowsAffected, _ := result.RowsAffected()
		middleware.LogWithTrace(ctx, "info", "数据库更新完成，影响行数: %d", rowsAffected)

		if rowsAffected == 0 {
			middleware.LogWithTrace(ctx, "warning", "数据库更新影响行数为0，可能是条件不匹配或数据未变化")
		}
	}

	middleware.LogWithTrace(ctx, "info", "更新站点基本设置成功 - SiteId: %d", siteId)
	return &v1.UpdateBasicSettingRes{
		Message: "设置成功",
	}, nil
}

// validateUpdateBasicSetting 验证更新基本设置的参数
func (s *sSite) validateUpdateBasicSetting(req *v1.UpdateBasicSettingReq) error {
	// 验证URL格式（简单验证）
	if req.UrlAgentPc != "" && len(req.UrlAgentPc) < 7 {
		return fmt.Errorf("代理链接地址格式错误")
	}

	if req.UrlMobile != "" && len(req.UrlMobile) < 7 {
		return fmt.Errorf("手机域名地址格式错误")
	}

	if req.UrlAgentRegister != "" && len(req.UrlAgentRegister) < 7 {
		return fmt.Errorf("代理推广地址格式错误")
	}

	// 验证提现金额
	if req.MinWithdraw < 0 {
		return fmt.Errorf("单笔最低提现金额不能为负数")
	}

	if req.MaxWithdraw < 0 {
		return fmt.Errorf("单笔最高提现金额不能为负数")
	}

	if req.MinWithdraw > req.MaxWithdraw && req.MaxWithdraw > 0 {
		return fmt.Errorf("单笔最低提现金额不能大于最高提现金额")
	}

	// 验证关闭原因长度
	if len(req.CloseReason) > 255 {
		return fmt.Errorf("暂时关闭网站原因不能超过255个字符")
	}

	return nil
}
