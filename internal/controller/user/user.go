package user

import (
	"context"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "jh_user_service/api/user/v1"
	"jh_user_service/internal/dao"
	"jh_user_service/internal/middleware"
	"jh_user_service/internal/model/do"
	"jh_user_service/internal/service"
)

type Controller struct {
	v1.UnimplementedUserServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterUserServer(s.Server, &Controller{})
}

// RegisterHTTP 注册 HTTP 路由
func RegisterHTTP(s *ghttp.Server) {
	s.BindHandler("/health", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{"code": 0, "msg": "success"})
	})

	s.Group("/api/user", func(group *ghttp.RouterGroup) {
		group.GET("/list", (*Controller).GetListHTTP)
		group.GET("/detail/:id", (*Controller).GetOneHTTP)
		group.POST("/create", (*Controller).CreateHTTP)
		group.DELETE("/delete/:id", (*Controller).DeleteHTTP)
	})
}

func (*Controller) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	middleware.LogWithTrace(ctx, "info", "创建用户请求 - 用户名: %s, 昵称: %s", req.Passport, req.Nickname)

	_, err = dao.User.Ctx(ctx).Data(do.User{
		Username: req.Passport,
		Password: req.Password,
		Realname: req.Nickname,
	}).Insert()

	if err != nil {
		middleware.LogWithTrace(ctx, "error", "创建用户失败: %v", err)
		return nil, err
	}

	middleware.LogWithTraceAndFields(ctx, "info", "create user success", g.Map{
		"username": req.Passport,
		"nickname": req.Nickname,
	})
	return
}

func (*Controller) GetOne(ctx context.Context, req *v1.GetOneReq) (res *v1.GetOneRes, err error) {
	middleware.LogWithTrace(ctx, "info", "查询用户详情 - ID: %d", req.Id)

	user, err := service.User().GetById(ctx, req.Id)
	if err != nil {
		middleware.LogWithTrace(ctx, "error", "查询用户失败 - ID: %d, 错误: %v", req.Id, err)
		return nil, err
	}

	res = &v1.GetOneRes{
		User: user,
	}

	middleware.LogWithTrace(ctx, "info", "查询用户成功 - ID: %d", req.Id)
	return
}

func (*Controller) GetList(ctx context.Context, req *v1.GetListReq) (res *v1.GetListRes, err error) {
	middleware.LogWithTrace(ctx, "info", "查询用户列表 - 页码: %d, 大小: %d", req.Page, req.Size)

	res = &v1.GetListRes{}
	err = dao.User.Ctx(ctx).Page(int(req.Page), int(req.Size)).Scan(&res.Users)

	if err != nil {
		middleware.LogWithTrace(ctx, "error", "查询用户列表失败: %v", err)
		return nil, err
	}

	middleware.LogWithTrace(ctx, "info", "查询用户列表成功 - 返回 %d 条记录", len(res.Users))
	return
}

func (*Controller) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	middleware.LogWithTrace(ctx, "info", "删除用户请求 - ID: %d", req.Id)

	err = service.User().DeleteById(ctx, req.Id)

	if err != nil {
		middleware.LogWithTrace(ctx, "error", "删除用户失败 - ID: %d, 错误: %v", req.Id, err)
		return nil, err
	}

	middleware.LogWithTrace(ctx, "info", "删除用户成功 - ID: %d", req.Id)
	return
}

// HTTP 处理函数

func (c *Controller) CreateHTTP(r *ghttp.Request) {
	var req v1.CreateReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(g.Map{"code": 400, "msg": "invalid request"})
		return
	}

	res, err := c.Create(r.Context(), &req)
	if err != nil {
		r.Response.WriteJson(g.Map{"code": 500, "msg": err.Error()})
		return
	}

	r.Response.WriteJson(g.Map{"code": 0, "msg": "success", "data": res})
}

func (c *Controller) GetOneHTTP(r *ghttp.Request) {
	id := r.Get("id").Uint64()
	req := &v1.GetOneReq{Id: id}

	res, err := c.GetOne(r.Context(), req)
	if err != nil {
		r.Response.WriteJson(g.Map{"code": 500, "msg": err.Error()})
		return
	}

	r.Response.WriteJson(g.Map{"code": 0, "msg": "success", "data": res})
}

func (c *Controller) GetListHTTP(r *ghttp.Request) {
	page := r.Get("page", 1).Int32()
	size := r.Get("size", 10).Int32()
	req := &v1.GetListReq{Page: page, Size: size}

	res, err := c.GetList(r.Context(), req)
	if err != nil {
		r.Response.WriteJson(g.Map{"code": 500, "msg": err.Error()})
		return
	}

	r.Response.WriteJson(g.Map{"code": 0, "msg": "success", "data": res})
}

func (c *Controller) DeleteHTTP(r *ghttp.Request) {
	id := r.Get("id").Uint64()
	req := &v1.DeleteReq{Id: id}

	res, err := c.Delete(r.Context(), req)
	if err != nil {
		r.Response.WriteJson(g.Map{"code": 500, "msg": err.Error()})
		return
	}

	r.Response.WriteJson(g.Map{"code": 0, "msg": "success", "data": res})
}
