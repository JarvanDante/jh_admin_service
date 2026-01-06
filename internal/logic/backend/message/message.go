package message

import (
	"context"
	"fmt"
	"jh_app_service/api/backend/message/v1"
	"jh_app_service/internal/dao"
	"jh_app_service/internal/middleware"
	"jh_app_service/internal/model/do"
	"jh_app_service/internal/model/entity"
	"jh_app_service/internal/service/backend"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type (
	sMessage struct{}
)

func init() {
	backend.RegisterMessage(&sMessage{})
}

// GetMessageList 获取消息列表
func (s *sMessage) GetMessageList(ctx context.Context, req *v1.GetMessageListReq) (*v1.GetMessageListRes, error) {
	middleware.LogWithTrace(ctx, "info", "获取消息列表请求 - Page: %d, Size: %d", req.Page, req.Size)

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

	// 获取消息列表
	var messages []*entity.Message
	total, err := dao.Message.Ctx(ctx).Where("site_id", siteId).Count()
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "获取消息总数失败: %v", err)
		return nil, err
	}

	// 分页查询
	offset := int((page - 1) * size)
	err = dao.Message.Ctx(ctx).
		Where("site_id", siteId).
		Fields("id, title, content, receiver, created_at, updated_at").
		Order("created_at DESC").
		Limit(offset, int(size)).
		Scan(&messages)
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "获取消息列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	messageList := make([]*v1.MessageItem, 0, len(messages))
	for _, msg := range messages {
		messageList = append(messageList, &v1.MessageItem{
			Id:        int32(msg.Id),
			Title:     msg.Title,
			Content:   msg.Content,
			Receiver:  msg.Receiver,
			CreatedAt: msg.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: msg.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// 获取等级列表（写死数据，后续快照用）
	gradeList := []*v1.GradeItem{
		{Id: 1, Name: "普通会员"},
		{Id: 2, Name: "VIP会员"},
		{Id: 3, Name: "SVIP会员"},
	}

	// 获取层级列表（写死数据，后续快照用）
	levelList := []*v1.LevelItem{
		{Id: 1, Name: "一级代理"},
		{Id: 2, Name: "二级代理"},
		{Id: 3, Name: "三级代理"},
	}

	// 获取状态列表
	statusList := []*v1.StatusItem{
		{Value: 1, Name: "在线会员"},
		{Value: 0, Name: "离线会员"},
	}

	middleware.LogWithTrace(ctx, "info", "获取消息列表成功 - 总数: %d, 当前页: %d", total, page)

	return &v1.GetMessageListRes{
		MessageList: messageList,
		GradeList:   gradeList,
		LevelList:   levelList,
		StatusList:  statusList,
		Total:       int32(total),
		CurrentPage: page,
		PerPage:     size,
	}, nil
}

// CreateMessage 创建消息
func (s *sMessage) CreateMessage(ctx context.Context, req *v1.CreateMessageReq) (*v1.CreateMessageRes, error) {
	middleware.LogWithTrace(ctx, "info", "创建消息请求 - Title: %s", req.Title)

	// 默认站点ID和管理员ID为1
	siteId := int32(1)
	adminId := int32(1)

	// 验证参数
	if req.Title == "" || len(req.Title) < 2 || len(req.Title) > 255 {
		middleware.LogWithTrace(ctx, "error", "标题长度必须在2-255字符之间")
		return nil, fmt.Errorf("标题长度必须在2-255字符之间")
	}
	if req.Content == "" || len(req.Content) < 2 || len(req.Content) > 255 {
		middleware.LogWithTrace(ctx, "error", "内容长度必须在2-255字符之间")
		return nil, fmt.Errorf("内容长度必须在2-255字符之间")
	}

	// 创建消息记录
	message := &do.Message{
		SiteId:    siteId,
		AdminId:   adminId,
		Title:     req.Title,
		Content:   req.Content,
		Receiver:  req.Receiver,
		CreatedAt: gtime.Now(),
		UpdatedAt: gtime.Now(),
	}

	// 保存消息
	result, err := dao.Message.Ctx(ctx).Data(message).Insert()
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "创建消息失败: %v", err)
		return nil, err
	}

	messageId, err := result.LastInsertId()
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "获取消息ID失败: %v", err)
		return nil, err
	}

	// 根据类型发送消息给用户
	err = s.sendMessageToUsers(ctx, messageId, req)
	if err != nil {
		g.Log().Errorf(ctx, "发送消息给用户失败: %v", err)
		// 不返回错误，消息已创建成功
	}

	middleware.LogWithTrace(ctx, "info", "创建消息成功 - Title: %s, ID: %d", req.Title, messageId)

	return &v1.CreateMessageRes{}, nil
}

