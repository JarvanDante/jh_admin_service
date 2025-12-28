package user

import (
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "jh_user_service/api/user/v1"
)

type Controller struct {
	v1.UnimplementedUserServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterUserServer(s.Server, &Controller{})
}

// RegisterHTTP 注册 HTTP 路由 - 用户相关路由已删除
func RegisterHTTP(s *ghttp.Server) {
	s.BindHandler("/health", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{"code": 0, "msg": "success"})
	})

	// 用户相关路由已删除
}

// 用户相关的 gRPC 和 HTTP 方法已删除
