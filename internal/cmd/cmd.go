package cmd

import (
	"context"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"google.golang.org/grpc"

	"jh_user_service/internal/controller/admin"
	"jh_user_service/internal/controller/user"
	"jh_user_service/internal/registry"
)

var (
	// Main is the main command.
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start grpc server of jinhuang",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 初始化 Consul 客户端
			registry.InitConsul()

			// 注册服务到 Consul
			if err := registry.RegisterService(); err != nil {
				g.Log().Fatalf(ctx, "register service failed: %v", err)
			}

			c := grpcx.Server.NewConfig()
			c.Options = append(c.Options, []grpc.ServerOption{
				grpcx.Server.ChainUnary(
					grpcx.Server.UnaryValidate,
				)}...,
			)
			s := grpcx.Server.New(c)
			user.Register(s)
			admin.Register(s)
			s.Run()
			return nil
		},
	}
)