// GetUserMessages 获取用户消息
func (s *sMessage) GetUserMessages(ctx context.Context, req *v1.GetUserMessagesReq) (*v1.GetUserMessagesRes, error) {
	middleware.LogWithTrace(ctx, "info", "获取用户消息请求 - UserId: %d", req.UserId)

	// 默认站点ID为1
	siteId := int32(1)

	page := req.Page
	if page <= 0 {
		page = 1
	}
	size := req.Size
	if size <= 0 {
		size = 20
	}

	// 定义查询结果结构体
	type UserMessageResult struct {
		Id        uint        `json:"id"`
		Title     string      `json:"title"`
		Content   string      `json:"content"`
		IsRead    int         `json:"is_read"`
		CreatedAt *gtime.Time `json:"created_at"`
	}

	// 获取用户消息列表
	var userMessages []*UserMessageResult
	total, err := dao.UserMessage.Ctx(ctx).
		LeftJoin("message m", "user_message.message_id = m.id").
		Where("user_message.site_id", siteId).
		Where("user_message.user_id", req.UserId).
		Count()
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "获取用户消息总数失败: %v", err)
		return nil, err
	}

	// 分页查询
	offset := int((page - 1) * size)
	err = dao.UserMessage.Ctx(ctx).
		LeftJoin("message m", "user_message.message_id = m.id").
		Where("user_message.site_id", siteId).
		Where("user_message.user_id", req.UserId).
		Fields("user_message.id, m.title, m.content, user_message.status as is_read, user_message.created_at").
		Order("user_message.created_at DESC").
		Limit(offset, int(size)).
		Scan(&userMessages)
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "获取用户消息失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	messageList := make([]*v1.UserMessageItem, 0, len(userMessages))
	for _, msg := range userMessages {
		messageList = append(messageList, &v1.UserMessageItem{
			Id:        int32(msg.Id),
			Title:     msg.Title,
			Content:   msg.Content,
			IsRead:    int32(msg.IsRead),
			CreatedAt: msg.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// 获取未读消息数量
	unreadCount, err := dao.UserMessage.Ctx(ctx).
		Where("site_id", siteId).
		Where("user_id", req.UserId).
		Where("status", 0). // 0表示未读
		Count()
	if err != nil {
		g.Log().Errorf(ctx, "获取未读消息数量失败: %v", err)
		unreadCount = 0
	}

	middleware.LogWithTrace(ctx, "info", "获取用户消息成功 - UserId: %d, 总数: %d", req.UserId, total)

	return &v1.GetUserMessagesRes{
		List:               messageList,
		Total:              int32(total),
		CurrentPage:        page,
		PerPage:            size,
		UnreadMessageCount: int32(unreadCount),
	}, nil
}

// GetUnreadCount 获取未读数量
func (s *sMessage) GetUnreadCount(ctx context.Context, req *v1.GetUnreadCountReq) (*v1.GetUnreadCountRes, error) {
	middleware.LogWithTrace(ctx, "info", "获取未读数量请求 - UserId: %d", req.UserId)

	// 默认站点ID为1
	siteId := int32(1)

	count, err := dao.UserMessage.Ctx(ctx).
		Where("site_id", siteId).
		Where("user_id", req.UserId).
		Where("status", 0). // 0表示未读
		Count()
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "获取未读消息数量失败: %v", err)
		return nil, err
	}

	middleware.LogWithTrace(ctx, "info", "获取未读数量成功 - UserId: %d, Count: %d", req.UserId, count)

	return &v1.GetUnreadCountRes{
		Count: int32(count),
	}, nil
}

