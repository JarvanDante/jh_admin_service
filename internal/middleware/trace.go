package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TraceInterceptor gRPC拦截器，用于提取和设置traceId
func TraceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 从gRPC metadata中提取traceId
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		var traceID string

		// 尝试从不同的key中获取traceId
		if values := md.Get("trace-id"); len(values) > 0 {
			traceID = values[0]
		} else if values := md.Get("x-trace-id"); len(values) > 0 {
			traceID = values[0]
		}

		if traceID != "" {
			// 将traceId添加到上下文中，供后续使用
			ctx = context.WithValue(ctx, "traceId", traceID)

			// 在日志中记录traceId的接收
			g.Log().Debugf(ctx, "接收到traceId: %s, 方法: %s", traceID, info.FullMethod)

			// 为了让GoFrame的日志系统能够识别traceId，我们可以设置到上下文中
			// 这样所有使用该上下文的日志都会包含traceId
			ctx = context.WithValue(ctx, "gf_trace_id", traceID)
		}
	}

	// 调用实际的处理函数
	return handler(ctx, req)
}

// GetTraceIDFromContext 从上下文中获取traceId
func GetTraceIDFromContext(ctx context.Context) string {
	if traceID := ctx.Value("traceId"); traceID != nil {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	return ""
}

// LogWithTrace 带traceId的日志记录 - JSON格式
func LogWithTrace(ctx context.Context, level string, message string, args ...interface{}) {
	traceID := GetTraceIDFromContext(ctx)

	// 构建结构化日志数据
	logData := g.Map{
		"service":  "user-service",
		"trace_id": traceID,
		"msg":      fmt.Sprintf(message, args...),
	}

	// 根据级别记录日志
	switch strings.ToLower(level) {
	case "debug":
		g.Log().Debug(ctx, logData)
	case "info":
		g.Log().Info(ctx, logData)
	case "warn", "warning":
		g.Log().Warning(ctx, logData)
	case "error":
		g.Log().Error(ctx, logData)
	case "fatal":
		g.Log().Fatal(ctx, logData)
	default:
		g.Log().Info(ctx, logData)
	}
}

// LogWithTraceAndFields 带traceId和额外字段的日志记录 - JSON格式
func LogWithTraceAndFields(ctx context.Context, level string, message string, fields g.Map, args ...interface{}) {
	traceID := GetTraceIDFromContext(ctx)

	// 构建结构化日志数据
	logData := g.Map{
		"service":  "user-service",
		"trace_id": traceID,
		"msg":      fmt.Sprintf(message, args...),
	}

	// 添加额外字段
	for k, v := range fields {
		logData[k] = v
	}

	// 根据级别记录日志
	switch strings.ToLower(level) {
	case "debug":
		g.Log().Debug(ctx, logData)
	case "info":
		g.Log().Info(ctx, logData)
	case "warn", "warning":
		g.Log().Warning(ctx, logData)
	case "error":
		g.Log().Error(ctx, logData)
	case "fatal":
		g.Log().Fatal(ctx, logData)
	default:
		g.Log().Info(ctx, logData)
	}
}
