package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	v1 "jh_user_service/api/admin/v1"
	"jh_user_service/internal/dao"
	"jh_user_service/internal/model/do"
	"jh_user_service/internal/model/entity"
)

type Controller struct {
	v1.UnimplementedAdminServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterAdminServer(s.Server, &Controller{})
}

// RegisterHTTP 注册 HTTP 路由
func RegisterHTTP(s *ghttp.Server) {
	s.Group("/api/admin", func(group *ghttp.RouterGroup) {
		group.POST("/login", (*Controller).LoginHTTP)
		group.GET("/refresh-token", (*Controller).RefreshTokenHTTP)
	})
}

// Login 管理员登录
func (*Controller) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	// 获取站点ID (这里需要根据实际情况获取，可能从上下文或配置中获取)
	siteId := 1 // 临时硬编码，实际应该从请求中获取

	// 查询管理员
	var admin *entity.Admin
	err = dao.Admin.Ctx(ctx).Where(do.Admin{
		Username: req.Username,
		SiteId:   siteId,
		Status:   1, // 正常状态
		DeleteAt: 0, // 未删除
	}).Scan(&admin)

	if err != nil {
		return nil, fmt.Errorf("数据库查询错误: %v", err)
	}

	if admin == nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 验证Google 2FA (如果开启)
	if admin.SwitchGoogle2Fa == 1 {
		if req.Code == "" {
			return nil, fmt.Errorf("请输入动态验证码")
		}
		// 这里需要实现Google 2FA验证逻辑
		// valid := validateGoogle2FA(admin.Google2FaSecret, req.Code)
		// if !valid {
		//     return nil, fmt.Errorf("动态验证码错误")
		// }
	}

	// 生成JWT token
	token, err := generateJWTToken(admin)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %v", err)
	}

	// 更新最后登录信息
	_, err = dao.Admin.Ctx(ctx).Where(do.Admin{Id: admin.Id}).Update(do.Admin{
		LastLoginIp:   getClientIP(ctx),
		LastLoginTime: gtime.Now(),
	})
	if err != nil {
		g.Log().Error(ctx, "更新登录信息失败:", err)
	}

	// 记录登录日志
	err = addAdminLog(ctx, admin, "登录成功")
	if err != nil {
		g.Log().Error(ctx, "记录登录日志失败:", err)
	}

	// 获取socket地址 (从配置中获取)
	socketAddr := g.Cfg().MustGet(ctx, "workerman.host", "").String()
	if port := g.Cfg().MustGet(ctx, "workerman.port", "").String(); port != "" {
		socketAddr = socketAddr + ":" + port
	}

	res = &v1.LoginRes{
		Token:  token,
		Socket: socketAddr,
	}

	return res, nil
}

// RefreshToken 刷新token
func (*Controller) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (res *v1.RefreshTokenRes, err error) {
	// 从上下文中获取当前用户信息 (需要中间件解析JWT)
	// 这里简化处理，实际需要从JWT中解析用户信息

	// 重新生成token
	// admin := getCurrentAdmin(ctx)
	// token, err := generateJWTToken(admin)
	// if err != nil {
	//     return nil, fmt.Errorf("刷新token失败: %v", err)
	// }

	res = &v1.RefreshTokenRes{
		Token: "new_token", // 临时返回
	}

	return res, nil
}

// HTTP 处理函数
func (c *Controller) LoginHTTP(r *ghttp.Request) {
	var req v1.LoginReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(g.Map{"code": 400, "msg": "参数错误"})
		return
	}

	res, err := c.Login(r.Context(), &req)
	if err != nil {
		r.Response.WriteJson(g.Map{"code": 403, "msg": err.Error()})
		return
	}

	r.Response.WriteJson(g.Map{"code": 0, "msg": "登录成功", "data": res})
}

func (c *Controller) RefreshTokenHTTP(r *ghttp.Request) {
	var req v1.RefreshTokenReq

	res, err := c.RefreshToken(r.Context(), &req)
	if err != nil {
		r.Response.WriteJson(g.Map{"code": 666, "msg": "刷新安全令牌失败"})
		return
	}

	r.Response.WriteJson(g.Map{"code": 0, "msg": "刷新安全令牌成功", "data": res.Token})
}

// 辅助函数

// generateJWTToken 生成JWT token
func generateJWTToken(admin *entity.Admin) (string, error) {
	// JWT密钥 (实际应该从配置文件获取)
	jwtSecret := []byte("your-secret-key")

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin_id": admin.Id,
		"username": admin.Username,
		"site_id":  admin.SiteId,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24小时过期
		"iat":      time.Now().Unix(),
	})

	// 签名token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// getClientIP 获取客户端IP
func getClientIP(ctx context.Context) string {
	// 从上下文中获取HTTP请求
	if r := g.RequestFromCtx(ctx); r != nil {
		return r.GetClientIp()
	}
	return "127.0.0.1"
}

// addAdminLog 添加管理员日志
func addAdminLog(ctx context.Context, admin *entity.Admin, message string) error {
	_, err := dao.AdminLog.Ctx(ctx).Insert(do.AdminLog{
		SiteId:        admin.SiteId,
		AdminId:       int(admin.Id),
		AdminUsername: admin.Username,
		Ip:            getClientIP(ctx),
		Remark:        message,
		CreatedAt:     gtime.Now(),
	})
	return err
}