// ReadMessage 读取消息
func (s *sMessage) ReadMessage(ctx context.Context, req *v1.ReadMessageReq) (*v1.ReadMessageRes, error) {
	middleware.LogWithTrace(ctx, "info", "读取消息请求 - UserId: %d, MessageId: %d", req.UserId, req.MessageId)

	// 默认站点ID为1
	siteId := int32(1)

	_, err := dao.UserMessage.Ctx(ctx).
		Where("site_id", siteId).
		Where("user_id", req.UserId).
		Where("message_id", req.MessageId).
		Data(g.Map{"status": 1}). // 1表示已读
		Update()
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "标记消息已读失败: %v", err)
		return nil, err
	}

	middleware.LogWithTrace(ctx, "info", "读取消息成功 - UserId: %d, MessageId: %d", req.UserId, req.MessageId)

	return &v1.ReadMessageRes{}, nil
}

// sendMessageToUsers 根据类型发送消息给用户
func (s *sMessage) sendMessageToUsers(ctx context.Context, messageId int64, req *v1.CreateMessageReq) error {
	// 默认站点ID为1
	siteId := int32(1)

	switch req.Type {
	case 1: // 指定用户
		// 解析用户名列表，创建用户消息记录
		// 这里简化处理，实际需要根据用户名查询用户ID
		userIds := []int{1, 2, 3} // 写死数据，后续快照用
		for _, userId := range userIds {
			userMessage := &do.UserMessage{
				SiteId:    siteId,
				UserId:    userId,
				MessageId: int(messageId),
				Status:    0, // 0表示未读
				CreatedAt: gtime.Now(),
			}
			_, err := dao.UserMessage.Ctx(ctx).Data(userMessage).Insert()
			if err != nil {
				g.Log().Errorf(ctx, "创建用户消息记录失败: %v", err)
			}
		}
	case 2: // 等级
		// 根据等级查询用户，创建消息记录
		// 写死数据，后续快照用
		userIds := []int{1, 2} // 假设等级1的用户
		for _, userId := range userIds {
			userMessage := &do.UserMessage{
				SiteId:    siteId,
				UserId:    userId,
				MessageId: int(messageId),
				Status:    0, // 0表示未读
				CreatedAt: gtime.Now(),
			}
			_, err := dao.UserMessage.Ctx(ctx).Data(userMessage).Insert()
			if err != nil {
				g.Log().Errorf(ctx, "创建用户消息记录失败: %v", err)
			}
		}
	case 3: // 层级
		// 根据层级查询用户，创建消息记录
		// 写死数据，后续快照用
		userIds := []int{3, 4} // 假设层级1的用户
		for _, userId := range userIds {
			userMessage := &do.UserMessage{
				SiteId:    siteId,
				UserId:    userId,
				MessageId: int(messageId),
				Status:    0, // 0表示未读
				CreatedAt: gtime.Now(),
			}
			_, err := dao.UserMessage.Ctx(ctx).Data(userMessage).Insert()
			if err != nil {
				g.Log().Errorf(ctx, "创建用户消息记录失败: %v", err)
			}
		}
	case 4: // 状态
		// 根据状态查询用户，创建消息记录
		// 写死数据，后续快照用
		userIds := []int{1, 2, 3, 4} // 假设在线用户
		for _, userId := range userIds {
			userMessage := &do.UserMessage{
				SiteId:    siteId,
				UserId:    userId,
				MessageId: int(messageId),
				Status:    0, // 0表示未读
				CreatedAt: gtime.Now(),
			}
			_, err := dao.UserMessage.Ctx(ctx).Data(userMessage).Insert()
			if err != nil {
				g.Log().Errorf(ctx, "创建用户消息记录失败: %v", err)
			}
		}
	}

	return nil
}
