package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"

	v1 "jh_admin_service/api/admin/v1"
	"jh_admin_service/internal/dao"
	"jh_admin_service/internal/middleware"
	"jh_admin_service/internal/model/do"
	"jh_admin_service/internal/model/entity"
	"jh_admin_service/internal/tracing"
)

type Controller struct {
	v1.UnimplementedAdminServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterAdminServer(s.Server, &Controller{})
}

// Login 管理员登录
func (*Controller) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	// 创建Jaeger span
	ctx, span := tracing.StartSpan(ctx, "admin.Login", trace.WithAttributes(
		attribute.String("username", req.Username),
		attribute.String("method", "Login"),
	))
	defer span.End()

	// 获取站点ID (这里需要根据实际情况获取，可能从上下文或配置中获取)
	siteId := 1 // 临时硬编码，实际应该从请求中获取
	tracing.SetSpanAttributes(span, attribute.Int("site_id", siteId))

	// 数据库查询span
	ctx, dbSpan := tracing.StartSpan(ctx, "db.query.admin", trace.WithAttributes(
		attribute.String("db.operation", "select"),
		attribute.String("db.table", "admin"),
	))

	var admin *entity.Admin
	err = dao.Admin.Ctx(ctx).Where("username = ? AND site_id = ?", req.Username, siteId).Scan(&admin)

	dbSpan.End()

	if err != nil {
		tracing.SetSpanError(span, err)
		tracing.SetSpanError(dbSpan, err)
		middleware.LogWithTrace(ctx, "error", "数据库查询错误: %v", err)
		return nil, fmt.Errorf("数据库查询错误: %v", err)
	}

	if admin == nil {
		tracing.AddSpanEvent(span, "admin_not_found", attribute.String("username", req.Username))
		middleware.LogWithTrace(ctx, "warning", "未找到管理员记录 - 用户名: %s, 站点ID: %d", req.Username, siteId)
		return nil, fmt.Errorf("用户名或密码错误")
	}

	tracing.SetSpanAttributes(span,
		attribute.Int64("admin_id", int64(admin.Id)),
		attribute.Int("admin_status", admin.Status),
	)

	// 检查状态（状态检查放在找到记录之后）
	if admin.Status != 1 {
		tracing.AddSpanEvent(span, "admin_status_invalid",
			attribute.String("username", req.Username),
			attribute.Int("status", admin.Status),
		)
		middleware.LogWithTrace(ctx, "warning", "管理员状态异常 - 用户名: %s, 状态: %d", req.Username, admin.Status)
		return nil, fmt.Errorf("账号已被禁用")
	}

	// 验证密码span
	ctx, authSpan := tracing.StartSpan(ctx, "auth.password_verify")
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password))
	authSpan.End()

	if err != nil {
		tracing.SetSpanError(span, err)
		tracing.AddSpanEvent(span, "password_verification_failed", attribute.String("username", req.Username))
		middleware.LogWithTrace(ctx, "warning", "密码验证失败 - 用户名: %s, 错误: %v", req.Username, err)
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 验证Google 2FA (如果开启)
	if admin.SwitchGoogle2Fa == 1 {
		tracing.AddSpanEvent(span, "google_2fa_required")
		if req.Code == "" {
			tracing.AddSpanEvent(span, "google_2fa_code_missing")
			return nil, fmt.Errorf("请输入动态验证码")
		}
		// 这里需要实现Google 2FA验证逻辑
		// valid := validateGoogle2FA(admin.Google2FaSecret, req.Code)
		// if !valid {
		//     return nil, fmt.Errorf("动态验证码错误")
		// }
	}

	// 生成JWT token span
	ctx, tokenSpan := tracing.StartSpan(ctx, "auth.generate_jwt_token")
	token, err := generateJWTToken(admin)
	tokenSpan.End()

	if err != nil {
		tracing.SetSpanError(span, err)
		tracing.SetSpanError(tokenSpan, err)
		return nil, fmt.Errorf("生成token失败: %v", err)
	}

	// 更新最后登录信息span
	ctx, updateSpan := tracing.StartSpan(ctx, "db.update.admin_login_info", trace.WithAttributes(
		attribute.String("db.operation", "update"),
		attribute.String("db.table", "admin"),
	))
	_, err = dao.Admin.Ctx(ctx).Where(do.Admin{Id: admin.Id}).Update(do.Admin{
		LastLoginIp:   getClientIP(ctx),
		LastLoginTime: gtime.Now(),
	})
	updateSpan.End()

	if err != nil {
		tracing.SetSpanError(updateSpan, err)
		middleware.LogWithTrace(ctx, "error", "更新登录信息失败: %v", err)
	}

	// 记录登录日志span
	ctx, logSpan := tracing.StartSpan(ctx, "db.insert.admin_log", trace.WithAttributes(
		attribute.String("db.operation", "insert"),
		attribute.String("db.table", "admin_log"),
	))
	err = addAdminLog(ctx, admin, "登录成功")
	logSpan.End()

	if err != nil {
		tracing.SetSpanError(logSpan, err)
		middleware.LogWithTrace(ctx, "error", "记录登录日志失败: %v", err)
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

	tracing.AddSpanEvent(span, "login_success",
		attribute.String("username", req.Username),
		attribute.Int64("admin_id", int64(admin.Id)),
	)
	tracing.SetSpanAttributes(span, attribute.Bool("success", true))

	middleware.LogWithTraceAndFields(ctx, "info", "登录成功", g.Map{
		"username": req.Username,
		"admin_id": admin.Id,
		"site_id":  siteId,
	})
	return res, nil
}

// RefreshToken 刷新token
func (*Controller) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (res *v1.RefreshTokenRes, err error) {
	middleware.LogWithTrace(ctx, "info", "刷新token请求")

	// 从上下文中获取当前用户信息 (需要中间件解析JWT)
	// 这里简化处理，实际需要从JWT中解析用户信息

	// 重新生成token
	// admin := getCurrentAdmin(ctx)
	// token, err := generateJWTToken(admin)
	// if err != nil {
	//     middleware.LogWithTrace(ctx, "error", "刷新token失败: %v", err)
	//     return nil, fmt.Errorf("刷新token失败: %v", err)
	// }

	res = &v1.RefreshTokenRes{
		Token: "new_token", // 临时返回
	}

	middleware.LogWithTrace(ctx, "info", "刷新token成功")
	return res, nil
}

// CreateAdmin 创建管理员
func (*Controller) CreateAdmin(ctx context.Context, req *v1.CreateAdminReq) (res *v1.CreateAdminRes, err error) {
	// 创建Jaeger span
	ctx, span := tracing.StartSpan(ctx, "admin.CreateAdmin", trace.WithAttributes(
		attribute.String("username", req.Username),
		attribute.String("nickname", req.Nickname),
		attribute.String("method", "CreateAdmin"),
	))
	defer span.End()

	// 获取站点ID (这里需要根据实际情况获取，可能从上下文或配置中获取)
	siteId := 1 // 临时硬编码，实际应该从请求中获取
	tracing.SetSpanAttributes(span, attribute.Int("site_id", siteId))

	middleware.LogWithTrace(ctx, "info", "创建管理员请求 - 用户名: %s, 昵称: %s", req.Username, req.Nickname)

	// 验证参数
	if req.Username == "" || req.Password == "" || req.Nickname == "" {
		tracing.AddSpanEvent(span, "validation_failed", attribute.String("reason", "missing_required_fields"))
		middleware.LogWithTrace(ctx, "error", "创建管理员参数验证失败 - 缺少必要字段")
		return nil, fmt.Errorf("用户名、密码和昵称不能为空")
	}

	// 验证用户名长度和格式
	if len(req.Username) < 4 || len(req.Username) > 12 {
		tracing.AddSpanEvent(span, "validation_failed",
			attribute.String("reason", "username_length_invalid"),
			attribute.Int("username_length", len(req.Username)),
		)
		middleware.LogWithTrace(ctx, "error", "创建管理员参数验证失败 - 用户名长度不符合要求: %d", len(req.Username))
		return nil, fmt.Errorf("用户名长度必须在4-12个字符之间")
	}

	// 验证密码长度
	if len(req.Password) < 6 || len(req.Password) > 20 {
		tracing.AddSpanEvent(span, "validation_failed",
			attribute.String("reason", "password_length_invalid"),
			attribute.Int("password_length", len(req.Password)),
		)
		middleware.LogWithTrace(ctx, "error", "创建管理员参数验证失败 - 密码长度不符合要求: %d", len(req.Password))
		return nil, fmt.Errorf("密码长度必须在6-20个字符之间")
	}

	// 验证昵称长度
	if len(req.Nickname) < 2 || len(req.Nickname) > 20 {
		tracing.AddSpanEvent(span, "validation_failed",
			attribute.String("reason", "nickname_length_invalid"),
			attribute.Int("nickname_length", len(req.Nickname)),
		)
		middleware.LogWithTrace(ctx, "error", "创建管理员参数验证失败 - 昵称长度不符合要求: %d", len(req.Nickname))
		return nil, fmt.Errorf("昵称长度必须在2-20个字符之间")
	}

	// 检查用户名是否已存在
	middleware.LogWithTrace(ctx, "info", "检查用户名是否存在 - 用户名: %s, 站点ID: %d", req.Username, siteId)

	ctx, checkSpan := tracing.StartSpan(ctx, "db.query.check_username_exists", trace.WithAttributes(
		attribute.String("db.operation", "select"),
		attribute.String("db.table", "admin"),
		attribute.String("username", req.Username),
	))

	var existingAdmin *entity.Admin
	err = dao.Admin.Ctx(ctx).Where(do.Admin{
		Username: req.Username,
		SiteId:   siteId,
	}).Scan(&existingAdmin)

	checkSpan.End()

	if err != nil {
		tracing.SetSpanError(span, err)
		tracing.SetSpanError(checkSpan, err)
		middleware.LogWithTrace(ctx, "error", "检查用户名存在性时数据库查询失败: %v", err)
		return nil, fmt.Errorf("数据库查询错误: %v", err)
	}

	if existingAdmin != nil {
		tracing.AddSpanEvent(span, "username_already_exists", attribute.String("username", req.Username))
		middleware.LogWithTrace(ctx, "warning", "用户名已存在 - 用户名: %s, 站点ID: %d", req.Username, siteId)
		return nil, fmt.Errorf("用户名已经被使用")
	}

	// 加密密码
	middleware.LogWithTrace(ctx, "info", "开始加密密码 - 用户名: %s", req.Username)

	ctx, hashSpan := tracing.StartSpan(ctx, "auth.hash_password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	hashSpan.End()

	if err != nil {
		tracing.SetSpanError(span, err)
		tracing.SetSpanError(hashSpan, err)
		middleware.LogWithTrace(ctx, "error", "密码加密失败 - 用户名: %s, 错误: %v", req.Username, err)
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	// 创建管理员
	middleware.LogWithTrace(ctx, "info", "开始创建管理员 - 用户名: %s, 昵称: %s, 角色: %d, 状态: %d",
		req.Username, req.Nickname, req.Role, req.Status)

	// 确保状态值正确 (如果为0则设为1)
	status := int(req.Status)
	if status == 0 {
		status = 1
	}

	tracing.SetSpanAttributes(span,
		attribute.Int("role", int(req.Role)),
		attribute.Int("status", status),
	)

	ctx, insertSpan := tracing.StartSpan(ctx, "db.insert.admin", trace.WithAttributes(
		attribute.String("db.operation", "insert"),
		attribute.String("db.table", "admin"),
	))

	_, err = dao.Admin.Ctx(ctx).Insert(do.Admin{
		SiteId:      siteId,
		Username:    req.Username,
		Nickname:    req.Nickname,
		Password:    string(hashedPassword),
		AdminRoleId: int(req.Role),
		Status:      status,
		CreatedAt:   gtime.Now(),
		UpdatedAt:   gtime.Now(),
	})

	insertSpan.End()

	if err != nil {
		tracing.SetSpanError(span, err)
		tracing.SetSpanError(insertSpan, err)
		middleware.LogWithTrace(ctx, "error", "创建管理员数据库操作失败 - 用户名: %s, 错误: %v", req.Username, err)
		return nil, fmt.Errorf("创建管理员失败: %v", err)
	}

	tracing.AddSpanEvent(span, "admin_created_successfully",
		attribute.String("username", req.Username),
		attribute.String("nickname", req.Nickname),
	)
	tracing.SetSpanAttributes(span, attribute.Bool("success", true))

	middleware.LogWithTrace(ctx, "info", "创建管理员成功 - 用户名: %s, 实际状态: %d", req.Username, status)

	// 记录操作日志
	// 这里需要获取当前操作的管理员信息，暂时跳过
	// err = addAdminLog(ctx, currentAdmin, "添加员工："+req.Username)

	res = &v1.CreateAdminRes{}
	return res, nil
}

// 辅助函数

// generateJWTToken 生成JWT token
func generateJWTToken(admin *entity.Admin) (string, error) {
	// 从配置文件获取JWT密钥
	jwtSecret := g.Cfg().MustGet(context.Background(), "jwt.secret").String()
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT secret not configured")
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  admin.Id, // 使用 user_id 字段名，与 Gateway 的 Claims 结构体匹配
		"admin_id": admin.Id, // 保留 admin_id 用于兼容
		"username": admin.Username,
		"site_id":  admin.SiteId,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24小时过期
		"iat":      time.Now().Unix(),
	})

	// 签名token
	tokenString, err := token.SignedString([]byte(jwtSecret))
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
